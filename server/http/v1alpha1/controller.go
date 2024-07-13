package v1alpha1

import (
	"fmt"
	"io"
	"net/http"

	"github.com/blend/go-sdk/uuid"
	"github.com/blend/go-sdk/web"
	"github.com/mat285/boardgames/games"

	engine "github.com/mat285/boardgames/pkg/engine/v1alpha1"
	game "github.com/mat285/boardgames/pkg/game/v1alpha1"
	websockets "github.com/mat285/boardgames/pkg/websockets"
	wire "github.com/mat285/boardgames/pkg/wire/v1alpha1"
)

func (s *Server) Register(app *web.App) {

	app.POST("/api/v1alpha1/user/login", s.Login)

	app.GET("/api/v1alpha1/registry", s.ListGamesNames)
	// app.GET("/api/v1alpha1/user/:name", s.GetUserID)

	app.GET("/api/v1alpha1/user/games", s.ListUserGames)

	app.POST("/api/v1alpha1/games/:name/new", s.NewGame)
	app.POST("/api/v1alpha1/game/:id/join", s.JoinGame)
	app.POST("/api/v1alpha1/game/:id/start", s.StartGame)
	app.GET("/api/v1alpha1/game/:id/state", s.GetGameState)
	app.POST("/api/v1alpha1/game/:id/packet", s.SendPacket)

	app.RouteTree.Handle("GET", "/api/v1alpha1/websockets/:name", s.OpenWebSocketsConnection)

}

func (s *Server) ListGamesNames(r *web.Ctx) web.Result {
	return web.JSON.Result(games.ListGames())
}

func (s *Server) Login(r *web.Ctx) web.Result {
	var p game.Player
	err := r.PostBodyAsJSON(&p)
	if err != nil {
		return web.JSON.BadRequest(err)
	}
	if len(p.Username) == 0 {
		return web.JSON.BadRequest(fmt.Errorf("missing username"))
	}
	id := s.GetOrSetUserID(p.Username)
	p.ID = id
	return web.JSON.Result(p)
}

func (s *Server) NewGame(r *web.Ctx) web.Result {
	userID, username, err := s.CurrentUser(r)
	if err != nil {
		return web.JSON.BadRequest(err)
	}
	name, _ := r.Param("name")
	if len(name) == 0 {
		return web.JSON.BadRequest(fmt.Errorf("Missing `name`"))
	}

	rg, has := games.RegisteredGames()[name]
	if !has {
		return web.JSON.NotFound()
	}
	cfg := rg.Config()
	if cfg != nil {
		err := r.PostBodyAsJSON(&cfg)
		if err != nil {
			return web.JSON.BadRequest(err)
		}
	}
	g, err := rg.New(cfg)
	if err != nil {
		return web.JSON.BadRequest(err)
	}
	e, err := s.Router.NewEngine(s.Ctx, g, nil)
	if err != nil {
		return web.JSON.InternalError(err)
	}
	e = s.Router.GetEngine(e.ID)
	if e == nil {
		return web.JSON.NotFound()
	}
	if s.Router.GetClient(userID) == nil {
		s.Router.ConnectClient(s.Ctx, NewWebsocket(userID, username, nil, s.InboundPackets))
	}
	err = s.Router.Join(s.Ctx, userID, e.ID)
	if err != nil {
		return web.JSON.InternalError(err)
	}
	return web.JSON.Result(e.ID)
}

func (s *Server) OpenWebSocketsConnection(w http.ResponseWriter, r *http.Request, _ *web.Route, params web.RouteParameters) {
	username := params.Get("name")
	s.usersLock.Lock()
	userID, has := s.Users[username]
	s.usersLock.Unlock()
	if !has {
		w.WriteHeader(404)
		return
		// write err
	}

	conn, err := websockets.Upgrader().Upgrade(w, r, nil)
	if err != nil {
		return
		// write err
	}
	client, ok := s.Router.GetClient(userID).(*Websocket)
	if client != nil && ok {
		client.NewConnection(s.Ctx, conn)
	} else {
		client = NewWebsocket(userID, username, conn, s.InboundPackets)
		s.Router.ConnectClient(s.Ctx, client)
		client.Open(s.Ctx)
	}
}

type UserGame struct {
	ID   uuid.UUID
	Game string
}

func (s *Server) ListUserGames(r *web.Ctx) web.Result {
	userID, _, err := s.CurrentUser(r)
	if err != nil {
		return web.JSON.BadRequest(err)
	}
	engines := s.Router.ClientEngines(r.Context(), userID)
	res := make([]UserGame, len(engines))
	for i, e := range engines {
		res[i] = UserGame{ID: e.ID, Game: e.Game.Name()}
	}
	return web.JSON.Result(res)
}

func (s *Server) JoinGame(r *web.Ctx) web.Result {
	userID, username, err := s.CurrentUser(r)
	if err != nil {
		return web.JSON.BadRequest(err)
	}
	id, err := web.UUIDValue(r.Param("id"))
	if err != nil {
		return web.JSON.BadRequest(err)
	}
	e := s.Router.GetEngine(id)
	if e == nil {
		return web.JSON.NotFound()
	}
	if s.Router.GetClient(userID) == nil {
		s.Router.ConnectClient(s.Ctx, NewWebsocket(userID, username, nil, s.InboundPackets))
	}
	err = s.Router.Join(s.Ctx, userID, id)
	if err != nil {
		return web.JSON.InternalError(err)
	}
	return web.JSON.OK()
}

func (s *Server) StartGame(r *web.Ctx) web.Result {
	id, err := web.UUIDValue(r.Param("id"))
	if err != nil {
		return web.JSON.BadRequest(err)
	}
	e := s.Router.GetEngine(id)
	if e == nil {
		return web.JSON.NotFound()
	}
	go s.Router.StartEngine(s.Ctx, id)
	return web.JSON.OK()
}

type GameResponse struct {
	ID     uuid.UUID   `json:"id"`
	Packet wire.Packet `json:"packet"`
}

func (s *Server) GetGameState(r *web.Ctx) web.Result {
	id, err := web.UUIDValue(r.Param("id"))
	if err != nil {
		return web.JSON.BadRequest(err)
	}
	e := s.Router.GetEngine(id)
	if e == nil {
		return web.JSON.NotFound()
	}
	return s.stateResponse(e)
}

func (s *Server) stateResponse(e *engine.Engine) web.Result {
	payload := []byte{}
	data, err := e.GetStateData()
	if err == nil {
		obj, err := e.MessageProvider.SerializeState(data)
		if err != nil {
			return web.JSON.InternalError(err)
		}
		payload = obj.Data
	}
	res := wire.NewPacket(wire.OptPacketHeaderValue("game", e.ID.String()), wire.OptPacketPayload(payload))
	return web.JSON.Result(res)
}

func (s *Server) SendPacket(r *web.Ctx) web.Result {
	userID, _, err := s.CurrentUser(r)
	if err != nil {
		return web.JSON.BadRequest(err)
	}

	id, err := web.UUIDValue(r.Param("id"))
	if err != nil {
		return web.JSON.BadRequest(err)
	}

	e := s.Router.GetEngine(id)
	if e == nil {
		return web.JSON.NotFound()
	}
	body, err := io.ReadAll(r.Request.Body)
	if err != nil {
		return web.JSON.BadRequest(err)
	}
	packet, err := wire.DeserializePacket(body)
	if err != nil {
		return web.JSON.BadRequest(err)
	}
	packet.Origin = userID
	packet.Destination = e.ID
	err = e.RecieveSync(r.Context(), *packet)
	if err != nil {
		return web.JSON.BadRequest(err)
	}
	return s.stateResponse(e)
}

func (s *Server) ClientPoll(r *web.Ctx) web.Result {
	id, err := web.UUIDValue(r.Param("id"))
	if err != nil {
		return web.JSON.BadRequest(err)
	}
	client := s.Router.GetClient(id)
	if client == nil {
		return web.JSON.BadRequest(fmt.Errorf("Not connected"))
	}
	poller, ok := client.(*PollClient)
	if !ok {
		return web.JSON.InternalError(fmt.Errorf("bad type"))
	}
	packet := poller.Poll(r.Context())
	if packet == nil {
		return web.JSON.Status(http.StatusNoContent, nil)
	}
	return web.JSON.Result(packet)
}

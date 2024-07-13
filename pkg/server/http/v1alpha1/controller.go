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

	app.GET("/api/v1alpha1/registry", s.ListGamesNames)
	app.GET("/api/v1alpha1/user/:name", s.GetUserID)

	app.POST("/api/v1alpha1/games/:name/new", s.NewGame)
	app.POST("/api/v1alpha1/game/:id/join", s.JoinGame)
	app.POST("/api/v1alpha1/game/:id/start", s.StartGame)
	app.GET("/api/v1alpha1/game/:id/state", s.GetGameState)
	app.POST("/api/v1alpha1/game/:id/move/:player", s.MakeMove)

	app.RouteTree.Handle("GET", "/api/v1alpha1/websockets/:name", s.OpenWebSocketsConnection)

}

func (s *Server) ListGamesNames(r *web.Ctx) web.Result {
	return web.JSON.Result(games.ListGames())
}

func (s *Server) GetUserID(r *web.Ctx) web.Result {
	name, _ := r.Param("name")
	if len(name) == 0 {
		return web.JSON.BadRequest(fmt.Errorf("Missing `name`"))
	}

	var res uuid.UUID
	s.usersLock.Lock()
	id, has := s.Users[name]
	if has {
		res = id
	} else {
		res = uuid.V4()
		s.Users[name] = res
	}
	s.usersLock.Unlock()

	return web.JSON.Result(res)
}

func (s *Server) NewGame(r *web.Ctx) web.Result {
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
	e, err := s.Router.NewEngine(s.Ctx, g)
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
	client := NewWebsocket(userID, username, conn, s.InboundPackets)
	s.Router.ConnectClient(s.Ctx, client)
	client.Open(s.Ctx)
}

func (s *Server) JoinGame(r *web.Ctx) web.Result {
	id, err := web.UUIDValue(r.Param("id"))
	if err != nil {
		return web.JSON.BadRequest(err)
	}
	var p game.Player
	err = r.PostBodyAsJSON(&p)
	if err != nil {
		return web.JSON.BadRequest(err)
	}
	e := s.Router.GetEngine(id)
	if e == nil {
		return web.JSON.NotFound()
	}
	if s.Router.GetClient(p.ID) == nil {
		s.Router.ConnectClient(s.Ctx, NewWebsocket(p.ID, p.Username, nil, s.InboundPackets))
	}
	err = s.Router.Join(s.Ctx, p.ID, id)
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
	if e.State != nil {
		obj, err := e.Game.SerializeState(e.State.Data)
		if err != nil {
			return web.JSON.InternalError(err)
		}
		payload = obj.Data
	}

	res := wire.NewPacket(wire.OptPacketHeaderValue("game", e.ID.String()), wire.OptPacketPayload(payload))
	return web.JSON.Result(res)
}

// type GameMove struct {
// 	ID       uuid.UUID `json:"id"`
// 	PlayerID uuid.UUID `json:"player"`
// 	Move     game.Move `json:"move"`
// }

func (s *Server) MakeMove(r *web.Ctx) web.Result {
	id, err := web.UUIDValue(r.Param("id"))
	if err != nil {
		return web.JSON.BadRequest(err)
	}

	pid, err := web.UUIDValue(r.Param("player"))
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
	err = e.ReceiveFromPlayer(r.Context(), pid, *packet)
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

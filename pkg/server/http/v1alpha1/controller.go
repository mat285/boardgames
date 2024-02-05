package v1alpha1

import (
	"fmt"
	"net/http"

	"github.com/blend/go-sdk/uuid"
	"github.com/blend/go-sdk/web"
	"github.com/mat285/boardgames/games"

	game "github.com/mat285/boardgames/pkg/game/v1alpha1"
	websockets "github.com/mat285/boardgames/pkg/websockets/v1alpha1"
)

func (s *Server) Register(app *web.App) {

	app.GET("/api/v1alpha1/registry", s.ListGamesNames)

	app.POST("/api/v1alpha1/games", s.NewGame)
	app.POST("/api/v1alpha1/game/:id/join", s.JoinGame)

	app.RouteTree.Handle("GET", "/api/v1alpha1/websockets/:id", s.OpenWebSocketsConnection)

}

func (s *Server) ListGamesNames(r *web.Ctx) web.Result {
	return web.JSON.Result(games.ListGames())
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
	e, err := s.Router.NewEngine(r.Context(), g)
	if err != nil {
		return web.JSON.InternalError(err)
	}
	return web.JSON.Result(e.ID)
}

func (s *Server) OpenWebSocketsConnection(w http.ResponseWriter, r *http.Request, _ *web.Route, params web.RouteParameters) {
	strID := params.Get("id")
	id, err := uuid.Parse(strID)
	if err != nil {
		// write err
	}
	conn, err := websockets.Upgrader().Upgrade(w, r, nil)
	if err != nil {
		// write err
	}
	client := websockets.NewClient(id, "user", conn, s.Router.Receive)
	client.Start(r.Context())
	return
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
	err = s.Router.Join(r.Context(), p.ID, id)
	if err != nil {
		return web.JSON.InternalError(err)
	}
	return web.JSON.OK()
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

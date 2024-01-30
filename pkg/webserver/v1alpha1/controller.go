package v1alpha1

// import (
// 	"fmt"

// 	"github.com/blend/go-sdk/web"
// 	"github.com/mat285/boardgames/games"
// )

// func (s *Server) Register(app *web.App) {

// }

// func (s *Server) ListGamesNames(r *web.Ctx) web.Result {
// 	return web.JSON.Result(games.ListGames())
// }

// func (s *Server) NewGame(r web.Ctx) web.Result {
// 	name, _ := r.Param("name")
// 	if len(name) == 0 {
// 		return web.JSON.BadRequest(fmt.Errorf("Missing `name`"))
// 	}

// 	rg, has := games.RegisteredGames()[name]
// 	if !has {
// 		return web.JSON.NotFound()
// 	}
// 	cfg := rg.Config()
// 	if cfg != nil {
// 		err := r.PostBodyAsJSON(&cfg)
// 		if err != nil {
// 			return web.JSON.BadRequest(err)
// 		}
// 	}
// 	g, err := rg.New(cfg)
// 	if err != nil {
// 		return web.JSON.BadRequest(err)
// 	}
// 	return
// 	// e := s.NewEngine(g)
// 	// return web.JSON.Result(e.ID)
// }

// func (s *Server) JoinGame(r *web.Ctx) web.Result {
// 	id, err := web.UUIDValue(r.Param("id"))
// 	if err != nil {
// 		return web.JSON.BadRequest(err)
// 	}
// 	e := s.GetEngine(id)
// 	if e == nil {
// 		return web.JSON.NotFound()
// 	}

// 	// TODO upgrade connection and add the user

// 	return web.JSON.OK()
// }

// func (s *Server) ListGames(r *web.Ctx) web.Result {
// 	gs := make([]Game, 0, len(s.engines))

// 	for _, e := range s.engines {
// 		g := GameFromEngine(e)
// 		if g == nil {
// 			continue
// 		}
// 		gs = append(gs, *g)
// 	}

// 	return web.JSON.Result(gs)
// }

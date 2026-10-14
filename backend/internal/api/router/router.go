package router

import (
	"golbugames/internal/api/middleware"
	"golbugames/internal/websocket/multiplayer"
	"net/http"
)

type Router struct {
	mux        *http.ServeMux
	hubManager *multiplayer.HubManager
	jwtManager *middleware.JWTManager
}

func NewRouter(hubManager *multiplayer.HubManager, jwtManager *middleware.JWTManager) *Router {
	return &Router{
		mux:        http.NewServeMux(),
		hubManager: hubManager,
		jwtManager: jwtManager,
	}
}

func (r *Router) InitRoutes() {
	InitRoutesSudoku(r.mux, r.hubManager, r.jwtManager)
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}

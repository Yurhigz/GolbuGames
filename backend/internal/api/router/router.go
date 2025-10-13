package router

import (
	"golbugames/backend/internal/websocket/multiplayer"
	"net/http"
)

type Router struct {
	mux        *http.ServeMux
	hubManager *multiplayer.HubManager
}

func NewRouter(hubManager *multiplayer.HubManager) *Router {
	return &Router{
		mux:        http.NewServeMux(),
		hubManager: hubManager,
	}
}

func (r *Router) InitRoutes() {
	InitRoutesSudoku(r.mux, r.hubManager)
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}

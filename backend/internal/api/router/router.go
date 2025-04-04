package router

import "net/http"

type Router struct {
	mux *http.ServeMux
}

func NewRouter() *Router {
	return &Router{
		mux: http.NewServeMux(),
	}
}

func (r *Router) InitRoutes() {
	// Initialiser les différents groupes de routes
	InitSudokuRoutes(r.mux)
	// Autres initialisations de routes futures
	// initChessRoutes(r.mux)
	// initAuthRoutes(r.mux)
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}

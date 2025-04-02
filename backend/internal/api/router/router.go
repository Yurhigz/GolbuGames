package router

import "net/http"

type Router struct {
	mux *http.ServeMux
}

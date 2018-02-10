package router

import (
	"net/http"

	mux "github.com/dimfeld/httptreemux"
)

type Handler interface {
	Routes() []*Route
}

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc mux.HandlerFunc
}

func NewRouter(h Handler) *mux.TreeMux {
	routes := h.Routes()

	router := mux.New()
	for _, r := range routes {
		switch r.Method {
		case http.MethodGet:
			router.GET(r.Pattern, r.HandlerFunc)
		case http.MethodPost:
			router.POST(r.Pattern, r.HandlerFunc)
		case http.MethodPut:
			router.PUT(r.Pattern, r.HandlerFunc)
		default:
			panic("unsupported http method")
		}
	}

	return router
}

package v1

import (
	"net/http"

	"github.com/gorilla/mux"
	logger "github.com/marove2000/hack-and-pay/logger/v1"
	routes "github.com/marove2000/hack-and-pay/routes/v1"
)

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes.RoutesIndex {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = logger.Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)

	}
	return router
}
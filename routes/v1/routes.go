package v1

import (
	"net/http"
	"github.com/marove2000/hack-and-pay/handler/v1"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var RoutesIndex = Routes{
	Route{
		"Index",
		"GET",
		"/",
		v1.Index,
	},
	Route{
		"UserIndex",
		"GET",
		"/get/user",
		v1.UserIndex,
	},
	Route{
		"AddUser",
		"POST",
		"/post/user",
		v1.AddUser,
	},
}
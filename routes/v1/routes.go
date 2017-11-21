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
		"/v1",
		v1.Index,
	},
	Route{
		"PublicUserIndex",
		"GET",
		"/v1/get/user",
		v1.PublicUserIndex,
	},
	Route{
		"AddUser",
		"POST",
		"/v1/post/user",
		v1.AddUser,
	},
	Route{
		"GetPublicUserDetail",
		"GET",
		"/v1/get/user/{id}",
		v1.GetPublicUserDetail,
	},
}
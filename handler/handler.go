package handler

import (
	"github.com/marove2000/hack-and-pay/repository/ldap"
	"github.com/marove2000/hack-and-pay/repository/sql"
	"github.com/marove2000/hack-and-pay/router"
)

const pkg = "handler."

type Handler struct {
	repo *sql.Mysql
	ldap *ldap.LDAP
}

func New(repo *sql.Mysql, ldap *ldap.LDAP) *Handler {
	return &Handler{
		repo: repo,
		ldap: ldap,
	}
}

func (h *Handler) Routes() []*router.Route {
	return []*router.Route{
		{
			"GetUserIndex",
			"GET",
			"/v1/user",
			wrapError(h.publicUserIndex), //TODO Edit Landing-Page
		},
		{
			"Login",
			"POST",
			"/v1/login",
			wrapError(h.login),
		},
		//{
		//	"AddUser",
		//	"POST",
		//	"/v1/user",
		//	h.addUser,
		//},
	}
}

//
//var RoutesIndex = []*router.Route{
//	{
//		"Index",
//		"GET",
//		"/v1",
//		h, //TODO Edit Landing-Page
//	},
//	{
//		"publicUserIndex",
//		"GET",
//		"/v1/user",
//		v1.publicUserIndex,
//	},
//	{
//		"AddUser",
//		"POST",
//		"/v1/user",
//		v1.AddUser,
//	},
//	{
//		"GetPublicUserDetail",
//		"GET",
//		"/v1/user/{id}",
//		v1.GetPublicUserDetail,
//	},
//	{
//		"GetAuthentication",
//		"POST",
//		"/v1/authentication",
//		v1.GetAuthentication,
//	},
//	{
//		"ChangeBalance",
//		"POST",
//		"/v1/user/{id}/transaction",
//		v1.ChangeBalance,
//	},
//	{
//		"PutProduct",
//		"PUT",
//		"/v1/put/product",
//		v1.PutProduct,
//	},
//	{
//		"GetProductIndex",
//		"GET",
//		"/v1/get/product",
//		v1.GetProductIndex,
//	},
//}

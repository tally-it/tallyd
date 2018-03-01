// Tally it - daemon
//
// A tool for hack spaces, clubs, groups and companies to manage a digital tally sheet.
//
// Schemes: http
// Version: 0.1
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
//
// Security:
// - bearer
//
// SecurityDefinitions:
//   bearer:
//    type: apiKey
//    name: Authorization
//    in: header
//
// swagger:meta
package handler

import (
	"github.com/tally-it/tallyd/log"
	"github.com/tally-it/tallyd/repository/ldap"
	"github.com/tally-it/tallyd/repository/sql"
	"github.com/tally-it/tallyd/router"
)

var pkgLogger = log.New("handler")

type Handler struct {
	repo       *sql.Mysql
	ldap       *ldap.LDAP
	authorizer authorizer
}

func New(repo *sql.Mysql, ldap *ldap.LDAP, authorizer authorizer) *Handler {
	return &Handler{
		repo:       repo,
		ldap:       ldap,
		authorizer: authorizer,
	}
}

func (h *Handler) Routes() []*router.Route {
	return []*router.Route{
		// swagger:operation GET /v1/user users GetUserIndex
		// ---
		// authorization: none
		//
		// summary: Lists all users with their public data.
		// responses:
		//   '200':
		//     $ref: '#/definitions/UserSlice'
		{
			"GetUserIndex",
			"GET",
			"/v1/user",
			wrap(h.publicUserIndex),
		},
		// swagger:operation POST /v1/user/login users userLogin
		// ---
		// summary: Used to login
		// description:
		//   >
		//     asd
		//     very
		//     long
		//     description
		// parameters:
		// - name: body
		//   required: true
		//   description: Optional description in *Markdown*
		//   in: body
		//   schema:
		//     $ref: '#/definitions/LoginRequestBody'
		// responses:
		//   '200':
		//     $ref: '#/definitions/LoginResponse'
		//   '400':
		//     $ref: '#/definitions/errorResponse'
		{
			"Login",
			"POST",
			"/v1/login",
			wrap(h.login),
		},
		// swagger:operation POST /v1/user users signUp
		// ---
		// summary: Registers a new user.
		// description: If author length is between 6 and 8, Error Not Found (404) will be returned.
		// parameters:
		// - name: body
		//   required: true
		//   description: Optional description in *Markdown*
		//   in: body
		//   schema:
		//     $ref: '#/definitions/AddUserRequestBody'
		// responses:
		//   '200':
		//     $ref: '#/definitions/AddUserResponseBody'
		//   '400':
		//     $ref: '#/definitions/errorResponse'
		//   '500':
		//     $ref: '#/definitions/errorResponse'
		{
			"AddUser",
			"POST",
			"/v1/user",
			wrap(h.signUp),
		},
		// swagger:operation POST /v1/user/{id} users userProfile
		// ---
		// summary: List the repositories owned by the given author.
		// description: If author length is between 6 and 8, Error Not Found (404) will be returned.
		// parameters:
		// - name: id
		//   in: path
		//   description: userId
		//   type: string
		//   required: true
		// - name: body
		//   required: true
		//   description: Optional description in *Markdown*
		//   in: body
		//   schema:
		//     $ref: '#/definitions/EditUserRequestBody'
		// responses:
		//   '200':
		//     $ref: '#/definitions/EditUserRequestBody'
        //   '400':
		//     $ref: '#/definitions/errorResponse'
        //   '500':
		//     $ref: '#/definitions/errorResponse'
		{
			"GetUserDetail",
			"GET",
			"/v1/user/:id",
			wrap(h.authorizer.Authorize(h.getUserDetail, authTypeNone)),
		},
		{
			"AddTransaction",
			"POST",
			"/v1/user/:id/transaction",
			wrap(h.authorizer.Authorize(h.addTransaction, authTypeAll)),
		},
		{
			"UpdateUser",
			"PUT",
			"/v1/user/:id",
			wrap(h.authorizer.Authorize(h.editUser, authTypePassword)),
		},
		{
			"DeleteUser",
			"DELETE",
			"/v1/user/:id",
			wrap(h.authorizer.Authorize(h.deleteUser, authTypePassword)),
		},
		{
			"GetProductIndex",
			"GET",
			"/v1/product",
			wrap(h.productIndex),
		},
		{
			"GetProductDetail",
			"GET",
			"/v1/product/:productID",
			wrap(h.productDetail),
		},
		{
			"AddProduct",
			"POST",
			"/v1/product",
			wrap(h.authorizer.Authorize(h.addProduct, authTypeAll)),
		},
		{
			"UpdateProduct",
			"PUT",
			"/v1/product/:productID",
			wrap(h.authorizer.Authorize(h.editProduct, authTypeAll)),
		},
		{
			"DeleteProduct",
			"DELETE",
			"/v1/product/:productID",
			wrap(h.authorizer.Authorize(h.deleteProduct, authTypeAll)),
		},
		{
			"ChangeStock",
			"POST",
			"/v1/product/:productID/stock",
			wrap(h.authorizer.Authorize(h.changeStock, authTypeAll)),
		},
		{
			"AddCategory",
			"POST",
			"/v1/category",
			wrap(h.authorizer.Authorize(h.addCategory, authTypeAll)),
		},
	}
}

package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/marove2000/hack-and-pay/config"
	"github.com/marove2000/hack-and-pay/contract"
	"github.com/marove2000/hack-and-pay/ctxutil"
	"github.com/marove2000/hack-and-pay/errors"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-validator/validator"
	"github.com/shopspring/decimal"
)

func (h *Handler) publicUserIndex(ctx context.Context, r *http.Request, pathParams map[string]string) (interface{}, error) {
	logger := pkgLogger.ForFunc(ctx, "publicUserIndex")
	logger.Debug("enter handler")

	// get all user data
	users, err := h.repo.GetUsersWithBalance(ctx)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (h *Handler) getUserDetail(ctx context.Context, r *http.Request, pathParams map[string]string) (interface{}, error) {
	logger := pkgLogger.ForFunc(ctx, "getUserDetail")
	logger.Debug("enter handler")

	// read id
	userID := pathParams["id"]
	id, err := strconv.Atoi(userID)
	if err != nil {
		logger.WithError(err).Error("failed to parse id")
		return nil, errors.BadRequest(err.Error())
	}

	var user *contract.User
	if id == ctxutil.GetUserID(ctx) || ctxutil.GetAdminStatus(ctx) {
		user, err = h.repo.GetUserWithBalance(ctx, id)
	} else {
		user, err = h.repo.GetPublicUserDataByUserID(ctx, id)
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (h *Handler) signUp(ctx context.Context, r *http.Request, pathParams map[string]string) (interface{}, error) {
	logger := pkgLogger.ForFunc(ctx, "signUp")
	logger.Debug("enter handler")

	user := &contract.AddUserRequestBody{}

	// get body data
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(user)
	if err != nil {
		logger.WithError(err).Error("failed to parse body")
		return nil, errors.BadRequest(err.Error())
	}
	defer r.Body.Close()

	// validate data
	if err = validator.Validate(user); err != nil {
		logger.WithError(err).Warn("bad request")
		return nil, errors.BadRequest(err.Error())
	}

	// check if ldap is active
	var id int
	switch {
	case h.ldap.IsActive():
		// login with ldap
		err = h.ldap.Login(ctx, user.Name, user.Password)
		if err != nil {
			return nil, err
		}

		// login correct, create user in DB
		id, err = h.repo.AddLDAPUser(ctx, user.Name, user.Email, false)
		if err != nil {
			return nil, err
		}

	default:
		// ldap not active create user account
		id, err = h.repo.AddLocalUser(ctx, user.Name, user.Email, user.Password, false)
		if err != nil {
			return nil, err
		}
	}

	return &contract.AddUserResponseBody{UserID: id}, err
}

func (h *Handler) deleteUser(ctx context.Context, r *http.Request, pathParams map[string]string) (interface{}, error) {
	logger := pkgLogger.ForFunc(ctx, "deleteUser")
	logger.Debug("enter handler")

	// read id
	userID := pathParams["id"]
	id, err := strconv.Atoi(userID)
	if err != nil {
		logger.WithError(err).Error(err.Error())
		return nil, errors.BadRequest("failed to parse id")
	}

	// user can only delete themself; admins can delete everybody
	// check if data, which can be only edited by admin, are changed
	if ctxutil.GetAdminStatus(ctx) == true || ctxutil.GetUserID(ctx) == id {
		err = h.repo.DeleteUser(ctx, id)
		if err != nil {
			return nil, errors.BadRequest("bad request")
		}
	} else {
		return nil, errors.Unauthorized()
	}

	return nil, nil
}

func (h *Handler) editUser(ctx context.Context, r *http.Request, pathParams map[string]string) (interface{}, error) {
	logger := pkgLogger.ForFunc(ctx, "editUser")
	logger.Debug("enter handler")

	// read id
	userID := pathParams["id"]
	id, err := strconv.Atoi(userID)
	if err != nil {
		logger.WithError(err).Error(err.Error())
		return nil, errors.BadRequest("failed to parse id")
	}

	user := &contract.EditUserRequestBody{}

	// get body data
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(user)
	if err != nil {
		logger.WithError(err).Error("failed to parse body")
		return nil, errors.BadRequest(err.Error())
	}
	defer r.Body.Close()

	// validate data
	if err = validator.Validate(user); err != nil {
		logger.WithError(err).Warn("bad request")
		return nil, errors.BadRequest(err.Error())
	}

	dbUser, err := h.repo.GetUserWithBalance(ctx, id)
	logger.Info(dbUser.IsBlocked)

	// check if data, which can be only edited by admin, are changed
	if ctxutil.GetAdminStatus(ctx) == false && (user.IsBlocked != dbUser.IsBlocked || user.IsAdmin != dbUser.IsAdmin) {
		logger.Warn("admin status is needed to execute the update")
		return nil, errors.Unauthorized()
	}

	err = h.repo.EditUser(ctx, id, user.Name, user.Email, user.IsBlocked, user.IsAdmin)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (h *Handler) addTransaction(ctx context.Context, r *http.Request, pathParams map[string]string) (interface{}, error) {
	logger := pkgLogger.ForFunc(ctx, "addTransaction")
	logger.Debug("enter handler")

	// read id
	userID := pathParams["id"]
	id, err := strconv.Atoi(userID)
	if err != nil {
		logger.WithError(err).Error(err.Error())
		return nil, errors.BadRequest("failed to parse id")
	}

	// read body
	transaction := &contract.ChangeBalanceRequestBody{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(transaction)
	if err != nil {
		logger.WithError(err).Error("failed to parse body")
		return nil, errors.BadRequest(err.Error())
	}
	defer r.Body.Close()

	if err = validator.Validate(transaction); err != nil {
		logger.WithError(err).Warn("bad request")
		return nil, errors.BadRequest(err.Error())
	}

	switch {
	case transaction.SKU != 0 && transaction.Value.Cmp(decimal.Zero) != 0:
		logger.Warn("both SKU and value are set")
		return nil, errors.BadRequest("both SKU and value are set")
	case transaction.SKU == 0 && transaction.Value.Cmp(decimal.Zero) == 0:
		logger.Warn("both SKU and value are zero")
		return nil, errors.BadRequest("both SKU and value are zero")
	}

	if (id == ctxutil.GetUserID(ctx) && id == transaction.UserID) || ctxutil.GetAdminStatus(ctx) {
		// add data
		err = h.repo.AddTransaction(ctx, transaction)
		if err != nil {
			return nil, err
		}
	} else {
		logger.Warn("bad request")
		return nil, errors.BadRequest("Bad request")
	}

	return nil, nil
}

func (h *Handler) login(ctx context.Context, r *http.Request, pathParams map[string]string) (interface{}, error) {
	logger := pkgLogger.ForFunc(ctx, "login")
	logger.Debug("enter handler")

	user := &contract.LoginRequestBody{}

	// get body data
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(user)
	if err != nil {
		logger.WithError(err).Error("failed to parse body")
		return nil, errors.BadRequest(err.Error())
	}
	defer r.Body.Close()

	if err = validator.Validate(user); err != nil {
		logger.WithError(err).Warn("bad request")
		return nil, errors.BadRequest(err.Error())
	}

	u, err := h.repo.GetPublicUserDataByUserName(ctx, user.Name)
	if err != nil {
		return nil, err
	}

	// get full user
	fullUser, err := h.repo.GetUserWithBalance(ctx, u.UserID)
	if err != nil {
		return nil, err
	}

	if h.ldap.IsActive() {
		err = h.ldap.Login(ctx, user.Name, user.Password)
		if err != nil {
			return nil, err
		}
	} else {
		err = h.repo.Login(ctx, user.Name, user.Password)
		if err != nil {
			return nil, err
		}
	}

	// read config
	conf := config.ReadConfig()

	// create JWT
	// TODO extract other auth-methods from db and insert into token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID":    u.UserID,
		"isBlocked": fullUser.IsBlocked,
		"isAdmin":   fullUser.IsAdmin,
		"authType":  "passwd",
		"exp":       time.Now().Add(time.Second * time.Duration(conf.JWT.ValidTime)),
	})
	tokenString, err := token.SignedString([]byte(conf.JWT.Secret))
	if err != nil {
		return nil, err
	}

	return &contract.LoginResponse{
		JWT:    tokenString,
		UserID: u.UserID,
	}, nil
}

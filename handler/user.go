package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
	"strconv"

	"github.com/marove2000/hack-and-pay/contract"
	"github.com/marove2000/hack-and-pay/errors"
	"github.com/marove2000/hack-and-pay/config"

	"github.com/go-validator/validator"
	"github.com/dgrijalva/jwt-go"
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

	// get all user data
	user, err := h.repo.GetPublicUserDataByUserID(ctx, id)
	if err != nil {
		logger.WithError(err).Error("failed to get user data")
		return nil, errors.BadRequest(err.Error())
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
		id, err = h.repo.AddLDAPUser(ctx, user)
		if err != nil {
			return nil, err
		}

	default:
		// ldap not active create user account
		id, err = h.repo.AddLocalUser(ctx, user)
		if err != nil {
			return nil, err
		}
	}
	return &contract.AddUserResponseBody{UserID: id}, err
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
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID":    u.UserID,
		"Name":      u.Name,
		"IsActive":  u.IsActive,
		"IsAdmin":   u.IsAdmin,
		"ExpiresAt": time.Now().Add(time.Second * time.Duration(conf.JWTValidTime)),
	})
	tokenString, err := token.SignedString([]byte(conf.JWTSecret))
	if err != nil {
		return nil, err
	}

	return tokenString, nil
}

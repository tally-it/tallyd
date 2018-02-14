package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/tally-it/tallyd/ctxutil"
	"github.com/tally-it/tallyd/errors"

	"github.com/dgrijalva/jwt-go"
	mux "github.com/dimfeld/httptreemux"
	"github.com/pborman/uuid"
)

type authType int

const (
	authTypeNone     authType = iota
	authTypePassword
	authTypePin
	authTypeAll
)

type handlerFunc func(ctx context.Context, r *http.Request, pathParams map[string]string) (interface{}, error)

type authorizer interface {
	Authorize(fn handlerFunc, authType authType) handlerFunc
}

type JWTAuthorizer struct {
	Secret string
}

func (a *JWTAuthorizer) Authorize(fn handlerFunc, authType authType) handlerFunc {
	return func(ctx context.Context, r *http.Request, pathParams map[string]string) (interface{}, error) {
		logger := pkgLogger.ForFunc(ctx, "Authorize")
		logger.Debug("enter wrapper")

		claims, err := a.extractAndValidateJWT(ctx, authType, r.Header)
		if err != nil {
			return nil, err
		}

		ctx, err = extractAdminStatus(ctx, claims)
		if err != nil {
			logger.WithError(err).Error("failed to extract admin status")
			return nil, errors.InternalServerError("something went wrong", nil)
		}

		ctx, err = extractUserID(ctx, claims)
		if err != nil {
			logger.WithError(err).Error("failed to extract userID")
			return nil, errors.InternalServerError("something went wrong", nil)
		}

		ctx, err = extractIsBlockedStatus(ctx, claims)
		if err != nil {
			logger.WithError(err).Error("failed to extract user is blocked status")
			return nil, errors.InternalServerError("something went wrong", nil)
		}

		ctx, err = extractAuthType(ctx, claims)
		if err != nil {
			logger.WithError(err).Error("failed to extract auth type")
			return nil, errors.InternalServerError("something went wrong", nil)
		}

		return fn(ctx, r, pathParams)
	}
}

func wrap(fn handlerFunc) mux.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, params map[string]string) {
		defer r.Body.Close()

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		ctx := context.Background()

		var corrId string
		if corrIdSlice := r.Header["correlation-id"]; len(corrIdSlice) != 0 {
			uid := uuid.Parse(corrIdSlice[0])
			if uid == nil {
				corrId = uid.String()
			}
		}
		if corrId == "" {
			corrId = uuid.NewRandom().String()
		}

		resp, err := fn(ctxutil.InjectCorrelationId(ctx, corrId), r, params)
		if err != nil {
			w.WriteHeader(err.(*errors.Error).Status)
			json.NewEncoder(w).Encode(err)
			return
		}

		if resp == nil {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		json.NewEncoder(w).Encode(resp)
	}
}

func (a *JWTAuthorizer) extractAndValidateJWT(ctx context.Context, authType authType, header http.Header) (jwt.MapClaims, error) {
	logger := pkgLogger.ForFunc(ctx, "extractAndValidateJWT")
	logger.Debug("enter wrapper")

	var JWT string
	bearerToken := strings.Split(header.Get("authorization"), " ")
	if len(bearerToken) != 2 {
		if authType == authTypeNone {
			return nil, nil
		}

		logger.Error("failed to read JWT token from header, maybe it's missing?")
		return nil, errors.Unauthenticated()
	}

	JWT = bearerToken[1]

	token, err := jwt.Parse(JWT, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			err := fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			return nil, err
		}
		return []byte(a.Secret), nil
	})
	if err != nil {
		logger.WithError(err).Warn("failed to parse JWT")
		return nil, errors.Unauthorized()
	}

	if !token.Valid {
		logger.WithError(err).Warn("token is invalid")
		return nil, errors.Unauthorized()
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		logger.Error("type assertion failed")
		return nil, errors.InternalServerError("something went wrong", nil)
	}

	switch authType {
	case authTypePassword:
		if claims["authType"] != "passwd" {
			logger.WithField("authType", claims["authType"]).Error("wrong authType, passwd needed")
			return nil, errors.InternalServerError("wrong authType", nil)
		}
	case authTypePin:
		if claims["authType"] != "pin" {
			logger.WithField("authType", claims["authType"]).Error("wrong authType, pin needed")
			return nil, errors.InternalServerError("wrong authType", nil)
		}
	}
	return claims, nil
}

func extractAdminStatus(ctx context.Context, claims jwt.MapClaims) (context.Context, error) {
	if claims == nil {
		return ctx, nil
	}

	isAdminInterface, ok := claims["isAdmin"]
	if !ok {
		return nil, fmt.Errorf("failed to find claim")
	}
	isAdmin, ok := isAdminInterface.(bool)
	if !ok {
		return nil, fmt.Errorf("type assertion failed")
	}
	return ctxutil.InjectAdminStatus(ctx, isAdmin), nil
}

func extractUserID(ctx context.Context, claims jwt.MapClaims) (context.Context, error) {
	if claims == nil {
		return ctx, nil
	}

	userIDInterface, ok := claims["userID"]
	if !ok {
		return nil, fmt.Errorf("failed to find claim")
	}
	userID, ok := userIDInterface.(float64)
	if !ok {
		return nil, fmt.Errorf("type assertion failed")
	}

	return ctxutil.InjectUserID(ctx, int(userID)), nil
}

func extractIsBlockedStatus(ctx context.Context, claims jwt.MapClaims) (context.Context, error) {
	if claims == nil {
		return ctx, nil
	}

	isBlockedInterface, ok := claims["isBlocked"]
	if !ok {
		return nil, fmt.Errorf("failed to find claim")
	}
	isBlocked, ok := isBlockedInterface.(bool)
	if !ok {
		return nil, fmt.Errorf("type assertion failed")
	}
	return ctxutil.InjectUserIsBlocked(ctx, isBlocked), nil
}

func extractAuthType(ctx context.Context, claims jwt.MapClaims) (context.Context, error) {
	if claims == nil {
		return ctx, nil
	}

	authTypeInterface, ok := claims["authType"]
	if !ok {
		return nil, fmt.Errorf("failed to find claim")
	}
	authType, ok := authTypeInterface.(string)
	if !ok {
		return nil, fmt.Errorf("type assertion failed")
	}
	return ctxutil.InjectAuthType(ctx, authType), nil
}

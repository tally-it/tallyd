package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/marove2000/hack-and-pay/errors"

	mux "github.com/dimfeld/httptreemux"
	"github.com/marove2000/hack-and-pay/ctxutil"
	"github.com/pborman/uuid"
	"strings"
	"github.com/dgrijalva/jwt-go"
	"fmt"
)

type authType int

const (
	authTypeAll      authType = iota
	authTypePassword
	authTypePin
)

type superFunc func(ctx context.Context, r *http.Request, pathParams map[string]string) (interface{}, error)

type authorizer interface {
	Authorize(fn superFunc, authType authType) superFunc
}

type JWTAuthorizer struct {
	Secret string
}

func (a *JWTAuthorizer) Authorize(fn superFunc, authType authType) superFunc {
	return func(ctx context.Context, r *http.Request, pathParams map[string]string) (interface{}, error) {
		logger := pkgLogger.ForFunc(ctx, "authorize")
		logger.Debug("enter authorize wrapper")

		// read JWT
		var JWT string
		bearerToken := strings.Split(r.Header.Get("authorization"), " ")
		if len(bearerToken) == 2 {
			JWT = bearerToken[1]
		} else {
			logger.Error("failed to read JWT token from header")
			return nil, errors.Unauthenticated()
		}

		token, err := jwt.Parse(JWT, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				err := fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
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


		return fn(ctx, r, pathParams)
	}
}

func wrap(fn superFunc) mux.HandlerFunc {
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

func extractAdminStatus(ctx context.Context, claims jwt.MapClaims) (context.Context, error) {
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
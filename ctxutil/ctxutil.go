package ctxutil

import (
	"context"
)

type correlationIdKey struct{}
type adminKey struct{}
type userIDKey struct{}
type userIsBlockedKey struct{}
type authTypeKey struct{}

func InjectAuthType(ctx context.Context, authType string) context.Context {
	return context.WithValue(ctx, authTypeKey{}, authType)
}

func GetAuthType(ctx context.Context) string {
	authType := ctx.Value(authTypeKey{})
	if authType != nil {
		return authType.(string)
	}
	return ""
}

func InjectAdminStatus(ctx context.Context, isAdmin bool) context.Context {
	return context.WithValue(ctx, adminKey{}, isAdmin)
}

func GetAdminStatus(ctx context.Context) bool {
	status := ctx.Value(adminKey{})
	if status != nil {
		return status.(bool)
	}
	return false
}

func InjectUserID(ctx context.Context, userID int) context.Context {
	return context.WithValue(ctx, userIDKey{}, userID)
}

func GetUserID(ctx context.Context) int {
	userID := ctx.Value(userIDKey{})
	if userID != nil {
		return userID.(int)
	}
	return 0
}

func InjectUserIsBlocked(ctx context.Context, isBlocked bool) context.Context {
	return context.WithValue(ctx, userIsBlockedKey{}, isBlocked)
}

func GetUserIsBlockedStatus(ctx context.Context) bool {
	status := ctx.Value(userIsBlockedKey{})
	if status != nil {
		return status.(bool)
	}
	return true
}

func InjectCorrelationId(ctx context.Context, corrId string) context.Context {
	return context.WithValue(ctx, correlationIdKey{}, corrId)
}

func GetCorrelationId(ctx context.Context) string {
	strInt := ctx.Value(correlationIdKey{})
	if strInt != nil {
		return strInt.(string)
	}

	return ""
}

package ctxutil

import "context"

type correlationIdKey struct{}

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

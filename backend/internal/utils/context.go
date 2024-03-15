package utils

import (
	"context"
)

type userIDKey struct{}

func NewContext(ctx context.Context, uid string) context.Context {
	return context.WithValue(ctx, userIDKey{}, uid)
}

func FromContext(ctx context.Context) string {
	uid, ok := ctx.Value(userIDKey{}).(string)
	if !ok {
		return ""
	}
	return uid
}

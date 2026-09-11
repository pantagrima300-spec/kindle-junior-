package auth

import (
	"context"

	"competition-backend/internal/storage"
)

type contextKey string

const userContextKey contextKey = "authenticated_user"

func WithUser(ctx context.Context, user *storage.AuthRecord) context.Context {
	return context.WithValue(ctx, userContextKey, user)
}

func GetUser(ctx context.Context) (*storage.AuthRecord, bool) {
	user, ok := ctx.Value(userContextKey).(*storage.AuthRecord)
	return user, ok
}

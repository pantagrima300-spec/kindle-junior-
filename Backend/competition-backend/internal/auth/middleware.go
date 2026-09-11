package auth

import (
	"net/http"
	"strings"

	"competition-backend/internal/storage"
)

func Middleware(pb *storage.PocketBaseClient) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			authHeader := r.Header.Get("Authorization")

			if authHeader == "" {
				http.Error(w, "authentication required", http.StatusUnauthorized)
				return
			}

			token := strings.TrimSpace(authHeader)

			if token == "" {
				http.Error(w, "authentication required", http.StatusUnauthorized)
				return
			}

			user, err := pb.ValidateToken(r.Context(), token)

			if err != nil {
				http.Error(w, "invalid or expired token", http.StatusUnauthorized)
				return
			}

			ctx := WithUser(r.Context(), &user.Record)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

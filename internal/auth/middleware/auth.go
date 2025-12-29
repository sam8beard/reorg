package middleware

import (
	"context"
	"github.com/sam8beard/reorg/internal/auth"
	"net/http"
	"strings"
)

type CtxKey string

const (
	CtxKeyUserID CtxKey = "user_id"
	CtxKeyGuest  CtxKey = "guest"
)

func AuthMiddleware(jwtService *auth.JWTService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Authorization header required", http.StatusUnauthorized)
				return
			}
			bearerToken := strings.Split(authHeader, " ")
			if len(bearerToken) != 2 || bearerToken[0] != "Bearer" {
				http.Error(w, "Invalid authorization header", http.StatusUnauthorized)
				return
			}

			userID, isGuest, err := jwtService.ValidateToken(bearerToken[1])
			if err != nil {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), CtxKeyUserID, userID)
			ctx = context.WithValue(ctx, CtxKeyGuest, isGuest)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

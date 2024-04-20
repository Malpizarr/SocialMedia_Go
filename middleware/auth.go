package middleware

import (
	"SocialMedia/utils"
	"context"
	"net/http"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}

		tokenString := authHeader[7:]
		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		ctx := r.Context()

		ctx = context.WithValue(ctx, "username", claims.Username)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

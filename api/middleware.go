package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"netshop/main/tools"
	"strings"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request to \x1b[33m'%s'\033[0m from %s", r.RequestURI, r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}

func RequireAuth(handler http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			tools.RespondWithError(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		var token string
		_, err := fmt.Sscanf(authHeader, "Bearer %s", &token)
		if err != nil {
			tools.RespondWithError(w, "Invalid authorization header", http.StatusBadRequest)
			return
		}

		claims, err := tools.ParseJWTToken(strings.TrimSpace(token))
		if err != nil {
			tools.RespondWithError(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, "user", claims)
		r = r.WithContext(ctx)

		handler(w, r)
	})
}

func RequireGuest(handler http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" {
			tools.RespondWithError(w, "User already authorized", http.StatusBadRequest)
			return
		}

		handler(w, r)
	})
}

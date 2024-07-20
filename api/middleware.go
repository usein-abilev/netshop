package api

import (
	"log"
	"net/http"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request to \x1b[33m'%s'\033[0m from %s", r.RequestURI, r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}

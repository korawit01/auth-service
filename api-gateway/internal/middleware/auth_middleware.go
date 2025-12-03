package middleware

import (
	"net/http"
)

// AuthMiddleware is a thin placeholder to demonstrate where token verification would occur.
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// No-op for now; place to validate JWTs from the Authorization header.
		next.ServeHTTP(w, r)
	})
}

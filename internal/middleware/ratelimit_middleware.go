package middleware

import (
	"net/http"

	"golang.org/x/time/rate"
)

func rateLimitMiddleware(r rate.Limit, b int, next http.Handler) http.Handler {
	limiter := rate.NewLimiter(r, b)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			http.Error(w, "rate limit exceded", http.StatusTooManyRequests)

			return
		}
		next.ServeHTTP(w, r)
	})
}

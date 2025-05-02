package middlewares

import (
	"cloud/Balancer/internal/service"
	"net/http"
)

func RateLimitMiddleware(rl *service.RateLimiter, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientID := rl.GetClientID(r)

		if !rl.Allow(clientID) {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}

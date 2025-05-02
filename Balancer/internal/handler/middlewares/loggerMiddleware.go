package middlewares

import (
	"cloud/Balancer/internal/service"
	"net/http"
	"time"
)

type responseWriter struct {
	http.ResponseWriter
	Status int
}

func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		service.AppLogger.Printf(
			"Started %s %s from %s",
			r.Method,
			r.URL.Path,
			r.RemoteAddr,
		)

		next.ServeHTTP(w, r)

		service.AppLogger.Printf(
			"Completed %s %s in %v",
			r.Method,
			r.URL.Path,
			time.Since(start),
		)
	})
}

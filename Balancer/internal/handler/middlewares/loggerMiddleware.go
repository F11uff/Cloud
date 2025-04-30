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

func LoggerMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		rw := &responseWriter{w, http.StatusOK}

		service.AppLogger.Printf(
			"Started %s %s from %s",
			r.Method,
			r.URL.Path,
			r.RemoteAddr,
		)

		next.ServeHTTP(w, r)

		duration := time.Since(startTime)
		service.AppLogger.Printf(
			"Completed %s %s | Status: %d | Duration: %v",
			r.Method,
			r.URL.Path,
			rw.Status,
			duration,
		)

	})
}

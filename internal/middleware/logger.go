package middleware

import (
	"log"
	"net/http"
	"time"
)

// LoggerMiddleware ghi log mọi request
func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("[%s] %s %s | %v", r.RemoteAddr, r.Method, r.URL.Path, time.Since(start))
	})
}

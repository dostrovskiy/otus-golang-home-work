package internalhttp

import (
	"net/http"
	"time"
)

type Middleware struct {
	logger Logger
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func NewMiddleware(log Logger) *Middleware {
	return &Middleware{logger: log}
}

func (mw *Middleware) loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ww := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next(ww, r)
		duration := time.Since(start)
		clientIP := r.RemoteAddr
		if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
			clientIP = ip
		}
		mw.logger.Info("%s [%s] %s %s %s %d %d \"%s\"",
			clientIP,
			time.Now().Format(time.RFC1123),
			r.Method,
			r.URL.Path,
			r.Proto,
			ww.statusCode,
			duration,
			r.UserAgent(),
		)
	})
}

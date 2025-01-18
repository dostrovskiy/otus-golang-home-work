package internalhttp

import (
	"context"
	"net/http"
	"time"

	"github.com/oapi-codegen/runtime/strictmiddleware/nethttp"
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

func (mw *Middleware) loggingMiddleware(next nethttp.StrictHTTPHandlerFunc, operationID string) nethttp.StrictHTTPHandlerFunc {
    return func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (response interface{}, err error) {
        start := time.Now()
        ww := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
        response, err = next(ctx, ww, r, request)
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
        return response, err
    }
}
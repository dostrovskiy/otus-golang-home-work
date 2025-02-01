package internalhttp

import (
	"context"
	"net/http"
	"regexp"
	"time"

	metrics "github.com/dostrovskiy/otus-golang-home-work/hw12_13_14_15_16_calendar/internal/metrics"
	"github.com/oapi-codegen/runtime/strictmiddleware/nethttp"
)

// ID or UUID.
var idRegex = regexp.MustCompile(`/[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}|/\d+`)

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

func (mw *Middleware) loggingMiddleware(next nethttp.StrictHTTPHandlerFunc, _ string) nethttp.StrictHTTPHandlerFunc {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (response interface{}, err error) { //nolint:lll
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

func normalizePath(path string) string {
	return idRegex.ReplaceAllString(path, ":id")
}

func getStatusGroup(status int) string {
	switch {
	case status >= 200 && status < 300:
		return "2xx"
	case status >= 400 && status < 500:
		return "4xx"
	case status >= 500:
		return "5xx"
	default:
		return "other"
	}
}

func (mw *Middleware) metricsMiddleware(next nethttp.StrictHTTPHandlerFunc, _ string) nethttp.StrictHTTPHandlerFunc {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (response interface{}, err error) { //nolint:lll
		endpoint := normalizePath(r.URL.Path)
		metrics.EventsConcurrentRequests.WithLabelValues(r.Method, endpoint).Inc()
		start := time.Now()
		rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		defer func() {
			status := getStatusGroup(rw.statusCode)
			duration := time.Since(start).Seconds()
			metrics.EventsRequestsCounter.WithLabelValues(r.Method, endpoint, status).Inc()
			metrics.EventsRequestsDuration.WithLabelValues(r.Method, endpoint).Observe(duration)
			metrics.EventsConcurrentRequests.WithLabelValues(r.Method, endpoint).Dec()
		}()
		return next(ctx, rw, r, request)
	}
}

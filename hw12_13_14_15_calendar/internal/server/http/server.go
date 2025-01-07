package internalhttp

import (
	"context"
	"net/http"
	"time"
)

type Server struct {
	app        Application
	httpServer *http.Server
	logger     Logger
	hr         *Handler
	mw         *Middleware
}

type Logger interface {
	Error(format string, a ...any)
	Warn(format string, a ...any)
	Info(format string, a ...any)
	Debug(format string, a ...any)
}

type Application interface { // TODO
}

func NewServer(logger Logger, app Application) *Server {
	return &Server{logger: logger, app: app, hr: NewHandler(logger), mw: NewMiddleware(logger)}
}

func (s *Server) routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/hello", s.mw.loggingMiddleware(s.hr.helloHandler))
	return mux
}

func (s *Server) Start(ctx context.Context, addr string) error {
	srv := http.Server{
		Addr:              addr,
		Handler:           s.routes(),
		ReadHeaderTimeout: 5 * time.Second,
	}
	s.httpServer = &srv

	go func() {
		s.logger.Info("Starting HTTP server on: %s...", s.httpServer.Addr)
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Error("failed to start HTTP server: %s", err.Error())
		}
	}()

	<-ctx.Done()

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	s.logger.Info("Stopping HTTP server...")
	stopCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := s.httpServer.Shutdown(stopCtx); err != nil {
		s.logger.Error("failed to stop HTTP server: %s", err.Error())
		return err
	}
	s.logger.Info("HTTP server stopped.")
	return nil
}

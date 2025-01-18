//go:generate oapi-codegen --config=oapicfg.yml --package=internalhttp ../../../api/event-openapi.yml

package internalhttp

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/dostrovskiy/otus-golang-home-work/hw12_13_14_15_16_calendar/internal/storage"
)

type Server struct {
	app        Application
	httpServer *http.Server
	logger     Logger
	mw         *Middleware
}

type Logger interface {
	Error(format string, a ...any)
	Warn(format string, a ...any)
	Info(format string, a ...any)
	Debug(format string, a ...any)
}

type Application interface {
	GetEvent(ctx context.Context, id string) (*storage.Event, error)
	CreateEvent(ctx context.Context, event *storage.Event) error
	UpdateEvent(ctx context.Context, id string, event *storage.Event) error
	DeleteEvent(ctx context.Context, id string) error
	FindEventsForPeriod(ctx context.Context, start time.Time, end time.Time) ([]*storage.Event, error)
}

var _ StrictServerInterface = (*Server)(nil)

func NewServer(logger Logger, app Application) *Server {
	return &Server{logger: logger, app: app, mw: NewMiddleware(logger)}
}

func (s *Server) Start(ctx context.Context, addr string) error {
	swagger, err := GetSwagger()
	if err != nil {
		return fmt.Errorf("error loading swagger spec: %w", err)
	}

	// Clear out the servers array in the swagger spec, that skips validating
	// that server names match. We don't know how this thing will be run.
	swagger.Servers = nil

	mux := http.NewServeMux()

	mws := []StrictMiddlewareFunc{s.mw.loggingMiddleware}
	h := NewStrictHandler(s, mws)

	HandlerFromMux(h, mux)

	srv := http.Server{
		Addr:              addr,
		Handler:           mux,
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

func mapHttpToStorageEvent(e Event) (*storage.Event, error) {
	notifyBefore, err := time.ParseDuration(*e.NotifyBefore)
	if err != nil {
		return nil, fmt.Errorf("error while parsing duration [%s]: %w", *e.NotifyBefore, err)
	}
	return &storage.Event{
		ID:           *e.Id,
		Title:        *e.Title,
		Start:        *e.Start,
		End:          *e.End,
		Description:  *e.Description,
		OwnerID:      *e.OwnerId,
		NotifyBefore: notifyBefore,
	}, nil
}

func mapStorageToHttpEvent(e storage.Event) (Event, error) {
	notifyBefore := e.NotifyBefore.String()
	return Event{
		Id:           &e.ID,
		Title:        &e.Title,
		Start:        &e.Start,
		End:          &e.End,
		Description:  &e.Description,
		OwnerId:      &e.OwnerID,
		NotifyBefore: &notifyBefore,
	}, nil
}

func mapStorageToHttpEvents(es []*storage.Event) ([]Event, error) {
	events := make([]Event, 0, len(es))
	for _, e := range es {
		event, err := mapStorageToHttpEvent(*e)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	return events, nil
}

func (s *Server) PostEvent(ctx context.Context, request PostEventRequestObject) (PostEventResponseObject, error) {
	event, err := mapHttpToStorageEvent(*request.Body)
	if err != nil {
		return nil, fmt.Errorf("error while posting event [%+v]: %w", *request.Body, err)
	}
	err = s.app.CreateEvent(ctx, event)
	if err != nil {
		return nil, err
	}
	return PostEvent201JSONResponse(*request.Body), nil
}

func (s *Server) DeleteEventId(ctx context.Context, request DeleteEventIdRequestObject) (DeleteEventIdResponseObject, error) {
	return DeleteEventId204Response{}, s.app.DeleteEvent(ctx, request.Id)
}

func (s *Server) GetEventId(ctx context.Context, request GetEventIdRequestObject) (GetEventIdResponseObject, error) {
	sevent, err := s.app.GetEvent(ctx, request.Id)
	if err != nil {
		return nil, fmt.Errorf("error while getting event by id [%s]: %w", request.Id, err)
	}
	event, err := mapStorageToHttpEvent(*sevent)
	if err != nil {
		return nil, fmt.Errorf("error while getting event by id [%s]: %w", request.Id, err)
	}
	return GetEventId200JSONResponse(event), nil
}

func (s *Server) PutEventId(ctx context.Context, request PutEventIdRequestObject) (PutEventIdResponseObject, error) {
	event, err := mapHttpToStorageEvent(*request.Body)
	if err != nil {
		return nil, fmt.Errorf("error while putting event [%+v]: %w", *request.Body, err)
	}
	err = s.app.UpdateEvent(ctx, event.ID, event)
	if err != nil {
		return nil, fmt.Errorf("error while putting event [%+v]: %w", *request.Body, err)
	}
	return PutEventId200JSONResponse(*request.Body), nil
}

func (s *Server) GetEventsByPeriod(ctx context.Context, request GetEventsByPeriodRequestObject) (GetEventsByPeriodResponseObject, error) {
	events, err := s.app.FindEventsForPeriod(ctx, request.Params.Start, request.Params.End)
	if err != nil {
		return nil, fmt.Errorf("error while getting events by period [%+v]: %w", request.Params, err)
	}
	hevents, err := mapStorageToHttpEvents(events)
	if err != nil {
		return nil, fmt.Errorf("error while getting events by period [%+v]: %w", request.Params, err)
	}
	return GetEventsByPeriod200JSONResponse(hevents), nil
}

func (s *Server) GetHello(ctx context.Context, request GetHelloRequestObject) (GetHelloResponseObject, error) {
	return GetHello204Response{}, nil
}

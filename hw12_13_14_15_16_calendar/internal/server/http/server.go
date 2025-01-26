//go:generate oapi-codegen --config=oapicfg.yml --package=internalhttp ../../../api/event-openapi.yml

package internalhttp

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/dostrovskiy/otus-golang-home-work/hw12_13_14_15_16_calendar/internal/app"
	"github.com/dostrovskiy/otus-golang-home-work/hw12_13_14_15_16_calendar/internal/storage"
)

var (
	ErrPostingEvent = func(event Event, err error) error {
		return fmt.Errorf("error while posting event [%+v]: %w", event, err)
	}
	ErrPuttingEvent = func(event Event, err error) error {
		return fmt.Errorf("error while putting event [%+v]: %w", event, err)
	}
	ErrGettingEventByID = func(id string, err error) error {
		return fmt.Errorf("error while getting event by id [%s]: %w", id, err)
	}
	ErrGettingEventByParams = func(params GetEventsByPeriodParams, err error) error {
		return fmt.Errorf("error while getting events by period [%+v]: %w", params, err)
	}
)

type Server struct {
	app        app.Application
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

var _ StrictServerInterface = (*Server)(nil)

func NewServer(logger Logger, app app.Application) *Server {
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

func mapHTTPToStorageEvent(e Event) (*storage.Event, error) {
	var notifyBefore time.Duration
	var err error
	notifyBeforeStr := deref(e.NotifyBefore)
	if notifyBeforeStr != "" {
		notifyBefore, err = time.ParseDuration(deref(e.NotifyBefore))
		if err != nil {
			return nil, fmt.Errorf("error while parsing duration [%s]: %w", deref(e.NotifyBefore), err)
		}
	}
	return &storage.Event{
		ID:           deref(e.Id),
		Title:        deref(e.Title),
		Start:        deref(e.Start),
		End:          deref(e.End),
		Description:  deref(e.Description),
		OwnerID:      deref(e.OwnerId),
		NotifyBefore: notifyBefore,
		NotifyStart:  deref(e.NotifyStart),
		Notified:     false,
	}, nil
}

func deref[T any](p *T) T {
	if p == nil {
		var zero T
		return zero
	}
	return *p
}

func mapStorageToHTTPEvent(e storage.Event) (Event, error) { //nolint:unparam
	notifyBefore := e.NotifyBefore.String()
	return Event{
		Id:           &e.ID,
		Title:        &e.Title,
		Start:        &e.Start,
		End:          &e.End,
		Description:  &e.Description,
		OwnerId:      &e.OwnerID,
		NotifyBefore: &notifyBefore,
		NotifyStart:  &e.NotifyStart,
		Notified:     &e.Notified,
	}, nil
}

func mapStorageToHTTPEvents(es []*storage.Event) ([]Event, error) {
	events := make([]Event, 0, len(es))
	for _, e := range es {
		event, err := mapStorageToHTTPEvent(*e)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	return events, nil
}

func (s *Server) PostEvent(ctx context.Context, request PostEventRequestObject) (PostEventResponseObject, error) {
	event, err := mapHTTPToStorageEvent(*request.Body)
	if err != nil {
		return nil, ErrPostingEvent(*request.Body, err)
	}
	storageEvent, err := s.app.CreateEvent(ctx, event)
	if err != nil {
		return nil, ErrPostingEvent(*request.Body, err)
	}
	httpEvent, err := mapStorageToHTTPEvent(*storageEvent)
	if err != nil {
		return nil, ErrPostingEvent(*request.Body, err)
	}
	return PostEvent201JSONResponse(httpEvent), nil
}

func (s *Server) DeleteEventId(ctx context.Context, request DeleteEventIdRequestObject) (DeleteEventIdResponseObject, error) { //nolint:lll,revive,stylecheck
	return DeleteEventId204Response{}, s.app.DeleteEvent(ctx, request.Id)
}

func (s *Server) GetEventId(ctx context.Context, request GetEventIdRequestObject) (GetEventIdResponseObject, error) { //nolint:revive,lll,stylecheck
	sevent, err := s.app.GetEvent(ctx, request.Id)
	if err != nil {
		return nil, ErrGettingEventByID(request.Id, err)
	}
	if sevent == nil {
		return GetEventId200JSONResponse{}, nil
	}
	event, err := mapStorageToHTTPEvent(*sevent)
	if err != nil {
		return nil, ErrGettingEventByID(request.Id, err)
	}
	return GetEventId200JSONResponse(event), nil
}

func (s *Server) PutEventId(ctx context.Context, request PutEventIdRequestObject) (PutEventIdResponseObject, error) { //nolint:revive,lll,stylecheck
	event, err := mapHTTPToStorageEvent(*request.Body)
	if err != nil {
		return nil, ErrPuttingEvent(*request.Body, err)
	}
	storageEvent, err := s.app.UpdateEvent(ctx, request.Id, event)
	if err != nil {
		return nil, ErrPuttingEvent(*request.Body, err)
	}
	httpEvent, err := mapStorageToHTTPEvent(*storageEvent)
	if err != nil {
		return nil, ErrPuttingEvent(*request.Body, err)
	}
	return PutEventId200JSONResponse(httpEvent), nil
}

func (s *Server) GetEventsByPeriod(ctx context.Context, request GetEventsByPeriodRequestObject) (GetEventsByPeriodResponseObject, error) { //nolint:lll
	events, err := s.app.FindEventsByPeriod(ctx, request.Params.Start, request.Params.End)
	if err != nil {
		return nil, ErrGettingEventByParams(request.Params, err)
	}
	hevents, err := mapStorageToHTTPEvents(events)
	if err != nil {
		return nil, ErrGettingEventByParams(request.Params, err)
	}
	return GetEventsByPeriod200JSONResponse(hevents), nil
}

func (s *Server) GetHello(_ context.Context, _ GetHelloRequestObject) (GetHelloResponseObject, error) {
	return GetHello204Response{}, nil
}

func (s *Server) GetEventsForNotify(ctx context.Context, request GetEventsForNotifyRequestObject) (GetEventsForNotifyResponseObject, error) { //nolint:lll
	events, err := s.app.FindEventsForNotify(ctx, request.Params.NotifyDate, request.Params.Notified)
	if err != nil {
		return nil, err
	}
	hevents, err := mapStorageToHTTPEvents(events)
	if err != nil {
		return nil, err
	}
	return GetEventsForNotify200JSONResponse(hevents), nil
}

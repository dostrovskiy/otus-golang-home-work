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

func mapHTTPToStorageEvent(e Event) *storage.Event {
	return &storage.Event{
		ID:           deref(e.ID),
		Title:        deref(e.Title),
		Start:        deref(e.Start),
		End:          deref(e.End),
		Description:  deref(e.Description),
		OwnerID:      deref(e.OwnerID),
		NotifyBefore: time.Duration(deref(e.NotifyBefore)),
		NotifyStart:  deref(e.NotifyStart),
		Notified:     false,
	}
}

func deref[T any](p *T) T {
	if p == nil {
		var zero T
		return zero
	}
	return *p
}

func mapStorageToHTTPEvent(e storage.Event) (Event, error) { //nolint:unparam
	notifyBefore := int64(e.NotifyBefore)
	return Event{
		ID:           &e.ID,
		Title:        &e.Title,
		Start:        &e.Start,
		End:          &e.End,
		Description:  &e.Description,
		OwnerID:      &e.OwnerID,
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
	event := mapHTTPToStorageEvent(*request.Body)
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

func (s *Server) DeleteEventID(ctx context.Context, request DeleteEventIDRequestObject) (DeleteEventIDResponseObject, error) { //nolint:lll
	return DeleteEventID204Response{}, s.app.DeleteEvent(ctx, request.ID)
}

func (s *Server) GetEventID(ctx context.Context, request GetEventIDRequestObject) (GetEventIDResponseObject, error) { //nolint:lll
	sevent, err := s.app.GetEvent(ctx, request.ID)
	if err != nil {
		return nil, ErrGettingEventByID(request.ID, err)
	}
	if sevent == nil {
		return GetEventID200JSONResponse{}, nil
	}
	event, err := mapStorageToHTTPEvent(*sevent)
	if err != nil {
		return nil, ErrGettingEventByID(request.ID, err)
	}
	return GetEventID200JSONResponse(event), nil
}

func (s *Server) PutEventID(ctx context.Context, request PutEventIDRequestObject) (PutEventIDResponseObject, error) { //nolint:lll
	event := mapHTTPToStorageEvent(*request.Body)
	storageEvent, err := s.app.UpdateEvent(ctx, request.ID, event)
	if err != nil {
		return nil, ErrPuttingEvent(*request.Body, err)
	}
	httpEvent, err := mapStorageToHTTPEvent(*storageEvent)
	if err != nil {
		return nil, ErrPuttingEvent(*request.Body, err)
	}
	return PutEventID200JSONResponse(httpEvent), nil
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

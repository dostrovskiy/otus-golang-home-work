package app

import (
	"context"
	"fmt"
	"time"

	"github.com/dostrovskiy/otus-golang-home-work/hw12_13_14_15_16_calendar/internal/storage"
	memorystorage "github.com/dostrovskiy/otus-golang-home-work/hw12_13_14_15_16_calendar/internal/storage/memory"
	sqlstorage "github.com/dostrovskiy/otus-golang-home-work/hw12_13_14_15_16_calendar/internal/storage/sql"
	"github.com/google/uuid"
)

type Application interface {
	GetEvent(ctx context.Context, id string) (*storage.Event, error)
	CreateEvent(ctx context.Context, event *storage.Event) (*storage.Event, error)
	UpdateEvent(ctx context.Context, id string, event *storage.Event) (*storage.Event, error)
	DeleteEvent(ctx context.Context, id string) error
	FindEventsByPeriod(ctx context.Context, start time.Time, end time.Time) ([]*storage.Event, error)
	FindEventsForNotify(ctx context.Context, notifyDate time.Time, notified bool) ([]*storage.Event, error)
	AddEventNotification(ctx context.Context, event *storage.Event) (*storage.Notification, error)
}

type App struct {
	logger  Logger
	storage Storage
}

type Logger interface {
	Error(format string, a ...any)
	Warn(format string, a ...any)
	Info(format string, a ...any)
	Debug(format string, a ...any)
}

type Storage interface {
	Add(event *storage.Event) (*storage.Event, error)
	Get(id string) (*storage.Event, error)
	FindByPeriod(start, end time.Time) ([]*storage.Event, error)
	FindForNotify(notifyDate time.Time, notified bool) ([]*storage.Event, error)
	Update(id string, event *storage.Event) (*storage.Event, error)
	Delete(id string) error
	Close() error
	AddNotification(notification *storage.Notification) (*storage.Notification, error)
}

func New(logger Logger, storage Storage) *App {
	return &App{logger: logger, storage: storage}
}

func NewEventStorage(storageType, dsn string) (Storage, error) {
	switch storageType {
	case "memory":
		return memorystorage.New(), nil
	case "sql":
		sqlstorage := sqlstorage.New(dsn)
		if err := sqlstorage.Open(context.Background()); err != nil {
			return nil, fmt.Errorf("failed to connect to storage: %w", err)
		}
		return sqlstorage, nil
	default:
		return nil, fmt.Errorf("unknown storage type: %s", storageType)
	}
}

func (a *App) CreateEvent(_ context.Context, event *storage.Event) (*storage.Event, error) {
	a.logger.Debug("App create event %+v", event)
	if event.ID == "" {
		event.ID = uuid.New().String()
	}
	if event.NotifyStart.IsZero() && event.NotifyBefore > 0 {
		event.NotifyStart = event.Start.Add(-event.NotifyBefore)
	}
	if !event.NotifyStart.IsZero() && event.NotifyBefore == 0 {
		event.NotifyBefore = event.Start.Sub(event.NotifyStart)
	}
	return a.storage.Add(event)
}

func (a *App) GetEvent(_ context.Context, id string) (*storage.Event, error) {
	a.logger.Debug("App get event by ID [%s]", id)
	return a.storage.Get(id)
}

func (a *App) FindEventsByPeriod(_ context.Context, start time.Time, end time.Time) ([]*storage.Event, error) {
	a.logger.Debug("App find events by period start [%v] and end [%v]", start, end)
	return a.storage.FindByPeriod(start, end)
}

func (a *App) UpdateEvent(_ context.Context, id string, event *storage.Event) (*storage.Event, error) {
	return a.storage.Update(id, event)
}

func (a *App) DeleteEvent(_ context.Context, id string) error {
	a.logger.Debug("App delete event id: [%s]", id)
	return a.storage.Delete(id)
}

func (a *App) FindEventsForNotify(_ context.Context, notifyDate time.Time, notified bool) ([]*storage.Event, error) {
	a.logger.Debug("App find events for notify by date [%v] and sign [%v]", notifyDate, notified)
	return a.storage.FindForNotify(notifyDate, notified)
}

func (a *App) AddEventNotification(_ context.Context, event *storage.Event) (*storage.Notification, error) {
	a.logger.Debug("App add event notification [%+v]", event)
	notification := &storage.Notification{
		ID:         uuid.New().String(),
		EventID:    event.ID,
		Title:      event.Title,
		EventStart: event.Start,
		EventEnd:   event.End,
		OwnerID:    event.OwnerID,
	}
	return a.storage.AddNotification(notification)
}

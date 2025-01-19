package app

import (
	"context"
	"fmt"
	"time"

	"github.com/dostrovskiy/otus-golang-home-work/hw12_13_14_15_16_calendar/internal/storage"
)

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
	GetForPeriod(start, end time.Time) ([]*storage.Event, error)
	Update(id string, event *storage.Event) (*storage.Event, error)
	Delete(id string) error
	Close() error
}

func New(logger Logger, storage Storage) *App {
	return &App{logger: logger, storage: storage}
}

func (a *App) CreateEvent(_ context.Context, event *storage.Event) (*storage.Event, error) {
	a.logger.Debug(fmt.Sprintf("App create event %v", event))
	return a.storage.Add(event)
}

func (a *App) GetEvent(_ context.Context, id string) (*storage.Event, error) {
	a.logger.Debug(fmt.Sprintf("App get event by ID [%s]", id))
	return a.storage.Get(id)
}

func (a *App) FindEventsForPeriod(_ context.Context, start time.Time, end time.Time) ([]*storage.Event, error) {
	a.logger.Debug(fmt.Sprintf("Find event by period start [%s] and end [%s]", start, end))
	return a.storage.GetForPeriod(start, end)
}

func (a *App) UpdateEvent(_ context.Context, id string, event *storage.Event) (*storage.Event, error) {
	return a.storage.Update(id, event)
}

func (a *App) DeleteEvent(_ context.Context, id string) error {
	return a.storage.Delete(id)
}

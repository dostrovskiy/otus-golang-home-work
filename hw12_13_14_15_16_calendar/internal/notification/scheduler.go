package notification //nolint

import (
	"context"
	"time"

	"github.com/dostrovskiy/otus-golang-home-work/hw12_13_14_15_16_calendar/internal/app"
	internalmessagebroker "github.com/dostrovskiy/otus-golang-home-work/hw12_13_14_15_16_calendar/internal/messagebroker"
	metrics "github.com/dostrovskiy/otus-golang-home-work/hw12_13_14_15_16_calendar/internal/metrics"
)

type Scheduler struct {
	logger   app.Logger
	app      app.Application
	producer *internalmessagebroker.Producer
}

func NewScheduler(logger app.Logger, app app.Application, sender *internalmessagebroker.Producer) *Scheduler {
	return &Scheduler{logger: logger, app: app, producer: sender}
}

func (s *Scheduler) Start(ctx context.Context, interval time.Duration) error {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			s.SendNotifications(ctx)
			s.DeletePastEvents(ctx)
		}
	}
}

func (s *Scheduler) SendNotifications(ctx context.Context) {
	metrics.EventsCleanupStatus.Set(1)
	defer func() {
		metrics.EventsCleanupStatus.Set(0)
	}()

	notifyDate := time.Now()
	s.logger.Debug("Sending notifications for %s", notifyDate)
	events, err := s.app.FindEventsForNotify(ctx, notifyDate, false)
	if err != nil {
		s.logger.Error("failed to get events: %s", err.Error())
		return
	}
	if len(events) == 0 {
		s.logger.Debug("No events to notify.")
		return
	}
	for _, event := range events {
		err = s.producer.Send(event)
		if err != nil {
			s.logger.Error("Scheduler failed to send notification: %s", err.Error())
			continue
		}
		event.Notified = true
		_, err = s.app.UpdateEvent(ctx, event.ID, event)
		if err != nil {
			s.logger.Error("Scheduler failed to send notification: %s", err.Error())
			continue
		}
		metrics.EventsNotificationCounter.Inc()
	}
}

func (s *Scheduler) DeletePastEvents(ctx context.Context) {
	start := time.Now()
	metrics.EventsCleanupStatus.Set(1)
	defer func() {
		metrics.EventsCleanupStatus.Set(0)
		metrics.EventsCleanupDuration.Observe(time.Since(start).Seconds())
	}()

	startDate := time.Unix(0, 0)
	endDate := time.Now().AddDate(-1, 0, 0)
	s.logger.Debug("Deleting past events from %s to %s", startDate, endDate)
	events, err := s.app.FindEventsByPeriod(ctx, startDate, endDate)
	if err != nil {
		s.logger.Error("failed to get events: %w", err)
		return
	}
	s.logger.Debug("Deleting %d events...", len(events))
	for _, event := range events {
		err = s.app.DeleteEvent(ctx, event.ID)
		if err != nil {
			s.logger.Error("failed to delete event: %w", err)
		}
	}
	metrics.EventsDeletedCounter.Add(float64(len(events)))
}
//nolint
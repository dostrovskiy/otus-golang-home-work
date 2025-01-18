package sqlstorage

import (
	"context"
	"fmt"
	"time"

	"github.com/dostrovskiy/otus-golang-home-work/hw12_13_14_15_16_calendar/internal/storage"
	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v4/stdlib" //nolint
	"github.com/jmoiron/sqlx"
)

type Storage struct {
	db  *sqlx.DB
	dsn string // "postgres://user:password@localhost:5432/dbname?sslmode=disable"
}

func New(dsn string) *Storage {
	return &Storage{
		dsn: dsn,
	}
}

func (s *Storage) Add(event *storage.Event) error {
	if event.ID == "" {
		event.ID = uuid.New().String()
	}
	_, err := s.db.NamedExec(
		`insert into events 
	       (id, title, event_start, event_end, description, owner_id, notify_before) 
	     values (:id, :title, :event_start, :event_end, :description, :owner_id, :notify_before)`,
		map[string]interface{}{
			"id":            event.ID,
			"title":         event.Title,
			"event_start":   event.Start,
			"event_end":     event.End,
			"description":   event.Description,
			"owner_id":      event.OwnerID,
			"notify_before": event.NotifyBefore,
		})
	if err != nil {
		return fmt.Errorf("error while adding event [%v]: %w", event, err)
	}
	return nil
}

func (s *Storage) Close() error {
	return s.db.Close()
}

func (s *Storage) Delete(id string) error {
	_, err := s.db.NamedExec("delete from events where id = :id", map[string]interface{}{"id": id})
	if err != nil {
		return fmt.Errorf("error while deleting event by id [%s]: %w", id, err)
	}
	return nil
}

func (s *Storage) Get(id string) (*storage.Event, error) {
	rows, err := s.db.NamedQuery(
		`select id, title, event_start, event_end, description, owner_id, notify_before
	       from events
	      where id = :id`, map[string]interface{}{"id": id})
	if err != nil {
		return nil, fmt.Errorf("error while getting event by id [%s]: %w", id, err)
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}
	var event storage.Event
	err = rows.StructScan(&event)
	if err != nil {
		return nil, fmt.Errorf("error while getting event by id [%s]: %w", id, err)
	}
	return &event, nil
}

func (s *Storage) GetForPeriod(start time.Time, end time.Time) ([]*storage.Event, error) {
	rows, err := s.db.NamedQuery(
		`select id, title, event_start, event_end, description, owner_id, notify_before
		   from events
		  where event_start <= :end
			and event_end   >= :start`, map[string]interface{}{
			"start": start,
			"end":   end,
		})
	if err != nil {
		return nil, fmt.Errorf("error while getting events for period [%v, %v]: %w", start, end, err)
	}
	defer rows.Close()

	events := make([]*storage.Event, 0, 10)
	for rows.Next() {
		event := &storage.Event{}
		err = rows.StructScan(event)
		if err != nil {
			return nil, fmt.Errorf("error while getting events for period [%v, %v]: %w", start, end, err)
		}
		events = append(events, event)
	}
	return events, nil
}

func (s *Storage) Open(ctx context.Context) (err error) {
	s.db, err = sqlx.Open("pgx", s.dsn)
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	return s.db.PingContext(ctx)
}

func (s *Storage) Update(id string, event *storage.Event) error {
	_, err := s.db.NamedExec(
		`update events
		    set title         = :title,
		        event_start   = :event_start,
		        event_end     = :event_end,
		        description   = :description,
		        owner_id      = :owner_id,
		        notify_before = :notify_before
		  where id = :id`, map[string]interface{}{
			"id":            id,
			"title":         event.Title,
			"event_start":   event.Start,
			"event_end":     event.End,
			"description":   event.Description,
			"owner_id":      event.OwnerID,
			"notify_before": event.NotifyBefore,
		})
	if err != nil {
		return fmt.Errorf("error while updating event [%v]: %w", event, err)
	}
	return nil
}

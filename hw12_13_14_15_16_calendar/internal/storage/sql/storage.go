package sqlstorage

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/dostrovskiy/otus-golang-home-work/hw12_13_14_15_16_calendar/internal/storage"
	_ "github.com/jackc/pgx/v4/stdlib" //nolint
	"github.com/jmoiron/sqlx"
)

var (
	ErrGettingEventByID = func(id string, err error) error {
		return fmt.Errorf("error while getting event by id [%s]: %w", id, err)
	}
	ErrFindingEventsByParams = func(start time.Time, end time.Time, err error) error {
		return fmt.Errorf("error while getting events for period [%v, %v]: %w", start, end, err)
	}
	ErrFindingEventsForNotify = func(notifyDate time.Time, notified bool, err error) error {
		return fmt.Errorf("error while getting events for notify date [%v] and sign [%v]: %w", notifyDate, notified, err)
	}
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

type SQLEvent struct {
	ID           string        `db:"id"`
	Title        string        `db:"title"`
	Start        time.Time     `db:"event_start"`
	End          sql.NullTime  `db:"event_end"`
	Description  string        `db:"description"`
	OwnerID      string        `db:"owner_id"`
	NotifyBefore sql.NullInt64 `db:"notify_before"`
	NotifyStart  sql.NullTime  `db:"notify_start"`
	Notified     sql.NullBool  `db:"notified"`
}

func (s *Storage) Add(event *storage.Event) (*storage.Event, error) {
	_, err := s.db.NamedExec(
		`insert into events 
	       (id, title, event_start, event_end, description, owner_id, notify_before, notify_start, notified) 
	     values (:id, :title, :event_start, :event_end, :description, :owner_id, :notify_before,
		         :notify_start, :notified)`,
		map[string]interface{}{
			"id":            event.ID,
			"title":         event.Title,
			"event_start":   event.Start,
			"event_end":     event.End,
			"description":   event.Description,
			"owner_id":      event.OwnerID,
			"notify_before": event.NotifyBefore,
			"notify_start":  event.NotifyStart,
			"notified":      event.Notified,
		})
	if err != nil {
		return nil, fmt.Errorf("error while adding event [%v]: %w", event, err)
	}
	return event, nil
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
		`select id, title, event_start, event_end, description, owner_id, notify_before,
				notify_start, notified
	       from events
	      where id = :id`, map[string]interface{}{"id": id})
	if err != nil {
		return nil, ErrGettingEventByID(id, err)
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}
	sqlEvent := SQLEvent{}
	err = rows.StructScan(&sqlEvent)
	if err != nil {
		return nil, ErrGettingEventByID(id, err)
	}
	return mapSQLEventToStorageEvent(sqlEvent), nil
}

func (s *Storage) FindByPeriod(start time.Time, end time.Time) ([]*storage.Event, error) {
	rows, err := s.db.NamedQuery(
		`select id, title, event_start, event_end, description, owner_id, notify_before,
				notify_start, notified
		   from events
		  where event_start <= :end
			and event_end   >= :start`, map[string]interface{}{
			"start": start,
			"end":   end,
		})
	if err != nil {
		return nil, ErrFindingEventsByParams(start, end, err)
	}
	defer rows.Close()

	events := make([]*storage.Event, 0, 10)
	for rows.Next() {
		sqlEvent := SQLEvent{}
		err = rows.StructScan(&sqlEvent)
		if err != nil {
			return nil, ErrFindingEventsByParams(start, end, err)
		}
		events = append(events, mapSQLEventToStorageEvent(sqlEvent))
	}
	return events, nil
}

func (s *Storage) FindForNotify(notifyDate time.Time, notified bool) ([]*storage.Event, error) {
	rows, err := s.db.NamedQuery(
		`select id, title, event_start, event_end, description, owner_id, notify_before,
				notify_start, notified
		   from events
		  where notify_start  > :zeroDate
		    and notify_start <= :notifyDate
			and event_start  >= :notifyDate
			and notified      = :notified`, map[string]interface{}{
			"zeroDate":   time.Time{},
			"notifyDate": notifyDate,
			"notified":   notified,
		})
	if err != nil {
		return nil, ErrFindingEventsForNotify(notifyDate, notified, err)
	}
	defer rows.Close()

	events := make([]*storage.Event, 0, 10)
	for rows.Next() {
		sqlEvent := SQLEvent{}
		err = rows.StructScan(&sqlEvent)
		if err != nil {
			return nil, ErrFindingEventsForNotify(notifyDate, notified, err)
		}
		events = append(events, mapSQLEventToStorageEvent(sqlEvent))
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

func (s *Storage) Update(id string, event *storage.Event) (*storage.Event, error) {
	event.ID = id
	_, err := s.db.NamedExec(
		`update events
		    set title         = :title,
		        event_start   = :event_start,
		        event_end     = :event_end,
		        description   = :description,
		        owner_id      = :owner_id,
		        notify_before = :notify_before,
		        notify_start  = :notify_start,
				notified      = :notified
		  where id = :id`, map[string]interface{}{
			"id":            event.ID,
			"title":         event.Title,
			"event_start":   event.Start,
			"event_end":     event.End,
			"description":   event.Description,
			"owner_id":      event.OwnerID,
			"notify_before": event.NotifyBefore,
			"notify_start":  event.NotifyStart,
			"notified":      event.Notified,
		})
	if err != nil {
		return nil, fmt.Errorf("error while updating event [%+v]: %w", event, err)
	}
	return event, nil
}

func (s *Storage) AddNotification(notification *storage.Notification) (*storage.Notification, error) {
	_, err := s.db.NamedExec(
		`insert into notifications 
	       (id, event_id, title, event_start, event_end, owner_id) 
	     values (:id, :event_id, :title, :event_start, :event_end, :owner_id)`,
		map[string]interface{}{
			"id":          notification.ID,
			"event_id":    notification.EventID,
			"title":       notification.Title,
			"event_start": notification.EventStart,
			"event_end":   notification.EventEnd,
			"owner_id":    notification.OwnerID,
		})
	if err != nil {
		return nil, fmt.Errorf("error while adding notification [%+v]: %w", notification, err)
	}
	return notification, nil
}

func mapSQLEventToStorageEvent(e SQLEvent) *storage.Event {
	return &storage.Event{
		ID:           e.ID,
		Title:        e.Title,
		Start:        e.Start,
		End:          e.End.Time,
		Description:  e.Description,
		OwnerID:      e.OwnerID,
		NotifyBefore: time.Duration(e.NotifyBefore.Int64),
		NotifyStart:  e.NotifyStart.Time,
		Notified:     e.Notified.Bool,
	}
}

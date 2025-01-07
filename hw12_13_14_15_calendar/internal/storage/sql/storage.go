package sqlstorage

import (
	"context"
	"fmt"
	"time"

	"github.com/dostrovskiy/otus-golang-home-work/hw12_13_14_15_calendar/internal/storage"
	_ "github.com/jackc/pgx/v4/stdlib" //nolint
	"github.com/jmoiron/sqlx"
)

type Storage struct {
	db  *sqlx.DB
	dsn string // "postgres://user:password@localhost:5432/dbname?sslmode=disable"
}

// Delete implements storage.EventStorage.
func (s *Storage) Delete(id string) error {
	_ = id
	return fmt.Errorf("unimplemented")
}

// Get implements storage.EventStorage.
func (s *Storage) Get(id string) (*storage.Event, error) {
	_ = id
	return nil, fmt.Errorf("unimplemented")
}

// GetForPeriod implements storage.EventStorage.
func (s *Storage) GetForPeriod(start time.Time, end time.Time) ([]storage.Event, error) {
	_ = start
	_ = end
	return nil, fmt.Errorf("unimplemented")
}

// Update implements storage.EventStorage.
func (s *Storage) Update(event storage.Event) error {
	_ = event
	return fmt.Errorf("unimplemented")
}

func New(dsn string) *Storage {
	return &Storage{
		dsn: dsn,
	}
}

func (s *Storage) Open(ctx context.Context) (err error) {
	s.db, err = sqlx.Open("pgx", s.dsn)
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	return s.db.PingContext(ctx)
}

func (s *Storage) Close() error {
	return s.db.Close()
}

func (s *Storage) Add(event storage.Event) error {
	_ = event
	return fmt.Errorf("unimplemented")
}

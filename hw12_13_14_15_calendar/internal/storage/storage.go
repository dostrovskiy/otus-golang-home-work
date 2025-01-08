package storage //nolint

import (
	"fmt"
	"time"
)

var ErrEventNotFoundByID = func(id string) error { return fmt.Errorf("event not found by id: %s", id) }
var ErrEventAlreadyExists = func(id string) error { return fmt.Errorf("event with id [%s] already exists", id) }
var ErrEventNotExists = func(id string) error { return fmt.Errorf("event with id [%s] not exists", id) }

type EventStorage interface {
	Add(event Event) error
	Get(id string) (*Event, error)
	GetForPeriod(start, end time.Time) ([]Event, error)
	Update(event Event) error
	Delete(id string) error
	Close() error
}
//nolint
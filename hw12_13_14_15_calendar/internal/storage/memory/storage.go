package memorystorage

import (
	"sync"
	"time"

	"github.com/dostrovskiy/otus-golang-home-work/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct {
	mu     sync.RWMutex
	events map[string]storage.Event
}

// Connect implements storage.EventStorage.
func (s *Storage) Close() error {
	return nil // no close for memory storage.
}

func New() *Storage {
	return &Storage{
		events: make(map[string]storage.Event),
	}
}

func (s *Storage) Add(event storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.events[event.ID]; ok {
		return storage.ErrEventAlreadyExists(event.ID)
	}
	s.events[event.ID] = event
	return nil
}

func (s *Storage) Get(id string) (*storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	event, ok := s.events[id]
	if !ok {
		return nil, storage.ErrEventNotFoundByID(id)
	}
	return &event, nil
}

func (s *Storage) GetForPeriod(start time.Time, end time.Time) ([]storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	events := make([]storage.Event, 0, len(s.events))
	for _, event := range s.events {
		if event.Start.Before(end) && event.End.After(start) {
			events = append(events, event)
		}
	}
	return events, nil
}

func (s *Storage) Update(event storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.events[event.ID]; !ok {
		return storage.ErrEventNotExists(event.ID)
	}
	s.events[event.ID] = event
	return nil
}

func (s *Storage) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.events, id)
	return nil
}

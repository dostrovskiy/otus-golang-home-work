package memorystorage

import (
	"fmt"
	"sync"
	"time"

	"github.com/dostrovskiy/otus-golang-home-work/hw12_13_14_15_16_calendar/internal/storage"
)

type Storage struct {
	mu     sync.RWMutex
	events map[string]*storage.Event
}

func (s *Storage) Close() error {
	return nil // nothing to close for memory storage.
}

func New() *Storage {
	return &Storage{
		events: make(map[string]*storage.Event),
	}
}

func (s *Storage) Add(event *storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.events[event.ID]; ok {
		return fmt.Errorf("event with id [%s] already exists", event.ID)
	}
	s.events[event.ID] = event
	return nil
}

func (s *Storage) Get(id string) (*storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	event, ok := s.events[id]
	if !ok {
		return nil, fmt.Errorf("event not found by id: %s", id)
	}
	return event, nil
}

func (s *Storage) GetForPeriod(start time.Time, end time.Time) ([]*storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	events := make([]*storage.Event, 0, len(s.events))
	for _, event := range s.events {
		if event.Start.Before(end) && event.End.After(start) {
			events = append(events, event)
		}
	}
	return events, nil
}

func (s *Storage) Update(id string, event *storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.events[id]; !ok {
		return fmt.Errorf("event with id [%s] not exists", id)
	}
	s.events[id] = event
	return nil
}

func (s *Storage) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.events, id)
	return nil
}

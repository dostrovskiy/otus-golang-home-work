package storage

import "time"

type Event struct {
	ID           string
	Title        string
	Start        time.Time
	End          time.Time
	Description  string
	OwnerID      string
	NotifyBefore time.Duration
}

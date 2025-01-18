package storage

import "time"

type Event struct {
	ID           string        `db:"id"`
	Title        string        `db:"title"`
	Start        time.Time     `db:"event_start"`
	End          time.Time     `db:"event_end"`
	Description  string        `db:"description"`
	OwnerID      string        `db:"owner_id"`
	NotifyBefore time.Duration `db:"notify_before"`
}

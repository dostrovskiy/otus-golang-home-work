package storage //nolint

import (
	"time"
)

type Notification struct {
	ID         string    `db:"id" json:"id"`
	EventID    string    `db:"event_id" json:"eventId"`
	Title      string    `db:"title" json:"title"`
	EventStart time.Time `db:"event_start" json:"eventStart"`
	EventEnd   time.Time `db:"event_end" json:"eventEnd"`
	OwnerID    string    `db:"owner_id" json:"ownerId"`
}
//nolint
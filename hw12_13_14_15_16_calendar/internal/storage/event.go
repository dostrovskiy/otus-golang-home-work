package storage

import (
	"time"
)

type Event struct {
	ID           string        `db:"id" json:"id"`
	Title        string        `db:"title" json:"title"`
	Start        time.Time     `db:"event_start" json:"eventStart"`
	End          time.Time     `db:"event_end" json:"eventEnd"`
	Description  string        `db:"description" json:"description"`
	OwnerID      string        `db:"owner_id" json:"ownerId"`
	NotifyBefore time.Duration `db:"notify_before" json:"notifyBefore"`
	NotifyStart  time.Time     `db:"notify_start" json:"notifyStart"`
	Notified     bool          `db:"notified" json:"notified"`
}

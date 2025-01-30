package integrationtests //nolint

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v4/stdlib" //nolint
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

type Notification struct {
	ID         string    `db:"id"`
	EventID    string    `db:"event_id"`
	Title      string    `db:"title"`
	EventStart time.Time `db:"event_start"`
	EventEnd   time.Time `db:"event_end"`
	OwnerID    string    `db:"owner_id"`
}

func TestNotification(t *testing.T) {
	id := uuid.New().String()
	event := Event{
		ID:           id,
		Title:        "Meeting for notification",
		Start:        time.Now().Add(time.Hour),
		End:          time.Now().Add(2 * time.Hour),
		OwnerID:      "Team Lead",
		NotifyBefore: 2 * time.Hour,
		Description:  "Meeting with team",
	}

	t.Run("post event and get notification", func(t *testing.T) {
		cli := &http.Client{}

		jsonData, err := json.Marshal(event)
		assert.NoError(t, err)

		postURL, err := url.JoinPath(config.Server.Address, "event")
		assert.NoError(t, err)

		req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, postURL, bytes.NewReader(jsonData))
		assert.NoError(t, err)
		postResp, err := cli.Do(req)
		assert.NoError(t, err)
		defer func() { assert.NoError(t, postResp.Body.Close()) }()
		assert.Equal(t, http.StatusCreated, postResp.StatusCode)

		db, err := sqlx.Open("pgx", config.DataSource.Dsn)
		assert.NoError(t, err)
		defer func() { assert.NoError(t, db.Close()) }()

		var newEvent Event
		err = db.Get(&newEvent, "SELECT * FROM events WHERE id = $1", id)
		assert.NoError(t, err)

		// wait for notification
		time.Sleep(10 * time.Second)

		var notification Notification
		err = db.Get(&notification, "SELECT * FROM notifications WHERE event_id = $1", id)
		assert.NoError(t, err)

		// clear events and notifications
		_, err = db.Exec("DELETE FROM events WHERE id = $1", id)
		assert.NoError(t, err)
		_, err = db.Exec("DELETE FROM notifications WHERE event_id = $1", id)
		assert.NoError(t, err)
	})
}
//nolint
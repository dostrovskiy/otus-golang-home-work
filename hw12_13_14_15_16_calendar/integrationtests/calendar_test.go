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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type Event struct {
	ID           string        `json:"id" db:"id"`
	Title        string        `json:"title" db:"title"`
	Start        time.Time     `json:"start" db:"event_start"`
	End          time.Time     `json:"end" db:"event_end"`
	Description  string        `json:"description" db:"description"`
	OwnerID      string        `json:"ownerId" db:"owner_id"`
	NotifyBefore time.Duration `json:"notifyBefore" db:"notify_before"`
	NotifyStart  time.Time     `json:"notifyStart" db:"notify_start"`
	Notified     bool          `json:"notified" db:"notified"`
}

func TestEvent(t *testing.T) {
	t.Run("post and get event", func(t *testing.T) {
		cli := &http.Client{}

		newEvent := genTestEvent(time.Now())
		id := newEvent.ID

		jsonData, err := json.Marshal(newEvent)
		require.NoError(t, err)

		eventURL, err := url.JoinPath(config.Server.Address, "event")
		require.NoError(t, err)

		postReq, err := http.NewRequestWithContext(context.Background(), http.MethodPost, eventURL, bytes.NewReader(jsonData))
		require.NoError(t, err)
		postResp, err := cli.Do(postReq)
		require.NoError(t, err)
		require.NoError(t, postResp.Body.Close())
		require.Equal(t, http.StatusCreated, postResp.StatusCode)

		getURL, err := url.JoinPath(eventURL, id)
		require.NoError(t, err)

		getReq, err := http.NewRequestWithContext(context.Background(), http.MethodGet, getURL, nil)
		require.NoError(t, err)
		getResp, err := cli.Do(getReq)
		require.NoError(t, err)
		require.NoError(t, getResp.Body.Close())
		require.Equal(t, http.StatusOK, getResp.StatusCode)

		delURL, err := url.JoinPath(eventURL, id)
		require.NoError(t, err)
		delReq, err := http.NewRequestWithContext(context.Background(), http.MethodDelete, delURL, nil)
		require.NoError(t, err)
		delResp, err := cli.Do(delReq)
		require.NoError(t, err)
		require.NoError(t, delResp.Body.Close())
		require.Equal(t, http.StatusNoContent, delResp.StatusCode)
	})

	t.Run("find event for period", func(t *testing.T) {
		events := make([]Event, 0, 10)
		// this day event
		events = append(events, genTestEvent(time.Now()))
		// this week, but not this day event
		if time.Now().Weekday() >= time.Wednesday {
			events = append(events, genTestEvent(time.Now().AddDate(0, 0, -2)))
		} else {
			events = append(events, genTestEvent(time.Now().AddDate(0, 0, 2)))
		}
		// this month, but not this week event
		if time.Now().Day() >= 15 {
			events = append(events, genTestEvent(time.Now().AddDate(0, 0, -8)))
		} else {
			events = append(events, genTestEvent(time.Now().AddDate(0, 0, 8)))
		}

		for _, event := range events {
			postEvent(t, event)
		}

		thisDayEvents := findEvents(t, BeginOfToday(), EndOfToday())
		assert.Equal(t, 1, len(thisDayEvents))

		thisWeekEvents := findEvents(t, BeginningOfThisWeek(), EndingOfThisWeek())
		assert.Equal(t, 2, len(thisWeekEvents))

		thisMonthEvents := findEvents(t, BeginningOfThisMonth(), EndingOfThisMonth())
		require.Equal(t, 3, len(thisMonthEvents))

		for _, event := range events {
			deleteEvent(t, event.ID)
		}
	})
}

func BeginOfToday() time.Time {
	y, m, d := time.Now().Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
}

func EndOfToday() time.Time {
	y, m, d := time.Now().Date()
	return time.Date(y, m, d+1, 0, 0, 0, 0, time.UTC)
}

func BeginningOfThisWeek() time.Time {
	y, m, d := time.Now().Date()
	return time.Date(y, m, d-int(time.Now().Weekday()), 0, 0, 0, 0, time.UTC)
}

func EndingOfThisWeek() time.Time {
	y, m, d := time.Now().Date()
	return time.Date(y, m, d-int(time.Now().Weekday())+7, 23, 59, 59, 999999999, time.UTC)
}

func BeginningOfThisMonth() time.Time {
	y, m, _ := time.Now().Date()
	return time.Date(y, m, 1, 0, 0, 0, 0, time.UTC)
}

func EndingOfThisMonth() time.Time {
	y, m, _ := time.Now().Date()
	return time.Date(y, m+1, 1, 23, 59, 59, 999999999, time.UTC).AddDate(0, 0, -1)
}

func genTestEvent(start time.Time) Event {
	return Event{
		ID:           uuid.New().String(),
		Title:        "Meeting",
		Start:        start,
		End:          start.Add(time.Hour),
		OwnerID:      "Team Lead",
		NotifyBefore: time.Hour,
		Description:  "Meeting with team",
	}
}

func findEvents(t *testing.T, from, to time.Time) []Event {
	t.Helper()
	cli := &http.Client{}
	pathURL, err := url.JoinPath(config.Server.Address, "events", "by-period")
	require.NoError(t, err)

	params := url.Values{
		"start": {from.Format("2006-01-02T00:00:00Z")},
		"end":   {to.Format("2006-01-02T00:00:00Z")},
	}
	findURL := pathURL + "?" + params.Encode()
	findReq, err := http.NewRequestWithContext(context.Background(), http.MethodGet, findURL, nil)
	require.NoError(t, err)

	findResp, err := cli.Do(findReq)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, findResp.StatusCode)

	var foundEvents []Event
	err = json.NewDecoder(findResp.Body).Decode(&foundEvents)
	require.NoError(t, err)
	defer func() { require.NoError(t, findResp.Body.Close()) }()

	return foundEvents
}

func postEvent(t *testing.T, event Event) {
	t.Helper()
	cli := &http.Client{}
	eventURL, err := url.JoinPath(config.Server.Address, "event")
	require.NoError(t, err)

	jsonData, err := json.Marshal(event)
	require.NoError(t, err)

	postReq, err := http.NewRequestWithContext(context.Background(), http.MethodPost, eventURL, bytes.NewReader(jsonData))
	require.NoError(t, err)

	postResp, err := cli.Do(postReq)
	require.NoError(t, err)
	defer func() { require.NoError(t, postResp.Body.Close()) }()
	require.Equal(t, http.StatusCreated, postResp.StatusCode)
}

func deleteEvent(t *testing.T, id string) {
	t.Helper()
	cli := &http.Client{}
	eventURL, err := url.JoinPath(config.Server.Address, "event")
	require.NoError(t, err)

	delURL, err := url.JoinPath(eventURL, id)
	require.NoError(t, err)

	delReq, err := http.NewRequestWithContext(context.Background(), http.MethodDelete, delURL, nil)
	require.NoError(t, err)

	delResp, err := cli.Do(delReq)
	require.NoError(t, err)
	defer func() { require.NoError(t, delResp.Body.Close()) }()
	require.Equal(t, http.StatusNoContent, delResp.StatusCode)
}
//nolint
package internalhttp //nolint

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/oapi-codegen/testutil"
	"github.com/stretchr/testify/require"

	app "github.com/dostrovskiy/otus-golang-home-work/hw12_13_14_15_16_calendar/internal/app"
	"github.com/dostrovskiy/otus-golang-home-work/hw12_13_14_15_16_calendar/internal/logger"
	memorystorage "github.com/dostrovskiy/otus-golang-home-work/hw12_13_14_15_16_calendar/internal/storage/memory"
)

func TestServer(t *testing.T) {
	var err error

	swagger, err := GetSwagger()
	require.NoError(t, err)

	// Clear out the servers array in the swagger spec, that skips validating
	// that server names match. We don't know how this thing will be run.
	swagger.Servers = nil

	log := logger.New("INFO")
	storage := memorystorage.New()
	app := app.New(log, storage)
	server := NewServer(log, app)
	hand := NewStrictHandler(server, nil)
	mux := http.NewServeMux()

	HandlerFromMux(hand, mux)

	t.Run("basic crud", func(t *testing.T) {
		title := "test event"
		start := time.Now()
		end := time.Now().Add(time.Hour)
		event := &Event{
			Title: &title,
			Start: &start,
			End:   &end,
		}

		var addEvent Event
		add := testutil.NewRequest().Post("/event").WithJsonBody(event).GoWithHTTPHandler(t, mux).Recorder
		require.Equal(t, http.StatusCreated, add.Code)
		err = json.NewDecoder(add.Body).Decode(&addEvent)
		require.NoError(t, err, "error unmarshaling post event response")
		require.NotEmpty(t, addEvent.Id)
		var id string
		if addEvent.Id != nil {
			id = *addEvent.Id
		}

		var getEvent Event
		get := testutil.NewRequest().Get(fmt.Sprintf("/event/%s", id)).GoWithHTTPHandler(t, mux).Recorder
		err = json.NewDecoder(get.Body).Decode(&getEvent)
		require.NoError(t, err, "error unmarshaling get event response")
		require.Equal(t, event.Title, getEvent.Title)
		require.True(t, event.Start.Round(0).Equal(getEvent.Start.Round(0)), "expected %v, got %v", event.Start.Round(0), getEvent.Start.Round(0))

		newTitle := "updated event"
		newStart := time.Now().Add(5 * time.Hour)
		newEnd := time.Now().Add(6 * time.Hour)
		newEvent := &Event{
			Title: &newTitle,
			Start: &newStart,
			End:   &newEnd,
		}
		var updEvent Event
		upd := testutil.NewRequest().Put(fmt.Sprintf("/event/%s", id)).WithJsonBody(newEvent).GoWithHTTPHandler(t, mux).Recorder //nolint:lll
		require.Equal(t, http.StatusOK, upd.Code)
		err = json.NewDecoder(upd.Body).Decode(&updEvent)
		require.NoError(t, err, "error unmarshaling update event response")
		require.Equal(t, newEvent.Title, updEvent.Title)
		require.True(t, newEvent.Start.Round(0).Equal(updEvent.Start.Round(0)), "expected %v, got %v", newEvent.Start.Round(0), updEvent.Start.Round(0))

		del := testutil.NewRequest().Delete(fmt.Sprintf("/event/%s", id)).GoWithHTTPHandler(t, mux).Recorder
		require.Equal(t, http.StatusNoContent, del.Code)

		getAft := testutil.NewRequest().Get(fmt.Sprintf("/event/%s", id)).GoWithHTTPHandler(t, mux).Recorder
		require.Equal(t, http.StatusOK, getAft.Code)
		require.Equal(t, "{}\n", getAft.Body.String())
	})
}
//nolint:gofumpt
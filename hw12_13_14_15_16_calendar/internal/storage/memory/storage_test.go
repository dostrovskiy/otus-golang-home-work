package memorystorage

import (
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/dostrovskiy/otus-golang-home-work/hw12_13_14_15_16_calendar/internal/storage"
	"github.com/stretchr/testify/require"
)

func TestStorage(t *testing.T) {
	store := New()

	event := &storage.Event{
		ID:    "1",
		Title: "Event 1",
		Start: time.Now(),
		End:   time.Now().Add(time.Hour),
	}

	t.Run("basic crud", func(t *testing.T) {
		err := store.Add(event)
		require.NoError(t, err)

		got, err := store.Get(event.ID)
		require.NoError(t, err)
		require.Equal(t, &event, got, "expected %v, got %v", event, got)

		events, err := store.GetForPeriod(time.Now(), time.Now().Add(time.Hour))
		require.NoError(t, err)
		require.Equal(t, []*storage.Event{event}, events, "expected %v, got %v", []*storage.Event{event}, events)

		err = store.Update("1", event)
		require.NoError(t, err)

		got, err = store.Get("2")
		require.Error(t, err)
		require.Nil(t, got)

		err = store.Delete(event.ID)
		require.NoError(t, err)

		got, err = store.Get(event.ID)
		require.Error(t, err)
		require.Nil(t, got)
	})
}

func TestStorageMultithread(t *testing.T) {
	wg := &sync.WaitGroup{}
	eventCount := 1000
	workerCount := 5

	t.Run("multithread adding", func(t *testing.T) {
		store := New()
		cnt := eventCount / workerCount
		wg.Add(workerCount)
		for g := 0; g < workerCount; g++ {
			go func() {
				defer wg.Done()
				for i := g * cnt; i < g*cnt+cnt; i++ {
					event := &storage.Event{
						ID:    strconv.Itoa(i),
						Title: "Event " + strconv.Itoa(i),
						Start: time.Now(),
						End:   time.Now().Add(time.Hour),
					}
					err := store.Add(event)
					require.NoError(t, err)
				}
			}()
		}
		wg.Wait()
		require.Equal(t, eventCount, len(store.events), "expected %d events, got %d", eventCount, len(store.events))
	})

	t.Run("multithread deleting", func(t *testing.T) {
		store := New()
		cnt := eventCount / workerCount
		for i := 0; i < eventCount; i++ {
			event := &storage.Event{
				ID:    strconv.Itoa(i),
				Title: "Event " + strconv.Itoa(i),
				Start: time.Now(),
				End:   time.Now().Add(time.Hour),
			}
			err := store.Add(event)
			require.NoError(t, err)
		}
		wg.Add(workerCount)
		for g := 0; g < workerCount; g++ {
			go func() {
				defer wg.Done()
				for i := g * cnt; i < g*cnt+cnt; i++ {
					err := store.Delete(strconv.Itoa(i))
					require.NoError(t, err)
				}
			}()
		}
		wg.Wait()
		require.Empty(t, store.events)
	})
}

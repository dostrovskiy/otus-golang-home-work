package hw05parallelexecution

import (
	"errors"
	"fmt"
	"math/rand"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

func TestRun(t *testing.T) {
	defer goleak.VerifyNone(t)

	t.Run("if were errors in first M tasks, than finished not more N+M tasks", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32

		for i := 0; i < tasksCount; i++ {
			err := fmt.Errorf("error from task %d", i)
			tasks = append(tasks, func() error {
				<-time.After(time.Millisecond * time.Duration(rand.Intn(100)))
				atomic.AddInt32(&runTasksCount, 1)
				return err
			})
		}

		workersCount := 10
		maxErrorsCount := 23
		err := Run(tasks, workersCount, maxErrorsCount)

		require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)
		require.LessOrEqual(t, runTasksCount, int32(workersCount+maxErrorsCount), "extra tasks were started")
	})

	t.Run("tasks without errors", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32
		var sumTime time.Duration

		for i := 0; i < tasksCount; i++ {
			taskSleep := time.Millisecond * time.Duration(rand.Intn(100))
			sumTime += taskSleep

			tasks = append(tasks, func() error {
				<-time.After(taskSleep)
				atomic.AddInt32(&runTasksCount, 1)
				return nil
			})
		}

		workersCount := 5
		maxErrorsCount := 1

		start := time.Now()
		err := Run(tasks, workersCount, maxErrorsCount)
		elapsedTime := time.Since(start)
		require.NoError(t, err)

		require.Equal(t, runTasksCount, int32(tasksCount), "not all tasks were completed")
		require.LessOrEqual(t, int64(elapsedTime), int64(sumTime/2), "tasks were run sequentially?")
	})
}

func TestRunNoErrors(t *testing.T) {
	defer goleak.VerifyNone(t)

	t.Run("run with no errors allowed", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)
		var runTasksCount int32
		var e error
		for i := 0; i < tasksCount; i++ {
			if i == tasksCount/2 {
				e = fmt.Errorf("error from task %d", i)
			}
			taskSleep := time.Millisecond * time.Duration(rand.Intn(100))
			tasks = append(tasks, func() error {
				<-time.After(taskSleep)
				atomic.AddInt32(&runTasksCount, 1)
				return e
			})
		}

		workersCount := 10
		err := Run(tasks, workersCount, 0)

		require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)
		require.LessOrEqual(t, runTasksCount, int32(workersCount+1), "extra tasks were started")
	})

	t.Run("ignore errors", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)
		var runTasksCount int32

		for i := 0; i < tasksCount; i++ {
			err := fmt.Errorf("error from task %d", i)
			taskSleep := time.Millisecond * time.Duration(rand.Intn(100))
			tasks = append(tasks, func() error {
				<-time.After(taskSleep)
				atomic.AddInt32(&runTasksCount, 1)
				return err
			})
		}

		err := Run(tasks, 10, -1)

		require.NoError(t, err)
		require.Equal(t, runTasksCount, int32(tasksCount), "not all tasks were completed")
	})
}

func TestRunErrors(t *testing.T) {
	t.Run("no goroutines to run task", func(t *testing.T) {
		tasks := []Task{}
		err := Run(tasks, 0, 0)

		require.Truef(t, errors.Is(err, ErrNoGoroutinesToRunTask), "actual err - %v", err)
	})
}

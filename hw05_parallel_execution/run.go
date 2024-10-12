/*
Package hw05parallelexecution is an implementation of parallel execution of tasks.
*/
package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

// ErrErrorsLimitExceeded is returned when errors limit exceeded.
var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

// ErrNoGoroutinesToRunTask is returned when no goroutines to run task.
var ErrNoGoroutinesToRunTask = errors.New("no goroutines to run task")

// Task is a function that returns error.
type Task func() error

func maxErrExeeded(errCount *atomic.Int64, maxErrCount int) bool {
	if e := errCount.Load(); maxErrCount >= 0 && e > 0 && e >= int64(maxErrCount) {
		return true
	}
	return false
}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
// if m == 0, then no errors are allowed.
// if m < 0, then ignore all errors.
func Run(tasks []Task, n, m int) error {
	if n < 1 {
		return ErrNoGoroutinesToRunTask
	}
	var errCount atomic.Int64
	var wg sync.WaitGroup
	wg.Add(n)
	taskChan := make(chan Task, n)
	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			for {
				task, open := <-taskChan
				if !open || maxErrExeeded(&errCount, m) {
					break
				}
				if err := task(); err != nil {
					errCount.Add(1)
				}
			}
		}()
	}
	var err error
	for _, task := range tasks {
		if maxErrExeeded(&errCount, m) {
			err = ErrErrorsLimitExceeded
			break
		}
		taskChan <- task
	}
	close(taskChan)
	wg.Wait()
	return err
}

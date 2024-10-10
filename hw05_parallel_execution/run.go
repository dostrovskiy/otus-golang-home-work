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

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
// if m == 0, then no errors are allowed.
// if m < 0, then ignore all errors.
func Run(tasks []Task, n, m int) error {
	if n < 1 {
		return ErrNoGoroutinesToRunTask
	}
	ticketChan := make(chan int, n) // Add n-tickets for goroutines to run task.
	for i := 0; i < n; i++ {
		ticketChan <- i
	}
	var ticket int
	var err error
	var e, errCount int64
	var wg sync.WaitGroup
	for _, task := range tasks {
		e = atomic.LoadInt64(&errCount)
		if m >= 0 && e > 0 && e >= int64(m) {
			err = ErrErrorsLimitExceeded
			break
		}
		ticket = <-ticketChan // Take a ticket for starting job, wait if all tikets are taken.
		wg.Add(1)
		go func(ticket int, task Task) {
			defer func() {
				ticketChan <- ticket // Return ticket after work.
				wg.Done()
			}()
			if err := task(); err != nil {
				atomic.AddInt64(&errCount, 1)
			}
		}(ticket, task)
	}
	wg.Wait()
	return err
}

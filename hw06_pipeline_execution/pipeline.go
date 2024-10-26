// Package hw06pipelineexecution provides a pipeline execution function
package hw06pipelineexecution

type (
	// In is a channel for input data for pipeline stage.
	In = <-chan interface{}
	// Out is a channel for output data for pipeline stage.
	Out = In
	// Bi is a bidirectional channel.
	Bi = chan interface{}
)

// Stage is a pipeline stage function.
type Stage func(in In) (out Out)

// ExecutePipeline starts pipeline executing stages in order.
func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := in
	for _, stage := range stages {
		midIn := out
		midOut := make(Bi)
		allDone := make(Bi)
		// Goroutine implement intermediate channels between stages to be able to interrupt the pipeline execution
		go func(midOut Bi, allDone Bi) {
			defer close(allDone)
			defer close(midOut)
			for v := range midIn {
				select {
				case <-done:
					return
				case midOut <- v:
				}
			}
		}(midOut, allDone)
		// Goroutine to read up the middle channels before they are closed, after done signal
		go func(midOut Bi, allDone Bi) {
			select {
			case <-done:
				for v := range midOut {
					_ = v
				}
				return
			case <-allDone:
				return
			}
		}(midOut, allDone)

		out = stage(midOut)
	}
	return out
}

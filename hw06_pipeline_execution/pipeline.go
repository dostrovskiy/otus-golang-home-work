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

func check(done In, in In) Out {
	out := make(Bi)
	go func() {
		defer close(out)
		for {
			select {
			case <-done:
				return
			case v, ok := <-in:
				if !ok {
					return
				}
				out <- v
			}
		}
	}()
	return out
}

// ExecutePipeline starts pipeline executing stages in order.
func ExecutePipeline(in In, done In, stages ...Stage) Out {
	for _, stage := range stages {
		in = stage(check(done, in))
	}
	return in
}

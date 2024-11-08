// Package hw06pipelineexecution provides a pipeline execution function
package hw06pipelineexecution

type (
	// In is read only channel for input data.
	In = <-chan interface{}
	// Out is read only channel for output results.
	Out = In
	// Bi is a read/write channel.
	Bi = chan interface{}
)

// Stage is a pipeline stage function, which takes data from read only channel "in"
// and return results to read only channel "out".
type Stage func(in In) (out Out)

func check(done In, in In) Out {
	out := make(Bi)
	go func() {
		defer func() {
			close(out)
			for range out {
				_ = out
			}
			for range in {
				_ = in
			}
		}()
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

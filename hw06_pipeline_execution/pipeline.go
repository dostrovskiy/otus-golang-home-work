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
	// Выстраиваем пайплайн: для каждой стадии создаем горутину, которая слушает исходящий канал от предыдущей стадии
	// Входящий канал для первой: входящий in функции ExecutePipeline
	// Исходящий канал для последней, out, возвращаем из функции ExecutePipeline
	out := in
	for _, stage := range stages {
		out = stage(out)
	}
	return out
}

package channels

import "context"

type TaskManager[I any] interface {
	Active() int
	Threshold() int
	Skipped() int
}

type OperationWithResutl[I, O any] func(context.Context, I, chan<- O) error

type Operation[I any] func(context.Context, I) error
type TaskManagerContext[I any] struct {
	Operation Operation[I]
	context.Context
	UnbufferedChannel[I]
}

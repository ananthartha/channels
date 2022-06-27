package channels

import "context"

type TaskManager[I any] interface {
	Size() int
	Active() int
	Threshold() int
	Skipped() int
}

type TaskManagerContext[I, O any] struct {
	Returns   chan<- O
	Operation Operation[I, O]
	context.Context
	Channel[I]
}

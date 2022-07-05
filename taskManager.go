package channels

import "context"

type TaskManager[I any] interface {
	Active() int
	Threshold() int
	Skipped() int
}

type TaskManagerContext[I, O any] struct {
	Returns   chan<- O
	Operation Operation[I, O]
	context.Context
	UnbufferedChannel[I]
	// Replace Close with Release to close the input channel and end the context and release the resources
}

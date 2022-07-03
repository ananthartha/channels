package channels

import (
	"context"
	"errors"
)

// Errors that are used throughout the Tunny API.
var (
	// TODO: Need to figure out a wau for channel closed logic
	ErrPipeHasBeanClosed = errors.New("the pipe has already benn closed")
)

type Operation[I, O any] func(context.Context, I, chan<- O) error

type Pipe[I, O any] struct {
	ctx  *TaskManagerContext[I, O]
	Stop context.CancelFunc
	UnbufferedInChannel[I]
	TaskManager[I]
}

func (p *Pipe[I, O]) Push(task I) {
	p.In() <- task
}

func (p *Pipe[I, O]) InTo(other UnbufferedInChannel[O]) UnbufferedInChannel[O] {
	p.ctx.Returns = other.In()
	return other
}

func NewPipe[I, O any](operation Operation[I, O], configs ...ConfigFn[I, O]) *Pipe[I, O] {
	config := buildConfig(configs...)
	ctx, stop := context.WithCancel(config.Context)

	instance := &Pipe[I, O]{
		Stop:                stop,
		UnbufferedInChannel: config.Channel,
		ctx: &TaskManagerContext[I, O]{
			Returns:           nil, // TODO: Do the black hole here
			Context:           ctx,
			UnbufferedChannel: config.Channel,
			Operation:         operation,
		},
	}

	instance.TaskManager = config.InitTaskManager(instance.ctx)
	return instance
}

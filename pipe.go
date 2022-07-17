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

type Pipe[I, O any] struct {
	ctx     *TaskManagerContext[I]
	Stop    context.CancelFunc
	Returns chan<- O
	UnbufferedInChannel[I]
	TaskManager[I]
}

func (p *Pipe[I, O]) Push(task I) {
	p.In() <- task
}

func (p *Pipe[I, O]) InTo(other UnbufferedInChannel[O]) UnbufferedInChannel[O] {
	p.Returns = other.In()
	return other
}

func NewPipe[I, O any](operation OperationWithResutl[I, O], configs ...ConfigFn[I]) *Pipe[I, O] {
	config := buildConfig(configs...)
	ctx, stop := context.WithCancel(config.context)

	instance := &Pipe[I, O]{
		Stop:                stop,
		UnbufferedInChannel: config.channel,
		// func(ctx, I) { operation(ctx, I, instance.Returns) }
	}

	instance.ctx = &TaskManagerContext[I]{
		Context:           ctx,
		UnbufferedChannel: config.channel,
		Operation:         func(ctx context.Context, i I) error { return operation(ctx, i, instance.Returns) },
	}

	instance.TaskManager = config.initTaskManager(instance.ctx)
	return instance
}

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
	ctx   TaskManagerContext[I, O]
	Close context.CancelFunc
	InChannel[I]
	TaskManager[I]
}

func (p *Pipe[I, O]) Push(task I) {
	p.In() <- task
}

func (p *Pipe[I, O]) InTo(other *Pipe[O, any]) *Pipe[O, any] {
	p.ctx.Returns = other.In()
	return other
}

func NewPipe[I, O any](operation Operation[I, O], configs ...ConfigFn[I, O]) *Pipe[I, O] {
	config := buildConfig(configs...)
	ctx, close := context.WithCancel(config.Context)

	instance := &Pipe[I, O]{
		Close:     close,
		InChannel: config.Channel,
	}

	instance.TaskManager = config.InitTaskManager(&TaskManagerContext[I, O]{
		Returns:   nil, // TODO: Do the black hole here
		Context:   ctx,
		Channel:   config.Channel,
		Operation: operation,
	})
	return instance
}

func NewPlug[I any](operation Operation[I, interface{}], configs ...ConfigFn[I, interface{}]) *Pipe[I, interface{}] {
	return NewPipe(operation, configs...)
}

package channels

import "context"

type ConsumerOperation[I any] func(context.Context, I) error

type Plug[I any] struct {
	ctx  *TaskManagerContext[I, any]
	Stop context.CancelFunc
	UnbufferedInChannel[I]
	TaskManager[I]
}

func NewPlug[I any](operation ConsumerOperation[I], configs ...ConfigFn[I, any]) *Plug[I] {
	config := buildConfig(configs...)
	ctx, stop := context.WithCancel(config.context)

	instance := &Plug[I]{
		Stop:                stop,
		UnbufferedInChannel: config.channel,
		ctx: &TaskManagerContext[I, any]{
			Returns:           nil, // No need for return as this is a plug
			Context:           ctx,
			UnbufferedChannel: config.channel,
			Operation:         func(ctx context.Context, i I, c chan<- any) error { return operation(ctx, i) },
		},
	}

	instance.TaskManager = config.initTaskManager(instance.ctx)
	return instance
}

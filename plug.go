package channels

import "context"

type ConsumerOperation[I any] func(context.Context, I) error

type Plug[I any] struct {
	ctx   TaskManagerContext[I, any]
	Close context.CancelFunc
	UnbufferedInChannel[I]
	TaskManager[I]
}

func NewPlug[I any](operation ConsumerOperation[I], configs ...ConfigFn[I, any]) *Plug[I] {
	config := buildConfig(configs...)
	ctx, close := context.WithCancel(config.Context)

	instance := &Plug[I]{
		Close:               close,
		UnbufferedInChannel: config.Channel,
	}

	instance.TaskManager = config.InitTaskManager(&TaskManagerContext[I, any]{
		Returns:           nil, // No need for return as this is a plug
		Context:           ctx,
		UnbufferedChannel: config.Channel,
		Operation:         func(ctx context.Context, i I, c chan<- any) error { return operation(ctx, i) },
	})
	return instance
}

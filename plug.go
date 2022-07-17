package channels

import "context"

type Plug[I any] struct {
	ctx  *TaskManagerContext[I]
	Stop context.CancelFunc
	UnbufferedInChannel[I]
	TaskManager[I]
}

func NewPlug[I any](operation Operation[I], configs ...ConfigFn[I]) *Plug[I] {
	config := buildConfig(configs...)
	ctx, stop := context.WithCancel(config.context)

	instance := &Plug[I]{
		Stop:                stop,
		UnbufferedInChannel: config.channel,
		ctx: &TaskManagerContext[I]{
			Context:           ctx,
			UnbufferedChannel: config.channel,
			Operation:         operation,
		},
	}

	instance.TaskManager = config.initTaskManager(instance.ctx)
	return instance
}

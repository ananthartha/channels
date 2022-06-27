package channels

type monoTaskManager[I, O any] struct {
	operation Operation[I, O]
	*TaskManagerContext[I, O]
}

func (tm *monoTaskManager[I, O]) Size() int      { return 1 }
func (tm *monoTaskManager[I, O]) Threshold() int { return 1 }
func (tm *monoTaskManager[I, O]) Active() int    { return 1 }
func (tm *monoTaskManager[I, O]) Skipped() int   { return -1 }

func NewmonoTaskManager[I, O any](ctx *TaskManagerContext[I, O]) *monoTaskManager[I, O] {
	instance := &monoTaskManager[I, O]{
		operation:          ctx.Operation,
		TaskManagerContext: ctx,
	}

	go instance.init()
	return instance
}

func (tm *monoTaskManager[I, O]) init() {
	// When the for loop ends close the channel to clean thing up
	defer tm.Close()
	for {
		select {
		case <-tm.Done():
			// TODO: Log error if any from context
			return
		case request, open := <-tm.Out():
			if !open {
				return
			}
			// The the avilable process from Chanel
			tm.operation(tm.Context, request, tm.Returns)
		}
	}
}

func WithMonoTaskManager[I, O any]() ConfigFn[I, O] {
	return func(c *Config[I, O]) error {
		c.InitTaskManager = func(tmc *TaskManagerContext[I, O]) TaskManager[I] {
			return NewmonoTaskManager(tmc)
		}
		return nil
	}
}

package channels

type monoTaskManager[I any] struct {
	operation Operation[I]
	*TaskManagerContext[I]
}

func (tm *monoTaskManager[I]) Size() int      { return 1 }
func (tm *monoTaskManager[I]) Threshold() int { return 1 }
func (tm *monoTaskManager[I]) Active() int    { return 1 }
func (tm *monoTaskManager[I]) Skipped() int   { return -1 }

func NewmonoTaskManager[I any](ctx *TaskManagerContext[I]) *monoTaskManager[I] {
	instance := &monoTaskManager[I]{
		operation:          ctx.Operation,
		TaskManagerContext: ctx,
	}

	go instance.init()
	return instance
}

func (tm *monoTaskManager[I]) init() {
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
			tm.operation(tm.Context, request)
		}
	}
}

func WithMonoTaskManager[I any]() ConfigFn[I] {
	return func(c *config[I]) error {
		c.initTaskManager = func(tmc *TaskManagerContext[I]) TaskManager[I] {
			return NewmonoTaskManager(tmc)
		}
		return nil
	}
}

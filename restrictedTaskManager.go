package channels

import (
	"fmt"
)

/*
type Process[I any] interface {
	Id() string
	In() chan<- I
	Terminate() // To Clean it self up after the current operation =
}
*/

type restrictedProcess[I, O any] struct {
	id   string
	task chan I
	*restrictedTaskManager[I, O]
}

func (p *restrictedProcess[I, O]) Id() string {
	return p.id
}

func (p *restrictedProcess[I, O]) run() {
	operation := p.operation
	// processes chan process[I]
	for {
		// Block untill we can push into the worker queue
		p.processes <- p
		// Wait for input to be recived
		// This a blocking call so select will only be executed after this blocking call
		task, open := <-p.task
		if !open {
			return
		}
		operation(p.ctx.Context, task, nil)
	}
}

type restrictedTaskManager[I, O any] struct {
	ctx       *TaskManagerContext[I, O]
	size      int
	processes chan *restrictedProcess[I, O]
	operation Operation[I, O]
}

func (tm *restrictedTaskManager[I, O]) Size() int      { return tm.size }
func (tm *restrictedTaskManager[I, O]) Threshold() int { return tm.size }
func (tm *restrictedTaskManager[I, O]) Active() int    { return -1 }
func (tm *restrictedTaskManager[I, O]) Skipped() int   { return -1 }

func NewRestrictedTaskManager[I, O any](ctx *TaskManagerContext[I, O], size int) *restrictedTaskManager[I, O] {
	instance := &restrictedTaskManager[I, O]{
		ctx:       ctx,
		size:      size,
		processes: make(chan *restrictedProcess[I, O]),
		operation: ctx.Operation,
	}

	// TODO: Create the processes and then go for the loop
	for count := 0; count < size; count++ {
		process := &restrictedProcess[I, O]{
			id:                    fmt.Sprintf("restricted-%d", count),
			restrictedTaskManager: instance,
			task:                  make(chan I),
		}
		go process.run()
	}

	go instance.init()
	return instance
}

func (tm *restrictedTaskManager[I, O]) init() {
	// When the for loop ends close the channel to clean thing up
	defer tm.ctx.Close()
	for {
		select {
		case <-tm.ctx.Done():
			// TODO: Log error if any from context
			return
		case request, open := <-tm.ctx.Out():
			if !open {
				return
			}
			// The the avilable process from Chanel
			process := <-tm.processes
			// Send the request to the process
			process.task <- request
		}
	}
}

func WithRestrictedTaskManager[I, O any](size int) ConfigFn[I, O] {
	return func(c *Config[I, O]) error {
		c.InitTaskManager = func(tmc *TaskManagerContext[I, O]) TaskManager[I] {
			return NewRestrictedTaskManager(tmc, size)
		}
		return nil
	}
}

package channels

import (
	"fmt"
	"sync/atomic"
)

/*
type Process[I any] interface {
	Id() string
	In() chan<- I
	Terminate() // To Clean it self up after the current operation =
}
*/

type restrictedProcess[I any] struct {
	id   string
	task chan I
	*restrictedTaskManager[I]
}

func (p *restrictedProcess[I]) Id() string {
	return p.id
}

func (p *restrictedProcess[I]) run() {
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
		operation(p.ctx.Context, task)
		atomic.AddInt32(&p.active, ^int32(0))
	}
}

type restrictedTaskManager[I any] struct {
	ctx       *TaskManagerContext[I]
	size      int
	active    int32
	processes chan *restrictedProcess[I]
	operation Operation[I]
}

func (tm *restrictedTaskManager[I]) Threshold() int { return tm.size }
func (tm *restrictedTaskManager[I]) Active() int    { return int(atomic.LoadInt32(&tm.active)) }
func (tm *restrictedTaskManager[I]) Skipped() int   { return -1 }

func NewRestrictedTaskManager[I any](ctx *TaskManagerContext[I], size int) *restrictedTaskManager[I] {
	instance := &restrictedTaskManager[I]{
		ctx:       ctx,
		size:      size,
		processes: make(chan *restrictedProcess[I], size),
		operation: ctx.Operation,
	}

	// TODO: Create the processes and then go for the loop
	for count := 0; count < size; count++ {
		process := &restrictedProcess[I]{
			id:                    fmt.Sprintf("restricted-%d", count),
			restrictedTaskManager: instance,
			task:                  make(chan I),
		}
		go process.run()
	}

	go instance.init()
	return instance
}

func (tm *restrictedTaskManager[I]) init() {
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
			atomic.AddInt32(&tm.active, 1)
			process.task <- request
		}
	}
}

func WithRestrictedTaskManager[I any](size int) ConfigFn[I] {
	return func(c *config[I]) error {
		c.initTaskManager = func(tmc *TaskManagerContext[I]) TaskManager[I] {
			return NewRestrictedTaskManager(tmc, size)
		}
		return nil
	}
}

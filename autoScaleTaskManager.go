package channels

import (
	"fmt"
	"sync"
	"sync/atomic"
)

/*
type Process[I any] interface {
	Id() string
	In() chan<- I
	Terminate() // To Clean it self up after the current operation =
}
*/

type autoScaleProcess[I any] struct {
	id   string
	task chan I
	*autoScaleTaskManager[I]
}

func (p *autoScaleProcess[I]) Id() string {
	return p.id
}

func (p *autoScaleProcess[I]) run() {
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
	}
}

type autoScaleTaskManager[I any] struct {
	ctx             *TaskManagerContext[I]
	active          int32
	processes       chan *autoScaleProcess[I]
	minThreshold    int
	threshold       int
	operation       Operation[I]
	reziseLock      sync.Mutex
	resizeStrategy  ResizingStrategy
	cancelProcesses chan *autoScaleProcess[I]
}

func (tm *autoScaleTaskManager[I]) Threshold() int { return tm.threshold }
func (tm *autoScaleTaskManager[I]) Active() int    { return int(atomic.LoadInt32(&tm.active)) }
func (tm *autoScaleTaskManager[I]) Skipped() int   { return -1 }

func NewAutoScaleTaskManager[I any](ctx *TaskManagerContext[I], threshold, minThreshold int) *autoScaleTaskManager[I] {
	instance := &autoScaleTaskManager[I]{
		ctx:             ctx,
		threshold:       threshold,
		minThreshold:    minThreshold,
		processes:       make(chan *autoScaleProcess[I], threshold/2), // try with just size 2
		cancelProcesses: make(chan *autoScaleProcess[I], threshold-minThreshold),
		operation:       ctx.Operation,
		resizeStrategy:  Balanced(),
	}

	for count := 0; count < minThreshold; count++ {
		process := &autoScaleProcess[I]{
			id:                   fmt.Sprintf("autoscale-%d", count),
			autoScaleTaskManager: instance,
			task:                 make(chan I),
		}
		go process.run()
		atomic.AddInt32(&instance.active, 1)
	}

	go instance.init()
	return instance
}

func (tm *autoScaleTaskManager[I]) scale() {
	if tm.reziseLock.TryLock() && tm.resizeStrategy.Resize(tm.Active(), tm.minThreshold, tm.threshold) {
		// We do it a go rouine so that the current taks can
		// be picked up by an process and we can resize in async
		go func() {
			// Can we me move this into its owne chain
			defer tm.reziseLock.Unlock()

			// Processes are buffered into channel (they by thre is reduction in workload)
			// TODO: Determine a way to find if load is low
			if len(tm.processes) == cap(tm.processes) && tm.Active() == tm.threshold {
				// if queued then we can kill a few off
				// TODO: Reduce the routines
				// int(float64(tm.Active()-tm.minThreshold)*0.9)
				desired := tm.Active() - int(float64(len(tm.processes)-tm.minThreshold)*0.9)
				for tm.Active() > desired {
					select {
					case process := <-tm.processes:
						// Stop the routine
						close(process.task)
						atomic.AddInt32(&tm.active, ^int32(0))
						tm.cancelProcesses <- process
					default:
					}
				}
				return
				// TODO: Determine a way to find if load is high to increase the rpocesses
			} else if tm.Active() < tm.threshold {
				// loadNow := float64(tm.Active()) / float64(tm.threshold)
				// concurrencyDesired := float64(tm.ctx) * loadAvg / p.targetLoad
				desired := tm.threshold + int(float64(tm.Active()-tm.threshold)*0.9)
				// fmt.Println(desired, tm.threshold, int(float64(tm.Active()-tm.threshold)*0.9))
				for tm.Active() <= desired {
					var process *autoScaleProcess[I]
					select {
					case process = <-tm.cancelProcesses:
						process.task = make(chan I)
					default:
						process = &autoScaleProcess[I]{
							id:                   fmt.Sprintf("autoscale-%d", 10),
							autoScaleTaskManager: tm,
							task:                 make(chan I),
						}
					}
					go process.run()
					atomic.AddInt32(&tm.active, 1)
				}
			}
		}()
	}
}

func (tm *autoScaleTaskManager[I]) init() {
	// When the for loop ends close the channel to clean thing up
	defer tm.ctx.Close()
	for {
		select {
		case <-tm.ctx.Done():
			select {
			// NOTE: Clean Up the active processes and exit
			case process := <-tm.processes:
				close(process.task)
			default:
				return
			}
		case request, open := <-tm.ctx.Out():
			if !open {
				return
			}
			// Re-size the process's as needed
			tm.scale()

			// The the avilable process from Chanel
			process := <-tm.processes
			// Send the request to the process
			process.task <- request
			// fmt.Println("X", tm.Active())
		}
	}
}

func WithAutoScaleTaskManager[I any](minThreshold, threshold int) ConfigFn[I] {
	return func(c *config[I]) error {
		c.initTaskManager = func(tmc *TaskManagerContext[I]) TaskManager[I] {
			return NewAutoScaleTaskManager(tmc, threshold, minThreshold)
		}
		return nil
	}
}

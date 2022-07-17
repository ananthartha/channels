package channels

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"go.uber.org/goleak"
)

func TestPipeWithRestrictedTaskManager(t *testing.T) {
	defer goleak.VerifyNone(t)
	var w sync.WaitGroup
	maxInts := 1000000

	q := NewQueueChannel[int](NewDefaultQueue[int]())

	p := NewPipe(func(ctx context.Context, msg int, out chan<- int) error {
		w.Done()
		return nil
	}, WithRestrictedTaskManager[int](50),
		WithChannel[int, int](q, "binary_packet", nil))

	if q.Len() != 0 {
		t.Error("empty queue length not 0")
	}

	w.Add(maxInts)

	for i := 0; i < maxInts; i++ {
		p.In() <- i
	}

	w.Wait()

	if q.Len() > 0 {
		t.Error("queue is not flushed out", q.Len())
	}

	p.Close()
}

func TestPipeWithRestrictedTaskManager2(t *testing.T) {
	defer goleak.VerifyNone(t)
	var w sync.WaitGroup
	maxInts := 1000000

	q := NewQueueChannel[int](NewDefaultQueue[int]())

	p := NewPipe(func(ctx context.Context, msg int, out chan<- int) error {
		w.Done()
		return nil
	}, WithAutoScaleTaskManager[int](10, 100),
		WithChannel[int, int](q, "binary_packet", nil))

	if q.Len() != 0 {
		t.Error("empty queue length not 0")
	}

	w.Add(maxInts)

	for i := 0; i < maxInts; i++ {
		p.In() <- i
	}

	w.Wait()

	if q.Len() > 0 {
		t.Error("queue is not flushed out", q.Len())
	}

	p.Close()
}

func BenchmarkPipeThoughPutTest(b *testing.B) {
	var w sync.WaitGroup
	fmt.Println(b.N)
	maxInts := b.N

	q := NewQueueChannel[int](NewDefaultQueue[int]())

	p := NewPipe(func(ctx context.Context, msg int, out chan<- int) error {
		w.Done()
		return nil
	}, WithRestrictedTaskManager[int](50),
		WithChannel[int, int](q, "binary_packet", nil))

	if q.Len() != 0 {
		b.Error("empty queue length not 0")
	}

	w.Add(maxInts)

	for i := 0; i < maxInts; i++ {
		p.In() <- i
	}

	w.Wait()

	if q.Len() > 0 {
		b.Error("queue is not flushed out", q.Len())
	}
}

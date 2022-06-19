package channels

type Queue[T any] interface {
	// Length returns the number of elements currently stored in the queue.
	Length() int
	// Peek returns the element at the head of the queue. will retirn nil
	// if the queue is empty.
	Peek() T
	// Poll removes and returns the element from the front of the queue. If the
	// queue is empty, the retirn nil.
	Poll() T
	// Add puts an element on the end of the queue.
	Add(T)
}

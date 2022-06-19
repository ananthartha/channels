package channels

// Native implements the Channel interface by wrapping a native go channel.
type Native[T any] chan T

// NewNative makes a new NativeChannel with the given buffer size. Just a convenience wrapper
// to avoid having to cast the result of make().
func NewNative[T any](size int) Native[T] {
	return make(chan T, size)
}

func (ch Native[T]) In() chan<- T {
	return ch
}

func (ch Native[T]) Out() <-chan T {
	return ch
}

func (ch Native[T]) Len() int {
	return len(ch)
}

func (ch Native[T]) Cap() int {
	return cap(ch)
}

func (ch Native[T]) Close() {
	close(ch)
}

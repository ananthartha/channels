package channels

// import "github.com/eapache/queue"

// QueueChannel implements the Channel interface with an infinite buffer between the input and the output.
type QueueChannel[T any] struct {
	input, output chan T
	length        chan int
	buffer        Queue[T]
	def           T
}

func NewQueueChannel[T any](buffer Queue[T]) *QueueChannel[T] {
	ch := &QueueChannel[T]{
		input:  make(chan T),
		output: make(chan T),
		length: make(chan int),
		buffer: buffer,
	}
	go ch.startBufferChannel()
	return ch
}

func NewQueueChannelWithInputSize[T any](buffer Queue[T], size int) *QueueChannel[T] {
	ch := &QueueChannel[T]{
		input:  make(chan T, size),
		output: make(chan T),
		length: make(chan int),
		buffer: buffer,
	}
	go ch.startBufferChannel()
	return ch
}

func (ch *QueueChannel[T]) In() chan T {
	return ch.input
}

func (ch *QueueChannel[T]) Out() <-chan T {
	return ch.output
}

func (ch *QueueChannel[T]) Len() int {
	return <-ch.length
}

func (ch *QueueChannel[T]) Cap() int {
	return -1
}

func (ch *QueueChannel[T]) Close() {
	close(ch.input)
}

func (ch *QueueChannel[T]) startBufferChannel() {
	var input, output chan T = ch.input, ch.output
	var inIten, nextItem T
	var IsInputOpen bool = true

	// Not sure how this is working
	// Should i check if inout is closed and buffer is nul then close the out channel
	for IsInputOpen || output != nil {
		if ch.buffer.Length() > 0 {
			output = ch.output
			nextItem = ch.buffer.Peek()
		} else {
			output = nil
			nextItem = ch.def
		}

		select {
		case inIten, IsInputOpen = <-input:
			if IsInputOpen {
				ch.buffer.Add(inIten)
			}
		// We have to do it like this, insted of ch.buffer.Poll() as it gets executed always and skips the message
		case output <- nextItem:
			ch.buffer.Poll()
		case ch.length <- ch.buffer.Length():
		}
	}

	close(ch.output)
	close(ch.length)
}

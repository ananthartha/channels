package channels

// import "github.com/eapache/queue"

// QueueChannel implements the Channel interface with an infinite buffer between the input and the output.
type QueueChannel[T any] struct {
	input, output chan T
	length        chan int
	buffer        Queue[T]
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

func (ch *QueueChannel[T]) In() chan<- T {
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

/*
func (ch *QueueChannel[T]) startBufferChannel() {
	var input, output chan *T
	var next *T
	input = ch.input

	// Not sure how this is working
	// Should i check if inout is closed and buffer is nul then close the out channel
	for input != nil || output != nil {
		select {
		case elem, open := <-input:
			if open {
				ch.buffer.Add(elem)
			} else {
				input = nil
			}
		case output <- next:
			ch.buffer.Poll()
		case ch.length <- ch.buffer.Length():
		}

		if ch.buffer.Length() > 0 {
			output = ch.output
			next = ch.buffer.Peek()
		} else {
			output = nil
			next = nil
		}
	}

	close(ch.output)
	close(ch.length)
}
*/
func (ch *QueueChannel[T]) startBufferChannel() {
	var input, output chan T = ch.input, ch.output
	var inIten T
	var IsInputClosed bool

	// Not sure how this is working
	// Should i check if inout is closed and buffer is nul then close the out channel
	for IsInputClosed && output != nil {
		select {
		case inIten, IsInputClosed = <-input:
			if IsInputClosed {
				ch.buffer.Add(inIten)
			}
		case output <- ch.buffer.Poll():
		case ch.length <- ch.buffer.Length():
		}

		if ch.buffer.Length() > 0 {
			output = ch.output
		} else {
			output = nil
		}
	}

	close(ch.output)
	close(ch.length)
}

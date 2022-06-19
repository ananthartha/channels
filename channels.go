package channels

type Buffer interface {
	Len() int // The number of elements currently buffered.
	Cap() int // The maximum number of elements that can be buffered.
}

type UnbufferedInChannel[T any] interface {
	In() chan<- *T // The writeable end of the channel.
	Close()        // Closes the channel. It is an error to write to In() after calling Close().
}

type InChannel[T any] interface {
	UnbufferedInChannel[T]
	Buffer
}

type UnbufferedOutChannel[T any] interface {
	Out() <-chan *T // The readable end of the channel.
}

type OutChannel[T any] interface {
	UnbufferedOutChannel[T]
	Buffer
}

type UnbufferedChannel[T any] interface {
	UnbufferedInChannel[T]
	UnbufferedOutChannel[T]
}

type Channel[T any] interface {
	UnbufferedChannel[T]
	Buffer
}

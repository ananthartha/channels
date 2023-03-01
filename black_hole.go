package channels

// BlackHole implements the InChannel interface and provides an analogue for the "Discard" variable in
// the ioutil package - it never blocks, and simply discards every value it reads. The number of items
// discarded in this way is counted and returned from Len.
type BlackHole[I any] struct {
	input  chan I
	length chan int
	count  int
}

func NewBlackHole[I any]() *BlackHole[I] {
	ch := &BlackHole[I]{
		input:  make(chan I),
		length: make(chan int),
	}
	go ch.discard()
	return ch
}

func (ch *BlackHole[I]) In() chan I {
	return ch.input
}

func (ch *BlackHole[I]) Len() int {
	val, open := <-ch.length
	if open {
		return val
	} else {
		return ch.count
	}
}

func (ch *BlackHole[I]) Cap() int {
	return -1
}

func (ch *BlackHole[I]) Close() {
	close(ch.input)
}

func (ch *BlackHole[I]) discard() {
	for {
		select {
		case _, open := <-ch.input:
			if !open {
				close(ch.length)
				return
			}
			ch.count++
		case ch.length <- ch.count:
		}
	}
}

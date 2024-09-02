package safechannel

import (
	"errors"
	"sync"
)

type SafeChannel[T any] struct {
	sendCh chan T
	recvCh chan T
	mu     sync.Mutex
	closed bool
}

func NewSafeChannel[T any](bufferSize int) *SafeChannel[T] {
	var sendCh, recvCh chan T

	if bufferSize > 0 {
		sendCh = make(chan T, bufferSize)
		recvCh = make(chan T, bufferSize)
	} else {
		sendCh = make(chan T)
		recvCh = make(chan T)
	}

	sc := &SafeChannel[T]{
		sendCh: sendCh,
		recvCh: recvCh,
		closed: false,
	}

	go sc.forwardMessages()
	return sc
}

func (sc *SafeChannel[T]) forwardMessages() {
	for {
		value, ok := <-sc.sendCh
		if !ok {
			// If sendCh is closed and all messages have been forwarded, close recvCh
			close(sc.recvCh)
			return
		}

		// Forward message to recvCh
		sc.recvCh <- value
	}
}

func (sc *SafeChannel[T]) Send(value T) error {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	if sc.closed {
		return errors.New("send on closed channel")
	}

	select {
	case sc.sendCh <- value:
		return nil
	default:
		return errors.New("channel buffer full")
	}
}

func (sc *SafeChannel[T]) Receive() (T, error) {
	value, ok := <-sc.recvCh
	if !ok {
		var zero T
		return zero, errors.New("receive on closed channel")
	}
	return value, nil
}

func (sc *SafeChannel[T]) Close() error {
	sc.mu.Lock()
	if sc.closed {
		sc.mu.Unlock()
		return errors.New("channel already closed")
	}

	close(sc.sendCh)
	sc.closed = true
	sc.mu.Unlock()

	return nil
}

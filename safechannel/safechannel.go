// package safechannel

// import (
// 	"errors"
// 	"sync"
// )

// // SafeChannel is a thread-safe wrapper around Go channels.
// type SafeChannel[T any] struct {
// 	sendCh chan T
// 	recvCh chan T
// 	closed bool
// 	mu     sync.Mutex
// }

// // NewSafeChannel creates a new SafeChannel with the given buffer size.
// func NewSafeChannel[T any](bufferSize int) *SafeChannel[T] {
// 	return &SafeChannel[T]{
// 		sendCh: make(chan T, bufferSize),
// 		recvCh: make(chan T, bufferSize),
// 		closed: false,
// 	}
// }

// // Send sends a value to the SafeChannel's sendCh and recvCh.
// func (sc *SafeChannel[T]) Send(value T) error {
// 	sc.mu.Lock()
// 	defer sc.mu.Unlock()

// 	if sc.closed {
// 		return errors.New("send on closed channel")
// 	}

// 	// Send to the sendCh channel and expect a goroutine to consume from recvCh
// 	select {
// 	case sc.sendCh <- value:
// 		// Block if necessary to ensure sendCh and recvCh are in sync
// 		// In case of unbuffered channels, ensure sync
// 		sc.recvCh <- value
// 	default:
// 		return errors.New("channel buffer full")
// 	}

// 	return nil
// }

// // Receive retrieves a value from the SafeChannel's recvCh.
// func (sc *SafeChannel[T]) Receive() (T, error) {
// 	value, ok := <-sc.recvCh
// 	if !ok {
// 		return *new(T), errors.New("receive on closed channel")
// 	}
// 	return value, nil
// }

// // Close closes the SafeChannel's sendCh and eventually recvCh after all messages are received.
// func (sc *SafeChannel[T]) Close() error {
// 	sc.mu.Lock()

// 	if sc.closed {
// 		sc.mu.Unlock()
// 		return errors.New("channel already closed")
// 	}

// 	sc.closed = true
// 	close(sc.sendCh) // Close the send channel to prevent further sends
// 	sc.mu.Unlock()

// 	// Start a goroutine to transfer remaining messages from sendCh to recvCh
// 	go func() {
// 		for msg := range sc.sendCh {
// 			sc.recvCh <- msg
// 		}
// 		close(sc.recvCh) // Close recvCh after all messages are transferred
// 	}()

// 	return nil
// }

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
	for value := range sc.sendCh {
		sc.mu.Lock()
		if sc.closed {
			sc.mu.Unlock()
			break
		}
		sc.recvCh <- value // Forward message to recvCh
		sc.mu.Unlock()
	}

	// Close recvCh after all messages are processed
	sc.mu.Lock()
	if !sc.closed {
		close(sc.recvCh)
		sc.closed = true
	}
	sc.mu.Unlock()
}

func (sc *SafeChannel[T]) Send(value T) error {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	if sc.closed {
		return errors.New("send on closed channel")
	}

	sc.sendCh <- value
	return nil
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
	defer sc.mu.Unlock()

	if sc.closed {
		return errors.New("channel already closed")
	}

	close(sc.sendCh)
	sc.closed = true
	return nil
}

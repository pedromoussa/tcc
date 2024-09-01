package safechannel

import (
	"errors"
	"sync"
)

// SafeChannel is a thread-safe wrapper around Go channels.
type SafeChannel[T any] struct {
	sendCh chan T
	recvCh chan T
	closed bool
	mu     sync.Mutex
}

// NewSafeChannel creates a new SafeChannel with the given buffer size.
func NewSafeChannel[T any](bufferSize int) *SafeChannel[T] {
	return &SafeChannel[T]{
		sendCh: make(chan T, bufferSize),
		recvCh: make(chan T, bufferSize),
		closed: false,
	}
}

// Send sends a value to the SafeChannel's sendCh and recvCh.
func (sc *SafeChannel[T]) Send(value T) error {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	if sc.closed {
		return errors.New("send on closed channel")
	}

	// Send to the sendCh channel and expect a goroutine to consume from recvCh
	select {
	case sc.sendCh <- value:
		// Block if necessary to ensure sendCh and recvCh are in sync
		// In case of unbuffered channels, ensure sync
		sc.recvCh <- value
	default:
		return errors.New("channel buffer full")
	}

	return nil
}

// Receive retrieves a value from the SafeChannel's recvCh.
func (sc *SafeChannel[T]) Receive() (T, error) {
	value, ok := <-sc.recvCh
	if !ok {
		return *new(T), errors.New("receive on closed channel")
	}
	return value, nil
}

// Close closes the SafeChannel's sendCh and eventually recvCh after all messages are received.
func (sc *SafeChannel[T]) Close() error {
	sc.mu.Lock()

	if sc.closed {
		sc.mu.Unlock()
		return errors.New("channel already closed")
	}

	sc.closed = true
	close(sc.sendCh) // Close the send channel to prevent further sends
	sc.mu.Unlock()

	// Start a goroutine to transfer remaining messages from sendCh to recvCh
	go func() {
		for msg := range sc.sendCh {
			sc.recvCh <- msg
		}
		close(sc.recvCh) // Close recvCh after all messages are transferred
	}()

	return nil
}

// type SafeChannel struct {
// 	sendCh chan interface{}
// 	recvCh chan interface{}
// 	closed bool
// 	// cache  []interface{}
// 	mu     sync.Mutex
// }

// func NewSafeChannel(bufferSize int) *SafeChannel {
// 	return &SafeChannel{
// 		sendCh: make(chan interface{}, bufferSize),
// 		recvCh: make(chan interface{}, bufferSize),
// 		// cache:  make([]interface{}, 0, bufferSize), // maybe keep it in the case of an unbuffered channel?
// 		closed: false,
// 	}
// }

// func (sc *SafeChannel) Send(value interface{}) error {
// 	sc.mu.Lock()
// 	defer sc.mu.Unlock()

// 	if sc.closed {
// 		return errors.New("send on closed channel")
// 	}

// 	// Armazenar no cache
// 	// sc.cache = append(sc.cache, value)

// 	select {
// 	case sc.sendCh <- value:
// 		sc.recvCh <- value // Propagar para recvCh
// 	default:
// 		return errors.New("channel buffer full")
// 	}

// 	return nil
// }

// func (sc *SafeChannel) Receive() (interface{}, error) {
// 	value, ok := <-sc.recvCh
// 	if !ok {
// 		return nil, errors.New("receive on closed channel")
// 	}
// 	return value, nil
// }

// func (sc *SafeChannel) Close() error {
// 	sc.mu.Lock()

// 	if sc.closed {
// 		sc.mu.Unlock()
// 		return errors.New("channel already closed")
// 	}

// 	// if sc.closed {
// 	// 	return
// 	// }

// 	sc.closed = true
// 	close(sc.sendCh) // Fechar o canal de envio
//   sc.mu.Unlock()

// 	// Propagar mensagens pendentes e limpar o cache
// 	go func() {
// 		for msg := range sc.sendCh {
// 			sc.recvCh <- msg
// 		}
// 		close(sc.recvCh)
// 	}()
// 	// go func() {
// 	// 	sc.mu.Lock()
// 	// 	defer sc.mu.Unlock()

// 	// 	for _, msg := range sc.cache {
// 	// 		sc.recvCh <- msg
// 	// 	}
// 	// 	close(sc.recvCh) // Fechar recvCh depois que todas as mensagens forem propagadas
// 	// 	sc.cache = nil   // Limpar o cache
// 	// }()

// 	return nil
// }

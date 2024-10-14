package safechannel

import (
	"errors"
	"sync"
)

type SafeChannel[T any] struct {
	sendCh chan T
	recvCh chan T
	notifyCh  chan Notification[T]
	mu     sync.Mutex
	closed bool
}

func NewSafeChannel[T any](bufferSize int) *SafeChannel[T] {
	var sendCh, recvCh chan T
	notifyCh := make(chan Notification[T], 10)

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
		notifyCh: notifyCh,
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

		sc.recvCh <- value
	}
}

func (sc *SafeChannel[T]) Send(value T) error {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	if sc.closed {
		sc.notify(Notification[T]{ReturnValue: -1, Message: "send on closed channel", Value: value})
		return errors.New("send on closed channel")
	}

	select {
	case sc.sendCh <- value:
		sc.notify(Notification[T]{ReturnValue: 0, Message: "sent successfully", Value: value})
		return nil
	default:
		sc.notify(Notification[T]{ReturnValue: -1, Message: "channel buffer full", Value: value})
		sc.sendCh <- value
		sc.notify(Notification[T]{ReturnValue: 0, Message: "sent after waiting", Value: value})
		return nil
	}
}

func (sc *SafeChannel[T]) Receive() (T, error) {
	value, ok := <-sc.recvCh
	if !ok {
		var zero T
		sc.notify(Notification[T]{ReturnValue: 0, Message: "receive on closed channel"})
		return zero, errors.New("receive on closed channel")
	}
	return value, nil
}

func (sc *SafeChannel[T]) Close() error {
	sc.mu.Lock()
	if sc.closed {
		sc.mu.Unlock()
		sc.notify(Notification[T]{ReturnValue: -1, Message: "channel already closed"})
		return errors.New("channel already closed")
	}

	close(sc.sendCh)
	sc.closed = true
	sc.mu.Unlock()

	return nil
}

type SelectCaseFunc func() (bool, int, interface{}, error)

func Select(cases ...SelectCaseFunc) (int, interface{}, error) {
	for {
		for _, c := range cases {
			ok, idx, value, err := c()
			if ok {
				return idx, value, err
			}
		}
	}
}

func CaseReceive[T any](sc *SafeChannel[T], onReceive func(T)) SelectCaseFunc {
	return func() (bool, int, interface{}, error) {
		select {
		case msg, ok := <-sc.recvCh:
			if !ok {
				sc.notify(Notification[T]{ReturnValue: -1, Message: "receive on closed channel", Value: msg})
				return true, -1, nil, errors.New("receive on closed channel")
			}
			if onReceive != nil {
				onReceive(msg)
			}
			return true, 0, msg, nil
		default:
			return false, 0, nil, nil
		}
	}
}

func CaseSend[T any](sc *SafeChannel[T], msg T) SelectCaseFunc {
	return func() (bool, int, interface{}, error) {
		if sc.closed {
			sc.notify(Notification[T]{ReturnValue: -1, Message: "send on closed channel", Value: msg})
			return true, -1, nil, errors.New("send on closed channel")
		}
		select {
		case sc.sendCh <- msg:
			return true, 0, nil, nil
		default:
			sc.notify(Notification[T]{ReturnValue: -1, Message: "channel buffer full", Value: msg})
			return true, -1, nil, errors.New("channel buffer full")
		}
	}
}

type Notification[T any] struct {
	ReturnValue int
	Message     string
	CallerID    string
	Value       T
}

func (sc *SafeChannel[T]) notify(notification Notification[T]) {
	select {
	case sc.notifyCh <- notification:
	default:
		// Ignore the notification if notifyCh is full
	}
}
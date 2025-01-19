package safechannel

import (
	"errors"
	"runtime"
	"sync"
	"sync/atomic"
)

type SafeChannel[T any] struct {
	sendCh chan T
	recvCh chan T
	notifyCh chan Notification[T]
	mu sync.Mutex
	cond *sync.Cond
	closed bool
	pendingReceivers int
	messagesInBuffer int64
}

func MakeSafechannel[T any](bufferSize ...int) *SafeChannel[T] {
	var sendCh, recvCh chan T
	var size = getBufferSize(bufferSize)
	
	sendCh = make(chan T, size)
	recvCh = make(chan T, size)

	sc := &SafeChannel[T]{
		sendCh: sendCh,
		recvCh: recvCh,
		closed: false,
	}
	sc.cond = sync.NewCond(&sc.mu)

	go sc.forwardMessages()
	return sc
}

func (sc *SafeChannel[T]) GetMessageCount() int64 {
	return atomic.LoadInt64(&sc.messagesInBuffer)
}

// func (sc *SafeChannel[T]) forwardMessages() {
// 	var channelsClosed bool

// 	for {

// 		if sc.closed {
// 			if !channelsClosed {
// 				close(sc.recvCh)
// 				if sc.notifyCh != nil {
// 					close(sc.notifyCh)
// 				}
// 				channelsClosed = true
// 			}
// 			return
// 		}

// 		// Attempt to forward the message to recvCh (same logic as in Send())
// 		select {
// 		case sc.recvCh <- <- sc.sendCh:
// 			sc.notify(Notification[T]{ReturnValue: 0, Message: "received successfully"})
// 		default:
// 			sc.notify(Notification[T]{ReturnValue: -1, Message: "channel buffer full"})
// 			sc.recvCh <- <-sc.sendCh
// 			sc.notify(Notification[T]{ReturnValue: 0, Message: "received after waiting"})
// 		}
// 	}
// }

func (sc *SafeChannel[T]) forwardMessages() {
	for {
		sc.mu.Lock()

		for sc.pendingReceivers == 0 && !sc.closed {
			sc.cond.Wait()
		}

		// if sc.closed {
		// 	close(sc.recvCh)
		// 	if sc.notifyCh != nil {
		// 		close(sc.notifyCh)
		// 	}
		// 	sc.mu.Unlock()
		// 	return
		// }

		sc.mu.Unlock()

		for {
			select {
			case value, ok := <-sc.sendCh:
				if !ok {
					close(sc.recvCh)
					if sc.notifyCh != nil {
						close(sc.notifyCh)
					}
					return
				}

				select {
				case sc.recvCh <- value:
					sc.notify(Notification[T]{ReturnValue: 0, Message: "message forwarded successfully", Value: value})
				default:
					sc.notify(Notification[T]{ReturnValue: -1, Message: "channel buffer full", Value: value})
					sc.recvCh <- value
					sc.notify(Notification[T]{ReturnValue: 0, Message: "forwarded after waiting", Value: value})
					sc.mu.Lock()
					sc.pendingReceivers--
					sc.mu.Unlock()
					goto StopForwarding
				}
			default:
				sc.mu.Lock()
				sc.pendingReceivers--
				sc.mu.Unlock()
				goto StopForwarding
			}
		}

	StopForwarding:
		continue
	}
}

func (sc *SafeChannel[T]) Send(value T) error {
	sc.mu.Lock()
	if sc.closed {
		sc.notify(Notification[T]{ReturnValue: -1, Message: "send on closed channel", Value: value})
		return errors.New("send on closed channel")
	}
	sc.mu.Unlock()

	select {
	case sc.sendCh <- value:
		sc.notify(Notification[T]{ReturnValue: 0, Message: "sent successfully", Value: value})
		atomic.AddInt64(&sc.messagesInBuffer, 1)
		return nil
	default:
		sc.notify(Notification[T]{ReturnValue: -1, Message: "channel buffer full", Value: value})
		sc.sendCh <- value
		atomic.AddInt64(&sc.messagesInBuffer, 1)
		sc.notify(Notification[T]{ReturnValue: 0, Message: "sent after waiting", Value: value})
		return nil
	}
}

// func (sc *SafeChannel[T]) Receive() (T, error) {
// 	value, ok := <-sc.recvCh
// 	if !ok {
// 		var zero T
// 		sc.notify(Notification[T]{ReturnValue: 0, Message: "receive on closed channel"})
// 		return zero, errors.New("receive on closed channel")
// 	}
// 	atomic.AddInt64(&sc.messagesInBuffer, -1)
// 	return value, nil
// }

func (sc *SafeChannel[T]) Receive() (T, error) {
	var zero T

	for {
		select {
		case value, ok := <-sc.recvCh:
			if !ok {
				sc.notify(Notification[T]{ReturnValue: -1, Message: "receive on closed channel"})
				return zero, errors.New("receive on closed channel")
			}
			atomic.AddInt64(&sc.messagesInBuffer, -1)
			return value, nil
		default:
			sc.mu.Lock()
			sc.pendingReceivers++
			sc.cond.Signal()
			sc.mu.Unlock()
		}
	}
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

// Trying to ensure default case is evaluated only when no other case is ready
func Select(cases ...SelectCaseFunc) (int, interface{}, error) {
	var defaultIdx int = -1
	for i, c := range cases {
			ok, idx, value, err := c()
			if ok {
					return idx, value, err
			}
			if idx == -1 {
					defaultIdx = i
			}
	}

	if defaultIdx != -1 {
			_, idx, value, err := cases[defaultIdx]()
			return idx, value, err
	}

	return -1, nil, nil
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

// Default case will now perform an action
func CaseDefault(action func()) SelectCaseFunc {
	return func() (bool, int, interface{}, error) {
			if action != nil {
					action()
			}
			return true, -1, nil, nil
	}
}

type Notification[T any] struct {
	ReturnValue int
	Message     string
	Value       T
	FuncName    string
	File        string
	Line        int
}

func (sc *SafeChannel[T]) EnableNotifications(bufferSize ...int) {
	var size = getBufferSize(bufferSize)
	sc.notifyCh = make(chan Notification[T], size)
}

func (sc *SafeChannel[T]) ReadNotification() (Notification[T], error) {
	if sc.notifyCh == nil {
			return Notification[T]{}, errors.New("notifications not enabled")
	}
	notification, ok := <-sc.notifyCh
	if !ok {
			return Notification[T]{}, errors.New("notification channel closed")
	}
	return notification, nil
}

func (sc *SafeChannel[T]) notify(notification Notification[T]) {
	if sc.notifyCh == nil {
		return
	}

	pc, file, line, ok := runtime.Caller(2)
	if ok {
		notification.FuncName = runtime.FuncForPC(pc).Name()
		notification.File = file
		notification.Line = line
	}
	select {
	case sc.notifyCh <- notification:
	default:
		<- sc.notifyCh
		sc.notifyCh <- notification
	}
}

func getBufferSize(bufferSize []int) int {
	if len(bufferSize) > 0 {
		return bufferSize[0]
	} else {
		return 0
	}
}

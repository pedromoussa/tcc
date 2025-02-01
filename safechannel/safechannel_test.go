package safechannel

import (
	"testing"
	"time"
)

func TestSafeChannel_NormalOperation(t *testing.T) {
	sc := MakeSafechannel[int](1)
	
	err := sc.Send(42)
	if err != nil {
		t.Errorf("Unexpected error on send: %v", err)
	}

	value, err := sc.Receive()
	if err != nil {
		t.Errorf("Unexpected error on receive: %v", err)
	}
	if value != 42 {
		t.Errorf("Expected 42, got %v", value)
	}
}

func TestSafeChannel_FullBuffer(t *testing.T) {
	sc := MakeSafechannel[int](1)
	sc.EnableNotifications(10)

	go func() {
		err := sc.Send(42)
		if err != nil {
			t.Errorf("Unexpected error on send: %v", err)
		}
		err = sc.Send(99)
		if err != nil {
			t.Errorf("Unexpected error on send: %v", err)
		}
	}()

	go func() {
		var messages [2]string
		index := 0

		for {
			notification, err := sc.ReadNotification()
			messages[index] = notification.Message
      index++
			if index == len(messages) {
				if messages[index-1] != "channel buffer full" || err != nil {
					t.Errorf("Expected 'channel buffer full' error, got %v", err)
				}
				break
			}
		}
	}()
}

/*
the check for nil channel was never implemented, the only risk is
if we have a nil SafeChannel (created but didn't call MakeSafechannel)
i do not know what's the expected behaviour in that case
might be worth looking into it
*/
// func TestSafeChannel_NilChannel(t *testing.T) {
// 	var sc *SafeChannel[int]

// 	err := sc.Send(42)
// 	if err == nil {
// 		t.Errorf("Expected error on sending to nil channel")
// 	}

// 	_, err = sc.Receive()
// 	if err == nil {
// 		t.Errorf("Expected error on receiving from nil channel")
// 	}
// }

func TestSafeChannel_SendToClosed(t *testing.T) {
	sc := MakeSafechannel[int](0)

	sc.Close()

	err := sc.Send(42)
	if err == nil || err.Error() != "send on closed channel" {
		t.Errorf("Expected 'send on closed channel' error, got %v", err)
	}
}

func TestSafeChannel_ReceiveFromClosed(t *testing.T) {
	sc := MakeSafechannel[int](1)

	err := sc.Send(42)
	if err != nil {
		t.Errorf("Unexpected error on send: %v", err)
	}

	sc.Close()

	value, err := sc.Receive()
	if err != nil {
		t.Errorf("Unexpected error on receive: %v", err)
	}
	if value != 42 {
		t.Errorf("Expected 42, got %v", value)
	}

	_, err = sc.Receive()
	if err == nil || err.Error() != "receive on closed channel" {
		t.Errorf("Expected 'receive on closed channel' error, got %v", err)
	}
}

// receive should listen until a message arrives from send
func TestSafeChannel_Concurrent(t *testing.T) {
	sc := MakeSafechannel[int]()

	go func() {
		time.Sleep(50 * time.Millisecond)
		err := sc.Send(99)
		if err != nil {
			t.Errorf("Unexpected error on send: %v", err)
		}
	}()

	value, err := sc.Receive()
	if err != nil {
		t.Errorf("Unexpected error on receive: %v", err)
	}
	if value != 99 {
		t.Errorf("Expected 99, got %v", value)
	}
}

func TestSafeChannel_SendWithoutReceive(t *testing.T) {
	sc := MakeSafechannel[int](1)
	done := make(chan error)

	err := sc.Send(42)
	if err != nil {
		t.Errorf("Unexpected error on send: %v", err)
	}

	go func() {
		err = sc.Send(99)
		done <- err
	}()

	select {
	case <-done:
		t.Fatal("Test failed: Send should block until a receive occurs")
	case <-time.After(500 * time.Millisecond):
		t.Log("Send is correctly blocked in unbuffered SafeChannel")
	}

	go func() {
		_, _ = sc.Receive()
	}()
	time.Sleep(200 * time.Millisecond)

	select {
	case err := <-done:
		if err != nil {
			t.Errorf("Unexpected error on send: %v", err)
		} else {
			t.Log("Send completed successfully after receive")
		}
	case <-time.After(1 * time.Second):
		t.Fatal("Test failed: Send did not unblock after a receive")
	}
}

func TestSafeChannel_SendWithoutReceiveUnbuffered(t *testing.T) {
	sc := MakeSafechannel[int]()
	sc.EnableNotifications(4)
	done := make(chan error)

	go func() {
		err := sc.Send(42)
		done <- err
	}()

	go func() {
		var messages [4]string
		index := 0

		for {
			notification, _ := sc.ReadNotification()
			t.Log(notification.Message)
      index++
			if index == len(messages) {
				break
			}
		}
	}()

	select { 
	case <-done: 
		t.Fatal("Test failed: Send should block until a receive occurs")
	case <-time.After(500 * time.Millisecond):
		t.Log("Send is correctly blocked in unbuffered SafeChannel")
	}

	go func() {
		_, _ = sc.Receive()
	}()
	time.Sleep(200 * time.Millisecond)

	select {
	case err := <-done:
		if err != nil {
			t.Errorf("Unexpected error on send: %v", err)
		} else {
			t.Log("Send completed successfully after receive")
		}
	case <-time.After(1 * time.Second):
		t.Fatal("Test failed: Send did not unblock after a receive")
	}
}

func TestSafeChannel_CloseTwice(t *testing.T) {
	sc := MakeSafechannel[int](0)

	err := sc.Close()
	if err != nil {
		t.Errorf("Unexpected error on first close: %v", err)
	}

	err = sc.Close()
	if err == nil || err.Error() != "channel already closed" {
		t.Errorf("Expected 'channel already closed' error, got %v", err)
	}
}

func TestSafeChannel_CustomSelect(t *testing.T) {
	sc := MakeSafechannel[int](1)

 	err := sc.Send(42)
 	if err != nil {
 		t.Errorf("Unexpected error on send: %v", err)
 	}

 	idx, value, err := Select(
 		CaseReceive(sc, nil),
 	)

	if err != nil {
		t.Errorf("Unexpected error on select: %v", err)
	}
	if idx != 0 {
		t.Errorf("Expected case index 0, got %v", idx)
	}
	if value.(int) != 42 {
		t.Errorf("Expected value 42, got %v", value)
	}
}

func TestSafeChannel_Notifications(t *testing.T) {
	sc := MakeSafechannel[int](1) 
	sc.EnableNotifications(10)

	go func() {
		err := sc.Send(42)
		if err != nil {
			t.Errorf("Unexpected error on send: %v", err)
		}
	}()

	notification, err := sc.ReadNotification()
	if err != nil {
		t.Errorf("Failed to read notification: %v", err)
	}
	if notification.Message != "sent successfully" {
		t.Errorf("Expected 'sent successfully', got %v", notification.Message)
	}
}

func TestSafeChannel_CustomSelectDefault(t *testing.T) {
	sc := MakeSafechannel[int](1)

	go func() {
		time.After(50 * time.Millisecond)
		sc.Send(42)
	}()

	idx, value, err := Select(
		CaseReceive(sc, nil),
		CaseDefault(func() {
			t.Log("Default case executed")
		}),
	)

	if idx != 0 || value.(int) != 42 || err != nil {
		t.Errorf("Expected case index 0 with value 42, got idx: %v, value: %v, err: %v", idx, value, err)
	}
}
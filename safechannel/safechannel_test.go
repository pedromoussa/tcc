package safechannel

import (
	// "errors"
	// "sync"
	"testing"
	"time"
)

func TestSafeChannel_NormalOperation(t *testing.T) {
	sc := MakeSafeChannel[int](1)
	
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

// FIX
func TestSafeChannel_FullBuffer(t *testing.T) {
	sc := MakeSafeChannel[int](1)

	err := sc.Send(42)
	if err != nil {
		t.Errorf("Unexpected error on send: %v", err)
	}

	err = sc.Send(99)
	if err == nil || err.Error() != "channel buffer full" {
		t.Errorf("Expected 'channel buffer full' error, got %v", err)
	}
}

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
	sc := MakeSafeChannel[int](0)

	sc.Close()

	err := sc.Send(42)
	if err == nil || err.Error() != "send on closed channel" {
		t.Errorf("Expected 'send on closed channel' error, got %v", err)
	}
}

func TestSafeChannel_ReceiveFromClosed(t *testing.T) {
	sc := MakeSafeChannel[int](1)

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

func TestSafeChannel_Concurrent(t *testing.T) {
	sc := MakeSafeChannel[int](0)

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

// FIX
func TestSafeChannel_SendWithoutReceive(t *testing.T) {
	sc := MakeSafeChannel[int](1)

	err := sc.Send(42)
	if err != nil {
		t.Errorf("Unexpected error on send: %v", err)
	}

	err = sc.Send(99)
	if err == nil || err.Error() != "channel buffer full" {
		t.Errorf("Expected 'channel buffer full' error, got %v", err)
	}
}

func TestSafeChannel_SendWithoutReceiveUnbuffered(t *testing.T) {
	sc := MakeSafeChannel[int](0)

	done := make(chan error)
    
	go func() {
		err := sc.Send(42)
		done <- err
	}()

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("Unexpected error on send: %v", err)
		} else {
			t.Log("Send completed successfully")
		}
	case <-time.After(1 * time.Second):
		t.Fatal("Test failed: Send blocked due to missing receive")
	}
}

func TestSafeChannel_CloseTwice(t *testing.T) {
	sc := MakeSafeChannel[int](0)

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
	sc := MakeSafeChannel[int](1)

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

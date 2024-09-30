package safechannel

import (
	// "errors"
	// "sync"
	"testing"
	"time"
)

// func TestSafeChannelSendReceive(t *testing.T) {
// 	sc := NewSafeChannel[int](0)

// 	// test sending and receiving a value
// 	err := sc.Send(42)
// 	if err != nil {
// 		t.Fatalf("Expected no error, got %v", err)
// 	}

// 	value, err := sc.Receive()
// 	if err != nil {
// 		t.Fatalf("Expected no error, got %v", err)
// 	}
// 	if value != 42 {
// 		t.Fatalf("Expected value 42, got %v", value)
// 	}
// }

// func TestSafeChannelBuffer(t *testing.T) {
// 	sc := NewSafeChannel[int](5) 

// 	// test sending into a buffered channel
// 	err := sc.Send(10)
// 	if err != nil {
// 		t.Fatalf("Expected no error, got %v", err)
// 	}

// 	err = sc.Send(20)
// 	if err != nil {
// 		t.Fatalf("Expected no error, got %v", err)
// 	}

// 	// test receiving from a buffered channel
// 	value, err := sc.Receive()
// 	if err != nil {
// 		t.Fatalf("Expected no error, got %v", err)
// 	}
// 	if value != 10 {
// 		t.Fatalf("Expected value 10, got %v", value)
// 	}

// 	value, err = sc.Receive()
// 	if err != nil {
// 		t.Fatalf("Expected no error, got %v", err)
// 	}
// 	if value != 20 {
// 		t.Fatalf("Expected value 20, got %v", value)
// 	}
// }

// func TestSafeChannelClose(t *testing.T) {
// 	sc := NewSafeChannel[int](0) 

// 	// test closing the channel
// 	err := sc.Close()
// 	if err != nil {
// 		t.Fatalf("Expected no error on first close, got %v", err)
// 	}

// 	err = sc.Close()
// 	if err == nil {
// 		t.Fatalf("Expected error on second close, got nil")
// 	}
// }

// func TestSafeChannelSendAfterClose(t *testing.T) {
// 	sc := NewSafeChannel[int](0) 

// 	sc.Close()

// 	// test sending to a closed channel (should return an error)
// 	err := sc.Send(42)
// 	if err == nil {
// 		t.Fatalf("Expected error when sending to closed channel, got nil")
// 	}
// }

// func TestSafeChannelReceiveAfterClose(t *testing.T) {
// 	sc := NewSafeChannel[int](0) 

// 	sc.Close()

// 	// test receiving from a closed channel (should return zero value and error)
// 	value, err := sc.Receive()
// 	if err == nil {
// 		t.Fatalf("Expected error when receiving from closed channel, got nil")
// 	}
// 	if value != 0 {
// 		t.Fatalf("Expected zero value, got %v", value)
// 	}
// }

// func TestSafeChannelNilSend(t *testing.T) {

// 	t.Skip("Skipping TestSafeChannelNilSend - nil channel blocking not testable yet in unit tests")
// }

// func TestSafeChannelNilReceive(t *testing.T) {

// 	t.Skip("Skipping TestSafeChannelNilReceive - nil channel blocking not testable yet in unit tests")
// }

// func TestSafeChannelBufferFull(t *testing.T) {
// 	sc := NewSafeChannel[int](2) 

// 	// send a value to fill the buffer
// 	err := sc.Send(1)
// 	if err != nil {
// 		t.Fatalf("Expected no error, got %v", err)
// 	}

// 	// send another value, which should result in buffer full error
// 	err = sc.Send(2)
// 	if err == nil {
// 		t.Fatalf("Expected error when sending to full buffer, got nil")
// 	}
// }

// func TestSafeChannelForwardMessages(t *testing.T) {
// 	sc := NewSafeChannel[int](2) 

// 	go func() {
// 		err := sc.Send(1)
// 		if err != nil {
// 			t.Fatalf("Expected no error, got %v", err)
// 		}
// 	}()

// 	value, err := sc.Receive()
// 	if err != nil {
// 		t.Fatalf("Expected no error, got %v", err)
// 	}
// 	if value != 1 {
// 		t.Fatalf("Expected value 1, got %v", value)
// 	}
// }

// /*
// func TestSafeChannelPanicOnClosedSend(t *testing.T) {
// 	defer func() {
// 		if r := recover(); r != nil {
// 			t.Logf("Recovered from panic: %v", r)
// 		} else {
// 			t.Fatalf("Expected panic, but no panic occurred")
// 		}
// 	}()

// 	sc := NewSafeChannel[int](0) 

// 	// enable panic mode (when implemented in the package)
// 	// sc.panicOnClosedSend = true

// 	sc.Close()

// 	// this should panic
// 	sc.Send(1)
// }
// */

// func TestSafeChannel_SendWithoutReceive(t *testing.T) {
//     sc := NewSafeChannel[int](0)

//     done := make(chan bool)
    
//     go func() {
//         err := sc.Send(42)
//         if err != nil {
//             t.Fatalf("Unexpected error on send: %v", err)
//         }
//         done <- true
//     }()
    
//     select {
//     case <-done:
//         t.Log("Send completed successfully")
//     case <-time.After(1 * time.Second):
//         t.Fatal("Test failed: Send blocked due to missing receive")
//     }
// }

// func TestSafeChannel_ChannelAndLock(t *testing.T) {
//     var mu sync.Mutex
//     sc := NewSafeChannel[int](0)

//     go func() {
//         mu.Lock()
//         defer mu.Unlock()
//         sc.Send(42) // this may block if no receiver
//     }()

//     go func() {
//         mu.Lock() // second goroutine tries to lock
//         defer mu.Unlock()
//         sc.Receive() // this may block if no sender
//     }()

//     select {
//     case <-time.After(1 * time.Second):
//         t.Fatal("Test failed: Deadlock occurred")
//     default:
//         t.Log("No deadlock occurred")
//     }
// }

// func TestSafeChannel_ProperClose(t *testing.T) {
//     sc := NewSafeChannel[int](0) 

//     go func() {
//         sc.Send(10)
//     }()

//     time.Sleep(100 * time.Millisecond) // wait a little
//     err := sc.Close()
//     if err != nil {
//         t.Fatalf("Unexpected error on close: %v", err)
//     }

//     _, err = sc.Receive()
//     if err == nil {
//         t.Fatal("Expected error on receiving from closed channel")
//     }
// }

// // Test to verify a successful send and receive operation using the custom Select function.
// func TestSelectSendReceive(t *testing.T) {
// 	sc := NewSafeChannel[string](0)
// 	msgToSend := "Hello SafeChannel"

// 	var wg sync.WaitGroup
// 	wg.Add(1)

// 	// Run a goroutine to send and receive using select
// 	go func() {
// 		defer wg.Done()

// 		_, receivedValue, err := Select(
// 			CaseSend(sc, msgToSend),
// 			CaseReceive(sc, func(msg string) {
// 				t.Logf("Received message: %s", msg)
// 			}),
// 		)

// 		if err != nil {
// 			t.Errorf("Unexpected error: %v", err)
// 		}

// 		if receivedValue != msgToSend {
// 			t.Errorf("Expected %v, but got %v", msgToSend, receivedValue)
// 		}
// 	}()

// 	// Wait for goroutine to finish
// 	wg.Wait()
// }

// // Test case when trying to receive from a closed channel
// func TestSelectReceiveClosedChannel(t *testing.T) {
// 	sc := NewSafeChannel[int](0)
// 	sc.Close() // Close the channel

// 	_, _, err := Select(
// 		CaseReceive(sc, nil),
// 	)

// 	if err == nil {
// 		t.Error("Expected an error, but got nil")
// 	} else if err.Error() != "receive on closed channel" {
// 		t.Errorf("Expected 'receive on closed channel', but got %v", err)
// 	}
// }

// // Test case when sending to a closed channel should fail
// func TestSelectSendClosedChannel(t *testing.T) {
// 	sc := NewSafeChannel[string](0)
// 	sc.Close() // Close the channel

// 	_, _, err := Select(
// 		CaseSend(sc, "some message"),
// 	)

// 	if err == nil {
// 		t.Error("Expected an error when sending to a closed channel")
// 	} else if err.Error() != "send on closed channel" {
// 		t.Errorf("Expected 'send on closed channel', but got %v", err)
// 	}
// }

// // Test case for buffered channels that are full
// func TestSelectSendBufferFull(t *testing.T) {
// 	sc := NewSafeChannel[string](1)
// 	sc.Send("First Message") // Fill the buffer

// 	_, _, err := Select(
// 		CaseSend(sc, "Second Message"),
// 	)

// 	if err == nil {
// 		t.Error("Expected error when sending to a full buffer")
// 	} else if err.Error() != "channel buffer full" {
// 		t.Errorf("Expected 'channel buffer full', but got %v", err)
// 	}
// }

// // Test to verify that multiple cases can work with Select
// func TestMultipleSelectCases(t *testing.T) {
// 	sc1 := NewSafeChannel[string](0)
// 	sc2 := NewSafeChannel[string](0)

// 	msgToSend := "Hello SafeChannel"

// 	// Send to both channels
// 	sc1.Send(msgToSend)
// 	sc2.Send(msgToSend)

// 	var wg sync.WaitGroup
// 	wg.Add(2)

// 	// Test multiple cases with select
// 	go func() {
// 		defer wg.Done()
// 		_, receivedValue, err := Select(
// 			CaseReceive(sc1, nil),
// 			CaseReceive(sc2, nil),
// 		)

// 		if err != nil {
// 			t.Errorf("Unexpected error: %v", err)
// 		}

// 		if receivedValue != msgToSend {
// 			t.Errorf("Expected %v, but got %v", msgToSend, receivedValue)
// 		}
// 	}()

// 	wg.Wait()
// }

// // Test for non-blocking receive using default case
// func TestNonBlockingSelect(t *testing.T) {
// 	sc := NewSafeChannel[int](0)

// 	_, value, err := Select(
// 		CaseReceive(sc, nil),
// 	)

// 	if err != nil && !errors.Is(err, nil) {
// 		t.Errorf("Expected no error but got %v", err)
// 	}

// 	if value != nil {
// 		t.Errorf("Expected nil value, but got %v", value)
// 	}
// }

/*
Select(
    CaseSend(sc1, msgToSend),
    CaseReceive(sc2, func(msg string) {
        fmt.Println("Received from sc2:", msg)
    }),
)
*/

func TestSafeChannel_NormalOperation(t *testing.T) {
	sc := NewSafeChannel[int](1)
	
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
	sc := NewSafeChannel[int](1)

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
	sc := NewSafeChannel[int](0)

	sc.Close()

	err := sc.Send(42)
	if err == nil || err.Error() != "send on closed channel" {
		t.Errorf("Expected 'send on closed channel' error, got %v", err)
	}
}

func TestSafeChannel_ReceiveFromClosed(t *testing.T) {
	sc := NewSafeChannel[int](1)

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
	sc := NewSafeChannel[int](0)

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
	sc := NewSafeChannel[int](1)

	err := sc.Send(42)
	if err != nil {
		t.Errorf("Unexpected error on send: %v", err)
	}

	err = sc.Send(99)
	if err == nil || err.Error() != "channel buffer full" {
		t.Errorf("Expected 'channel buffer full' error, got %v", err)
	}
}

func TestSafeChannel_CloseTwice(t *testing.T) {
	sc := NewSafeChannel[int](0)

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
	sc := NewSafeChannel[int](1)

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

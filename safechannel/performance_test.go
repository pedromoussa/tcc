package safechannel

import (
	"testing"
)

const numOperations, bufferSize = 1000000, 100

func BenchmarkNativeChannel(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ch := make(chan int, bufferSize)

		b.StartTimer()
		go func() {
			for j := 0; j < numOperations; j++ {
				ch <- j
			}
			close(ch)
		}()

		for v:= range ch {
			_ = v
		}
		b.StopTimer()
	}
}

func BenchmarkSafechannel(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sc := MakeSafechannel[int](bufferSize)

		b.StartTimer()
		go func() {
			for j := 0; j < numOperations; j++ {
				sc.Send(j)
			}
			sc.Close()
		}()

		for {
			_, err := sc.Receive()
			if err != nil {
				break
			}
		}
		b.StopTimer()
	}
}

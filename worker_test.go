package goueue

import (
	"sync"
	"testing"
	"time"
)

func TestWorker_NewStartStop(t *testing.T) {
	pool := make(chan chan Command)
	tests := []struct {
		name string
		pool chan chan Command
	}{
		{
			name: "worker test",
			pool: pool,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			timeout := time.After(30 * time.Second)
			done := make(chan bool)
			go func() {
				wg := &sync.WaitGroup{}
				w := NewWorker(tt.pool, wg)
				w.Start()
				go w.Stop()
				wg.Wait()
				done <- true
			}()

			select {
			case <-timeout:
				t.Fatalf("Test \"%s\" reached timeout", tt.name)
			case <-done:
			}
		})
	}
}

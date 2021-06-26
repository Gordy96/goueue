# goueue
Basic command queue/worker implementation

Basic usage:
```go
package main

import (
	"fmt"
	"time"

	"github.com/Gordy96/goueue"
)

type task struct {
	timestamp	int64
}

func (t *task) Handle() error {
	fmt.Printf("EPOCH NOW IS: %d\n", t.timestamp)
	return nil
}

func main() {
	numOfWorkerRoutines := 10
	q := goueue.New(numOfWorkerRoutines)
	q.Start()
	go func ()  {
		c := time.Tick(50 * time.Millisecond)
		for range c {
			q.Enqueue(&task{timestamp: time.Now().Unix()})
		}
	}()
	time.AfterFunc(10*time.Second, func() {
		q.Stop()
	})
	q.Wait()
}
```

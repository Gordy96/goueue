# goueue
Basic command queue/worker implementation

Basic usage:
```go
package main

import (
	"fmt"

	"github.com/Gordy96/goueue"
)

type task struct {
	id	int
}

func (t *task) Handle() error {
	fmt.Printf("Task #%d is being handled\n", t.id)
	return nil
}

func main() {
	numOfWorkerRoutines := 10
	q := goueue.New(numOfWorkerRoutines)
	q.Start()
	for i := 0; i < 10000; i++ {
		q.Enqueue(&task{id: i})
	}
	q.Stop()
	defer q.Wait()
}
```
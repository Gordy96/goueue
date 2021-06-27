package goueue

import "sync"

func NewWorker(pool chan chan Command, wg *sync.WaitGroup) *Worker {
	return &Worker{
		pool,
		make(chan Command),
		make(chan bool),
		wg,
	}
}

//Worker simple worker that executes commands / notifies queue about readiness
type Worker struct {
	pool  chan chan Command
	input chan Command
	quit  chan bool
	wg    *sync.WaitGroup
}

//Start increases sync.WaitGroup delta by 1 and spawns a goroutine
//that notifies queue about worker being ready to receive a new tasks concurrently
//and waits for either a new command to complete or signal to stop and decrease sync.WaitGroup by 1
func (w *Worker) Start() {
	w.wg.Add(1)
	go func() {
		for {
			go func() {
				w.pool <- w.input
			}()
			select {
			case job := <-w.input:
				err := job.Handle()
				if err != nil {
					break
				}
			case <-w.quit:
				break
			}
		}
		w.wg.Done()
	}()
}

//Stop notifies worker's own quit-channel
func (w *Worker) Stop() {
	w.quit <- true
}

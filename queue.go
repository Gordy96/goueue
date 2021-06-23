package goueue

import "sync"

//Queue implements queue/worker pattern with channels. Uses sync.WaitGroup optionally to wait for all workers to end
type Queue struct {
	workerCount int
	workers     []*Worker
	pool        chan chan Command
	jobQueue    chan Command
	wg          sync.WaitGroup
}

//Start creates and starts workers (launches a new goroutine)
func (q *Queue) Start() {
	for i := 0; i < q.workerCount; i++ {
		worker := NewWorker(q.pool, &(q.wg))
		worker.Start()
		q.workers[i] = worker
	}
	go func() {
		for {
			select {
			case job := <-q.jobQueue:
				go func(c Command) {
					workerChan := <-q.pool
					workerChan <- c
				}(job)
			}
		}
	}()
}

//Enqueue writes passed command to a channel
func (q *Queue) Enqueue(command Command) {
	q.jobQueue <- command
}

//Wait waits for sync.WaitGroup to deplete
func (q *Queue) Wait() {
	q.wg.Wait()
}

//Stop commands all workers to stop
func (q *Queue) Stop() {
	for _, worker := range q.workers {
		worker.Stop()
	}
}

//New creates a new Queue object with numWorkers workers
func New(numWorkers int) *Queue {
	q := &Queue{
		pool:        make(chan chan Command),
		workerCount: numWorkers,
		jobQueue:    make(chan Command),
		workers:     make([]*Worker, numWorkers),
		wg:          sync.WaitGroup{},
	}
	return q
}

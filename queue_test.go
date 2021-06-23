package goueue

import (
	"fmt"
	"sync"
	"testing"
)

func TestNew(t *testing.T) {
	type args struct {
		numWorkers int
	}
	numWorkers := 10
	nominal := &Queue{
		pool:        make(chan chan Command),
		workerCount: numWorkers,
		jobQueue:    make(chan Command),
		workers:     make([]*Worker, numWorkers),
		wg:          sync.WaitGroup{},
	}
	tests := []struct {
		name string
		args args
		want *Queue
	}{
		{
			name: "nominal queue",
			args: args{numWorkers: numWorkers},
			want: nominal,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New(tt.args.numWorkers)
			if got == nil {
				t.Error("New() should not return nil")
			}
		})
	}
}

type test_task struct {
	id string
	wg *sync.WaitGroup
	t  *testing.T
}

func (task *test_task) Handle() error {
	defer task.wg.Done()
	task.t.Logf("Task #%s is being handled", task.id)
	return nil
}

func TestQueue(t *testing.T) {
	makeTask := func(id string, wg *sync.WaitGroup, t *testing.T) *test_task {
		wg.Add(1)
		return &test_task{
			id: id,
			wg: wg,
			t:  t,
		}
	}
	tests := []struct {
		name        string
		workerCount int
		taskCount   int
	}{
		{
			name:        "queue test run",
			workerCount: 5,
			taskCount:   10,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := New(tt.workerCount)
			q.Start()
			alldone := sync.WaitGroup{}
			for i := 0; i < tt.taskCount; i++ {
				task := makeTask(fmt.Sprintf("task-%d", i), &alldone, t)
				q.Enqueue(task)
			}
			alldone.Wait()
			q.Stop()
		})
	}
}

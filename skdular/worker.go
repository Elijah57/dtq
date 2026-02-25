package skdular

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

// define worker struct
type Worker struct {
	ID         string
	ctx        context.Context
	numWorkers int
	MaxWorker  int
	MinWorker  int
	TaskQueue  <-chan Task
	wg         *sync.WaitGroup
}

type workerOptions func(*Worker)

func WithNumWorkers(n int) workerOptions {
	return func(w *Worker) {
		if n > 0 {
			w.numWorkers = n
		}
	}
}

func WithMaxWorker(n int) workerOptions {
	return func(w *Worker) {
		if n > 0 {
			w.MaxWorker = n
		}
	}
}

func WithMinWorker(n int) workerOptions {
	return func(w *Worker) {
		if n > 0 {
			w.MinWorker = n
		}
	}
}

// init worker
func NewWorker(ctx context.Context, taskQueue <-chan Task, wg *sync.WaitGroup, opts ...workerOptions) *Worker {

	w := &Worker{
		ID:         fmt.Sprintf("worker-%d", time.Now().UnixNano()),
		ctx:        ctx,
		numWorkers: 5,
		MaxWorker:  10,
		MinWorker:  2,
		TaskQueue:  taskQueue,
		wg:         wg,
	}

	for _, opt := range opts {
		opt(w)
	}

	return w
}

// start workers with default numWorkers
func (w *Worker) StartWorkers() {
	for i := 0; i < w.numWorkers; i++ {
		go w.WorkerNodeRun(i)
	}
	log.Println("Workers started")
}

// each worker node runs, listening on two channels - one for tasks and one for signals
func (w *Worker) WorkerNodeRun(id int) {
	for {
		select {
		case <-w.ctx.Done():
			fmt.Println("Worker", id, "channel closed, exiting ..")
			return
		case task, ok := <-w.TaskQueue: // check if channel is closed
			if !ok { // check if channel is closed and drained
				fmt.Println("Worker", id, "stopped")
				return
			}
			fmt.Println("Worker", id, "processing task", task.ID)
			w.wg.Add(1)
			_ = w.process(task, id)
		}
	}
}

func (w *Worker) process(task Task, workerId int) error {
	start := time.Now()
	log.Printf("Processing task %d", task.ID)
	// Simulate work
	// time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
	sum := cpuHeavyTask(task.Payload)
	fmt.Println("Worker", workerId, "finished task", task.ID, "with sum", sum)
	log.Printf("Finished task %d in %v", task.ID, time.Since(start))
	w.wg.Done()
	return nil
}

// func cpuHeavyTask(n int) float64 {
// 	sum := 0.0
// 	for i := 0; i < 1000000; i++ { // simulate CPU work
// 		sum += math.Sqrt(float64(i + n))
// 	}
// 	return sum
// }

func cpuHeavyTask(n int) int {
	a, b := 0, 1
	for i := 0; i < n; i++ { // large n like 50_000_000
		a, b = b, a+b
		a %= 1_000_000_007
	}
	return a
}

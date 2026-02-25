package skdular

import (
	"math/rand"
	"time"
)

// define task struct
type Task struct {
	ID        int
	Priority  int
	Payload   int
	CreatedAt time.Time
}

// task queue
// var taskQueue = make(chan Task, 100)

func GenerateTasks(queue chan Task, n int) {
	for i := 0; i < n; i++ {

		// generate a random number between 1_000_000 and 10_000_000
		num := 1_000_000_000 + rand.Intn(10_000_000-1_000_000+1)
		queue <- Task{ID: i, Payload: num}
	}

	close(queue) // optional: indicates no more tasks
}

package main

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/Elijah57/dtq/skdular"
)

func main() {

	wg := &sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	taskQueue := make(chan skdular.Task, 100)

	// init worker instance
	worker := skdular.NewWorker(ctx, taskQueue, wg, skdular.WithNumWorkers(100))
	worker.StartWorkers()

	// start task generation
	numTask := 1000
	start := time.Now()
	skdular.GenerateTasks(taskQueue, numTask)

	wg.Wait()
	log.Printf("Finished %d task(s) in %v", numTask, time.Since(start))
}

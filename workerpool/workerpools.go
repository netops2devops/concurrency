package main

import (
	"fmt"
	"os"
	"runtime/trace"
	"time"
)

const (
	MAX_WORKER_COUNT = 5
	MAX_JOB_COUNT    = 100
)

func Worker(id int, jobs <-chan int, results chan<- int) {
	for j := range jobs {
		fmt.Println("worker", id, "job", j)
		time.Sleep(time.Millisecond * 500)
		// fmt.Println("worker", id, "finished job", j)
		results <- j * 2
	}
}

func main() {
	f, _ := os.Create("trace.out")
	trace.Start(f)
	defer trace.Stop()

	jobs := make(chan int, 100)
	results := make(chan int, 100)

	for w := 1; w <= MAX_WORKER_COUNT; w++ {
		go Worker(w, jobs, results)
	}

	for j := range MAX_JOB_COUNT {
		jobs <- j
	}
	close(jobs)

	for range MAX_JOB_COUNT {
		<-results
	}
	close(results)
}

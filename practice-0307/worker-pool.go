// package main

// import (
// 	"fmt"
// 	"math/rand"
// 	"sync"
// 	"time"
// )

// type Job struct {
// 	id          int
// 	description string
// }

// type WorkerPool struct {
// 	numWorker int
// 	Process func(chan Job)
// 	results   []Job
// }

// func initJob() chan Job {
// 	randomNum := rand.Intn(8) + 3
// 	jobs := make(chan Job, randomNum)
// 	for i := 1; i <= randomNum; i++ {
// 		jobs <- Job{
// 			id:          i,
// 			description: fmt.Sprintf("Job %d", i),
// 		}
// 	}
// 	close(jobs) 
// 	return jobs
// }

// // Worker is the worker goroutine function.
// func Worker(id int, jobs <-chan Job, results chan<- Job, wg *sync.WaitGroup, ) {
// 	defer wg.Done()
// 	for j := range jobs {
// 		fmt.Printf("Worker %d: Process job %d\n", id, j.id)
// 		duration := rand.Intn(400) + 100
// 		time.Sleep(time.Duration(duration) * time.Millisecond)
// 		fmt.Printf("Worker %d: Finished job %d\n", id, j.id)
// 		results <- j
// 	}
// }



// // Init is an exported alias for init
// func (wp *WorkerPool) Init(numWorker int) {
// 	wp.numWorker = numWorker
// }


// func (wp *WorkerPool) process(jobs <-chan Job) {
// 	var wg sync.WaitGroup
// 	resultsChan := make(chan Job, 100)

// 	for i := 1; i <= wp.numWorker; i++ {
// 		wg.Add(1)
// 		go Worker(i, jobs, resultsChan, &wg)
// 	}

// 	go func() {
// 		wg.Wait()
// 		close(resultsChan)
// 	}()

// 	wp.results = nil
// 	for res := range resultsChan {
// 		wp.results = append(wp.results, res)
// 	}
// }

// func main() {
// 	var wp WorkerPool
// 	jobs := initJob()

// 	wp.Init(3)
// 	wp.process(jobs)

// 	fmt.Printf("\n--- All jobs finished! Total processed: %d ---\n", len(wp.results))
// 	for _, res := range wp.results {
// 		fmt.Printf("Result -> ID: %d, Description: %s\n", res.id, res.description)
// 	}
// }
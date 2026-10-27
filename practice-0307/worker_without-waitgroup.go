// package main

// import (
// 	"fmt"
// 	"time"
// )

// const numJobs = 5
// const numWorker = 3

// func worker(id int, jobs <-chan int, results chan<- int) {
// 	for j := range jobs {
// 		// fmt.Println("worker", id, "started  job", j)
// 		time.Sleep(time.Second)
// 		fmt.Println("worker", id, "finished job", j)
// 		results <- j * 2
// 	}
// }

// func initJob(jobs chan <- int,_numJobs int){
// 	for j := 1; j <= _numJobs; j++ {
// 		jobs <- j
// 	}
// 	close(jobs)
// }

// func receiveResult(results <-chan int){
// 	for a := 1; a <= numJobs; a++ {
// 		// fmt.Println("result", a,"is", <-results)
// 		<-results
// 	}
// }

// func main() {
// 	jobs := make(chan int, numJobs)
// 	results := make(chan int, numJobs)

// 	for w :=1; w<= numWorker; w++ {
// 		go worker(w, jobs, results)
// 	}

// 	initJob(jobs, numJobs)

// 	receiveResult(results)


// }	
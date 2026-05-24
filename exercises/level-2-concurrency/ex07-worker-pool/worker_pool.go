package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Job struct {
	ID        int
	ImageName string
}

type Result struct {
	JobID    int
	Duration int
}

func Worker(id int, jobs <-chan Job, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()
	for job := range jobs {
		// Giả lập xử lý ảnh mất 100-500ms
		duration := rand.Intn(400) + 100
		time.Sleep(time.Duration(duration) * time.Millisecond)

		results <- Result{
			JobID:    job.ID,
			Duration: duration,
		}
		fmt.Printf("Worker #%d processed image_%d.jpg in %dms\n", id, job.ID, duration)
	}
}

func RunWorkerPool(numWorkers int, jobsList []Job) []Result {
	jobs := make(chan Job, len(jobsList))
	results := make(chan Result, len(jobsList))

	var wg sync.WaitGroup

	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go Worker(w, jobs, results, &wg)
	}

	for _, j := range jobsList {
		jobs <- j
	}
	close(jobs)

	// Đợi các worker hoàn thành
	wg.Wait()
	close(results)

	var output []Result
	for res := range results {
		output = append(output, res)
	}

	return output
}

func main() {
	rand.Seed(time.Now().UnixNano())

	jobs := make([]Job, 20)
	for i := 0; i < 20; i++ {
		jobs[i] = Job{ID: i + 1, ImageName: fmt.Sprintf("image_%d.jpg", i+1)}
	}

	start := time.Now()
	results := RunWorkerPool(5, jobs)
	duration := time.Since(start)

	fmt.Printf("\nAll %d images processed in %v\n", len(results), duration)
}

package main

import (
	"fmt"
	"sync"
)

// worker function takes an ID, a channel for jobs, and a channel for results.
// It processes each job received from the jobs channel and sends the processed result to the results channel.
func worker(id int, jobs <-chan int, results chan<- int) {
	for job := range jobs { // Loop over each job received from the jobs channel
		fmt.Printf("Worker %d processing job %d\n", id, job) // Log the job being processed by the worker
		results <- job * 2 // Process the job (in this case, doubling it) and send the result to the results channel
	}
}

func main() {
	numJobs := 10      // Define the total number of jobs that need to be processed
	numWorkers := 3    // Define the number of worker goroutines that will process the jobs

	// Create buffered channels for jobs and results, allowing for up to numJobs in each channel
	jobs := make(chan int, numJobs)    // Channel to send jobs to workers
	results := make(chan int, numJobs) // Channel to receive results from workers

	var wg sync.WaitGroup // Create a WaitGroup to manage synchronization of worker completion

	// Start worker goroutines
	for i := 1; i <= numWorkers; i++ {
		wg.Add(1) // Increment the WaitGroup counter for each worker being started
		go func(workerID int) {
			defer wg.Done() // Ensure the counter is decremented when the worker completes its execution
			worker(workerID, jobs, results) // Call the worker function with the worker's ID and the channels
		}(i) // Pass the current worker ID to the goroutine
	}

	// Enqueue jobs into the jobs channel
	for i := 1; i <= numJobs; i++ {
		jobs <- i // Send each job (an integer) to the jobs channel for processing
	}
	close(jobs) // Close the jobs channel to indicate that no more jobs will be sent

	// Start a goroutine to wait for all workers to finish and then close the results channel
	go func() {
		wg.Wait() // Block until all workers have finished processing
		close(results) // Close the results channel to signal that no more results will be sent
	}()

	// Collect results from the results channel
	for result := range results { // Loop over the results channel until it is closed
		fmt.Printf("Result: %d\n", result) // Print each result received from the results channel
	}
}

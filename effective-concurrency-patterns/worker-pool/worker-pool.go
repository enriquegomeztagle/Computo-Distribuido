package main

import (
	"fmt"
	"sync"
)

// Task represents a unit of work with an ID and a Result field.
type Task struct {
	ID     int // Unique identifier for the task
	Result int // Result of processing the task
}

// worker function processes tasks from the tasks channel and sends the results to the results channel.
// It takes a worker ID, channels for tasks and results, and a WaitGroup to signal completion.
func worker(id int, tasks <-chan Task, results chan<- Task, wg *sync.WaitGroup) {
	defer wg.Done() // Decrement the WaitGroup counter when the function completes, signaling that this worker is done

	// Process tasks received from the tasks channel
	for task := range tasks {
		// Simulate task processing
		// Here you can add any processing logic, such as computations or I/O operations.
		// For demonstration, we will simply double the task ID to simulate work.

		// Store the result in the task
		task.Result = task.ID * 2 // Process the task (e.g., double the ID) and store the result

		// Send the updated task to the results channel
		results <- task // Send the processed task with its result to the results channel
	}
}

func main() {
	numTasks := 20  // Define the total number of tasks to be processed
	numWorkers := 4 // Define the number of worker goroutines to be launched

	// Create a WaitGroup to synchronize the completion of worker goroutines
	var wg sync.WaitGroup

	// Create channels for tasks and results
	tasks := make(chan Task, numTasks)   // Channel to hold tasks for processing
	results := make(chan Task, numTasks) // Channel to hold results from workers

	// Launch worker goroutines
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)                         // Increment the WaitGroup counter for each worker being started
		go worker(i, tasks, results, &wg) // Start the worker goroutine, passing the worker ID, tasks channel, results channel, and WaitGroup
	}

	// Generate tasks and send them to the tasks channel
	for i := 0; i < numTasks; i++ {
		tasks <- Task{ID: i} // Create a new Task with a unique ID and send it to the tasks channel
	}

	// Close the tasks channel to signal that no more tasks will be added
	close(tasks) // This allows workers to finish processing once they have consumed all tasks

	// Wait for all workers to finish processing
	wg.Wait() // Block until all workers have completed their tasks

	// Close the results channel
	close(results) // Close the results channel to signal that no more results will be sent

	// Collect and process results from the results channel
	for result := range results { // Loop over the results channel until it is closed
		fmt.Printf("Processed result for task %d: %d\n", result.ID, result.Result) // Print the ID and result of each processed task
		// This output represents the final processed values from all workers, demonstrating the results of the concurrent processing
	}
}

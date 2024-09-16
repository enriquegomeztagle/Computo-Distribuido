package main

import (
	"context" // Import the context package for managing cancellation and timeouts
	"fmt"     // Import the fmt package for formatted I/O
	"sync"    // Import the sync package for synchronization primitives
	"time"    // Import the time package for time-related functions
)

// worker function simulates a worker that processes tasks.
// It takes a context for cancellation, a worker ID, and a WaitGroup to signal completion.
func worker(ctx context.Context, id int, wg *sync.WaitGroup) {
	defer wg.Done() // Decrement the WaitGroup counter when the function completes, signaling that this worker is done

	select {
	// Listen for cancellation signal from the context
	case <-ctx.Done():
		// If the context is canceled, print a cancellation message and exit the function
		fmt.Printf("Worker %d: Canceled\n", id)
		return

	// Simulate work by waiting for 1 second
	case <-time.After(time.Second):
		// If the work completes before cancellation, print a completion message
		fmt.Printf("Worker %d: Done\n", id)
	}
}

func main() {
	numWorkers := 4 // Define the number of worker goroutines to be launched

	// Create a context with a cancellation function
	ctx, cancel := context.WithCancel(context.Background())
	// The cancel function can be called to signal cancellation to all workers
	defer cancel() // Ensure that the cancel function is called when main exits

	// Create a WaitGroup to synchronize the completion of worker goroutines
	var wg sync.WaitGroup

	// Launch worker goroutines
	for i := 0; i < numWorkers; i++ {
		wg.Add(1) // Increment the WaitGroup counter for each worker being started
		go worker(ctx, i, &wg) // Start the worker goroutine, passing the context, worker ID, and WaitGroup
	}

	// Cancel the context after a brief delay
	go func() {
		time.Sleep(2 * time.Second) // Wait for 2 seconds before canceling
		cancel() // Call the cancel function to signal all workers to stop
	}()

	// Wait for all workers to finish processing
	wg.Wait() // Block until all workers have completed their tasks
}

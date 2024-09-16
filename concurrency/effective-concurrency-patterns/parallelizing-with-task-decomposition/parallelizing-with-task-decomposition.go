package main

import (
	"fmt"
	"sync"
)


// processItem is a function that simulates the processing of an individual item.
// It takes an integer item, a pointer to a WaitGroup to signal when processing is complete,
// and a channel to send the processed result back to the main function.
func processItem(item int, wg *sync.WaitGroup, results chan int) {
	defer wg.Done() // Decrement the WaitGroup counter when the function completes, signaling that this worker is done

	// Simulate item processing
	// In a real-world scenario, this could involve complex computations, I/O operations, or any other processing logic.
	// For demonstration purposes, we will simply double the value of the item.
	// This simulates a workload that each worker will handle concurrently.

	// Send the result to the channel
	results <- item * 2 // Send the processed result (item doubled) to the results channel for collection
}

func main() {
	numItems := 100    // Define the total number of items to be processed by the workers
	numWorkers := 4    // Define the number of worker goroutines that will be launched to process the items concurrently

	// Create a WaitGroup to synchronize the completion of worker goroutines
	var wg sync.WaitGroup // WaitGroup will help us wait for all workers to finish their tasks

	// Create a buffered channel to collect results from the workers
	results := make(chan int, numItems) // Buffered channel to hold results, sized to the number of items to avoid blocking

	// Launch worker goroutines
	for i := 0; i < numWorkers; i++ {
		wg.Add(1) // Increment the WaitGroup counter for each worker being started
		go processItem(i, &wg, results) // Start the worker goroutine, passing the worker index, WaitGroup, and results channel
		// Each worker will process an item and send the result back through the results channel
	}

	// Close the results channel when all workers are done
	go func() {
		wg.Wait() // Block the goroutine until all workers have finished processing their items
		close(results) // Close the results channel to signal that no more results will be sent
		// This is important to prevent deadlocks when reading from the results channel
	}()

	// Collect and process results from the results channel
	for result := range results { // Loop over the results channel until it is closed
		fmt.Printf("Processed result: %d\n", result) // Print each processed result received from the results channel
		// This output represents the final processed values from all workers, demonstrating the results of the concurrent processing
	}
}

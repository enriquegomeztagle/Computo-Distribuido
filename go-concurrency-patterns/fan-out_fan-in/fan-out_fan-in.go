package main

import (
	"fmt"
	"sync"
)

func main() {
	// Sample data to be processed
	data := []int{1, 2, 3, 4, 5} 
	input := make(chan int, len(data)) // Create a buffered channel to hold input data, sized to the length of the data slice

	// Populate the input channel with data
	for _, d := range data {
		input <- d // Send each data item to the input channel for processing
	}
	close(input) // Close the input channel to signal that no more data will be sent

	// Fan-out: Launch multiple worker goroutines to process the input data concurrently
	numWorkers := 3 // Define the number of worker goroutines to be launched
	results := make(chan int, len(data)) // Create a buffered channel to hold results from workers

	var wg sync.WaitGroup // Create a WaitGroup to synchronize the completion of worker goroutines

	// Launch worker goroutines
	for i := 0; i < numWorkers; i++ {
		wg.Add(1) // Increment the WaitGroup counter for each worker being started
		go func() {
			defer wg.Done() // Ensure the counter is decremented when the worker completes
			for num := range input { // Loop over each number received from the input channel
				// Simulate some processing (in this case, doubling the number)
				result := num * 2 // Process the number (e.g., double it)
				results <- result // Send the processed result to the results channel
			}
		}()
	}

	// Fan-in: Aggregate results from all worker goroutines
	go func() {
		wg.Wait() // Block until all workers have finished processing
		close(results) // Close the results channel to signal that no more results will be sent
	}()

	// Process aggregated results from the results channel
	for result := range results { // Loop over the results channel until it is closed
		fmt.Println(result) // Print each result received from the results channel
		// This output represents the final processed values from all workers
	}
}

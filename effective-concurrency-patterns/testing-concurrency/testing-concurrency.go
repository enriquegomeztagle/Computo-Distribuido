package main

import (
	"sync"

	"testing"
)

// ParallelFunction executes a parallel computation using goroutines.
// It creates a specified number of worker goroutines that each contribute
// to a shared result. The function returns the total sum of the worker IDs.
func ParallelFunction() int {

	var wg sync.WaitGroup // WaitGroup to manage synchronization of goroutines
	var result int        // Variable to store the cumulative result
	numWorkers := 4       // Number of worker goroutines to spawn

	wg.Add(numWorkers) // Set the number of goroutines to wait for

	// Launching worker goroutines
	for i := 0; i < numWorkers; i++ {
		go func(id int) {
			defer wg.Done() // Mark this goroutine as done when it finishes
			result += id    // Add the worker's ID to the result
		}(i) // Pass the current index as the worker ID
	}

	wg.Wait()     // Wait for all goroutines to complete
	return result // Return the final computed result
}

// TestParallelFunction tests the ParallelFunction to ensure it produces
// the expected result. It compares the output of ParallelFunction with
// the expected sum of worker IDs and reports any discrepancies.
func TestParallelFunction(t *testing.T) {

	expected := 6                // Expected sum of integers from 0 to 3 (0 + 1 + 2 + 3)
	result := ParallelFunction() // Call the function to test

	// Check if the result matches the expected value
	if result != expected {
		t.Errorf("Expected %d, but got %d", expected, result) // Report error if not
	}
}

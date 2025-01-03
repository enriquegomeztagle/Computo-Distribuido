package main

import (
	"fmt"
)

func main() {
	// Create the initial channel with some data
	data := []int{1, 2, 3, 4, 5} // Sample data to be processed through the pipeline
	input := make(chan int, len(data)) // Create a buffered channel to hold input data, sized to the length of the data slice

	// Populate the input channel with data
	for _, d := range data {
		input <- d // Send each data item to the input channel for processing
	}
	close(input) // Close the input channel to signal that no more data will be sent, allowing downstream consumers to know when to stop reading

	// First stage of the pipeline: Doubles the input values
	doubleOutput := make(chan int) // Create a channel to hold the results of the doubling operation
	go func() {
		defer close(doubleOutput) // Ensure the doubleOutput channel is closed when this goroutine completes
			// The goroutine will read from the input channel and process each value
		for num := range input { // Loop over each number received from the input channel
			// Perform the doubling operation
			doubledValue := num * 2 // Double the number
			doubleOutput <- doubledValue // Send the doubled value to the doubleOutput channel
			// At this point, the doubled value is ready for the next stage of processing
		}
		// Once all input values are processed, the doubleOutput channel is closed
	}()

	// Second stage of the pipeline: Squares the doubled values
	squareOutput := make(chan int) // Create a channel to hold the results of the squaring operation
	go func() {
		defer close(squareOutput) // Ensure the squareOutput channel is closed when this goroutine completes
			// The goroutine will read from the doubleOutput channel and process each value
		for num := range doubleOutput { // Loop over each number received from the doubleOutput channel
			// Perform the squaring operation
			squaredValue := num * num // Square the doubled value
			squareOutput <- squaredValue // Send the squared value to the squareOutput channel
			// At this point, the squared value is ready for final output
		}
		// Once all doubled values are processed, the squareOutput channel is closed
	}()

	// Third stage of the pipeline: Prints the squared values
	// This stage collects the final results from the squareOutput channel
	for result := range squareOutput { // Loop over the squareOutput channel until it is closed
		fmt.Println(result) // Print each squared result received from the squareOutput channel
		// This output represents the final processed values from the pipeline, showing the result of the squaring operation
		// The printed results will be the squares of the doubled input values
	}
}

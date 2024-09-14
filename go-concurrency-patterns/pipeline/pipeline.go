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
		for num := range input { // Loop over each number received from the input channel
			// Perform the doubling operation
			doubledValue := num * 2 // Double the number
			doubleOutput <- doubledValue // Send the doubled value to the doubleOutput channel
		}
	}()

	// Second stage of the pipeline: Squares the doubled values
	squareOutput := make(chan int) // Create a channel to hold the results of the squaring operation
	go func() {
		defer close(squareOutput) // Ensure the squareOutput channel is closed when this goroutine completes
		for num := range doubleOutput { // Loop over each number received from the doubleOutput channel
			// Perform the squaring operation
			squaredValue := num * num // Square the doubled value
			squareOutput <- squaredValue // Send the squared value to the squareOutput channel
		}
	}()

	// Third stage of the pipeline: Prints the squared values
	for result := range squareOutput { // Loop over the squareOutput channel until it is closed
		fmt.Println(result) // Print each squared result received from the squareOutput channel
		// This output represents the final processed values from the pipeline
	}
}

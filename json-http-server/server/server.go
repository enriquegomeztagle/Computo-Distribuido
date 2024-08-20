package server

import (
	"fmt"
	"net/http"

	"json-http-server/log"
)

// Start server
func StartServer() {
	// New Instance for our log struct
	log := &log.Log{}

	// Setting up HTTP handlers to produce and consume logs
	http.HandleFunc("/produce", log.Append)
	http.HandleFunc("/consume", log.Fetch)

	// Serve on port 8080
	fmt.Println("Server started on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		// Error to handle if server fails to start
		fmt.Println("Error starting the server:", err)
	}
}

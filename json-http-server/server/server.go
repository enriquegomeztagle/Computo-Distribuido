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

	// Setting up HTTP handlers to append and fetch logs
	http.HandleFunc("/append", log.Append)
	http.HandleFunc("/fetch", log.Fetch)

	// Serve on port 8080
	fmt.Println("Server started on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		// Error to handle if server fails to start
		fmt.Println("Error starting the server:", err)
	}
}

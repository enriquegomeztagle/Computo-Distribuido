package server

import (
	"fmt"
	"net/http"

	"json-http-server/log"

	"github.com/gorilla/mux"
)

// StartServer starts the HTTP server
func StartServer() {
	// New instance for our log struct
	log := &log.Log{}

	// Create a new Gorilla Mux router
	r := mux.NewRouter()

	// Setting up routes with Gorilla Mux
	r.HandleFunc("/append", log.Append).Methods("POST")
	r.HandleFunc("/fetch", log.Fetch).Methods("GET")

	// Serve on port 8080
	fmt.Println("Server started on port 8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		// Error handling if the server fails to start
		fmt.Println("Error starting the server:", err)
	}
}

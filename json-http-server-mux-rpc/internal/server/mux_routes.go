package server

import (
	"fmt"
	"json-http-server-mux-rpc/internal/log"
	"net/http"

	"github.com/gorilla/mux"
)

func StartHTTPServer() {
	logInstance := &log.Log{}

	r := mux.NewRouter()

	r.HandleFunc("/append", logInstance.Append).Methods("POST")
	r.HandleFunc("/fetch", logInstance.Fetch).Methods("GET")

	fmt.Println("HTTP Server started on port 8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Println("Error starting the HTTP server:", err)
	}
}

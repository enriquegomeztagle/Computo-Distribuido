package main

import (
	"json-http-server-mux-rpc/internal/server"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(2)

	// Start HTTP server
	go func() {
		defer wg.Done()
		server.StartHTTPServer()
	}()

	// Start RPC server
	go func() {
		defer wg.Done()
		server.StartRPCServer()
	}()

	wg.Wait()
}

package main

import (
	"flag"
	"fmt"
	"os"
)

// Entry point
func main() {
	// Parse CLI args to know in which mode to run
	mode := flag.String("mode", "server", "Mode to run: server or client")
	flag.Parse()

	switch *mode {
	case "server":
		fmt.Println("Running in server mode...")
		startServerMode()
	case "client":
		fmt.Println("Running in client mode...")
		startClientMode()
	default:
		fmt.Println("Unknown mode. Use either 'server' or 'client'.")
		os.Exit(1)
	}
}

func startServerMode() {
	if err := initializeRPCServer(); err != nil {
		fmt.Println("Failed to start the server:", err)
	}
}

func startClientMode() {
	if err := initializeRPCClient(); err != nil {
		fmt.Println("Failed to start the client:", err)
	}
}

func initializeRPCServer() error {
	fmt.Println("Starting the RPC server...")
	return nil
}

func initializeRPCClient() error {
	fmt.Println("Starting the RPC client...")
	return nil
}

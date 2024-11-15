package main

import (
	"log"
	"net"
	"os"
	"testing"

	api "github.com/enriquegomeztagle/log_server/api/v1"
	log2 "github.com/enriquegomeztagle/log_server/internal/log"
	"github.com/enriquegomeztagle/log_server/internal/server"

	"google.golang.org/grpc"
)

// TestMain sets up the testing environment for gRPC server and client.
func TestMain(t *testing.T) {
	// Define the directory for log data.
	dir := "./logdata"

	// Check if the log data directory exists.
	if _, err := os.Stat(dir); err == nil {
		// If it exists, remove all contents to start fresh.
		os.RemoveAll(dir)
	}

	// Create the log data directory if it does not exist.
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, 0755) // Create directory with appropriate permissions.
	}

	// Setup log configuration.
	config := log2.Config{}
	config.Segment.MaxStoreBytes = 1024 // Set maximum storage bytes for log segments.
	config.Segment.MaxIndexBytes = 1024 // Set maximum index bytes for log segments.

	// Create a new commit log with the specified directory and configuration.
	commitLog, err := log2.NewLog(dir, config)
	if err != nil {
		t.Fatalf("failed to create log: %v", err) // Fail the test if log creation fails.
	}

	// Setup gRPC server configuration with the commit log.
	serverConfig := &server.Config{
		CommitLog: commitLog,
	}

	// Create a new gRPC server with the specified configuration.
	grpcServer, err := server.NewGRPCServer(serverConfig)
	if err != nil {
		t.Fatalf("failed to create gRPC server: %v", err) // Fail the test if server creation fails.
	}

	// Start listening on TCP port 9000.
	listener, err := net.Listen("tcp", ":9000")
	if err != nil {
		t.Fatalf("failed to listen on port 9000: %v", err) // Fail the test if listening fails.
	}
	defer grpcServer.Stop() // Ensure the server stops when the test completes.

	// Start the gRPC server in a separate goroutine.
	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatalf("failed to serve gRPC server: %v", err) // Log fatal error if server fails to serve.
		}
	}()

	// Test gRPC client connection to the server.
	conn, err := grpc.Dial(":9000", grpc.WithInsecure())
	if err != nil {
		t.Fatalf("failed to dial gRPC server: %v", err) // Fail the test if client connection fails.
	}
	defer conn.Close() // Ensure the connection closes when the test completes.

	// Create a new log client using the established connection.
	client := api.NewLogClient(conn)

	// Define test cases for producing and consuming log records.
	tests := []struct {
		value string // The value to be produced and consumed.
	}{
		{"hello world"},   // Test case 1
		{"hello world 2"}, // Test case 2
		{"hello world 3"}, // Test case 3
	}

	// Iterate over each test case.
	for _, test := range tests {
		// Produce a log record using the client.
		offset, err := produce(client, test.value)
		if err != nil {
			t.Fatalf("failed to produce record: %v", err) // Fail the test if production fails.
		}

		// Consume the log record using the client.
		record, err := consume(client, offset)
		if err != nil {
			t.Fatalf("failed to consume record: %v", err) // Fail the test if consumption fails.
		}

		// Verify that the consumed record matches the expected value.
		if record != test.value {
			t.Errorf("expected %s, got %s", test.value, record) // Log an error if the values do not match.
		}
	}
}

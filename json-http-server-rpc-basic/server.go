package main

import (
	"fmt"
	"net"
	"net/rpc"
)

type Service struct{}

// Module with provided args
func (s *Service) Module(args *OperationArgs, result *int) error {
	*result = args.X % args.Y
	return nil
}

type OperationArgs struct {
	X, Y int
}

func main() {
	service := new(Service)

	// Register service (remote calls)
	rpc.Register(service)

	// Start listening on 3214 for RPC connections
	listener, err := net.Listen("tcp", ":3214")
	if err != nil {
		fmt.Println("Error starting RPC server:", err)
		return
	}

	fmt.Println("Server listening on port 3214...")
	for {
		// Accept incoming connections
		conn, err := listener.Accept()
		if err != nil {
			// If error -> Skip to next
			continue
		}
		// Serve conn with goroutines for concurrency
		go rpc.ServeConn(conn)
	}
}

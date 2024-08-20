package main

import (
	"fmt"
	"net/rpc"
)

// Args for RPC operation
type Args struct {
	X, Y int
}

func main() {
	// Connect to RPC server
	client, err := rpc.Dial("tcp", "server:3214")
	if err != nil {
		fmt.Println("Error connecting to RPC server:", err)
		return
	}

	// Define args for RPC call
	args := Args{X: 10, Y: 5}
	var result int

	// Call the Modulo method
	err = client.Call("Service.Module", &args, &result)
	if err != nil {
		fmt.Println("Error while calling RPC method:", err)
	} else {
		fmt.Printf("Result: %d %% %d = %d\n", args.X, args.Y, result)
	}
}

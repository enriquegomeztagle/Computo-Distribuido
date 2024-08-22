package server

import (
	"fmt"
	"json-http-server-mux-rpc/internal/log"
	"net"
	"net/rpc"
)

// Wrapper for Log to expose RPC methods
type LogRPC struct {
	Log *log.Log
}

// Adding record with RPC
func (l *LogRPC) Append(args *log.AppendArgs, reply *log.AppendReply) error {
	l.Log.Mu.Lock()
	defer l.Log.Mu.Unlock()

	args.Record.Offset = uint64(len(l.Log.Entries))
	l.Log.Entries = append(l.Log.Entries, args.Record)

	reply.Offset = args.Record.Offset
	return nil
}

// Fetching record with RPC
func (l *LogRPC) Fetch(args *log.FetchArgs, reply *log.FetchReply) error {
	l.Log.Mu.Lock()
	defer l.Log.Mu.Unlock()

	if args.Offset >= uint64(len(l.Log.Entries)) {
		return fmt.Errorf("offset out of range")
	}

	reply.Record = l.Log.Entries[args.Offset]
	return nil
}

func StartRPCServer() {
	logInstance := &log.Log{}

	// Register RPC service
	logRPC := &LogRPC{Log: logInstance}
	rpc.Register(logRPC)

	// Listen connections
	listener, err := net.Listen("tcp", ":3214")
	if err != nil {
		fmt.Println("Error starting the RPC server:", err)
		return
	}
	defer listener.Close()
	fmt.Println("RPC Server started on port 3214")

	// Accept connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go rpc.ServeConn(conn)
	}
}

package main

import (
    "context"
    "log"
    "net"

    "google.golang.org/grpc"
    pb "grpc-example/api/v1"
)

// server es la implementación del servicio ExampleService
type server struct {
    pb.UnimplementedExampleServiceServer
}

// Implementamos el método SayHello
func (s *server) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
    return &pb.HelloResponse{Message: "Hello, " + req.Name}, nil
}

func main() {
    lis, err := net.Listen("tcp", ":50051")
    if err != nil {
        log.Fatalf("Failed to listen: %v", err)
    }

    s := grpc.NewServer()
    pb.RegisterExampleServiceServer(s, &server{})

    log.Println("gRPC server is running on port 50051...")
    if err := s.Serve(lis); err != nil {
        log.Fatalf("Failed to serve: %v", err)
    }
}

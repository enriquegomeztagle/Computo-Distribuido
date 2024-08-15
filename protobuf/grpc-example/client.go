package main

import (
    "context"
    "log"
    "time"

    "google.golang.org/grpc"
    pb "grpc-example/api/v1"
)

func main() {
    conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure(), grpc.WithBlock())
    if err != nil {
        log.Fatalf("Failed to connect: %v", err)
    }
    defer conn.Close()

    client := pb.NewExampleServiceClient(conn)

    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    defer cancel()

    response, err := client.SayHello(ctx, &pb.HelloRequest{Name: "World"})
    if err != nil {
        log.Fatalf("Error calling SayHello: %v", err)
    }

    log.Printf("Response from server: %s", response.Message)
}

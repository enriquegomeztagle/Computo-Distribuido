package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	api "github.com/enriquegomeztagle/log_server/api/v1"
	log2 "github.com/enriquegomeztagle/log_server/internal/log"
	"github.com/enriquegomeztagle/log_server/internal/server"

	"google.golang.org/grpc"
)

func main() {
	dir := "./logdata"
	if _, err := os.Stat(dir); err == nil {
		os.RemoveAll(dir)
	}
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, 0755)
	}

	config := log2.Config{}
	config.Segment.MaxStoreBytes = 1024
	config.Segment.MaxIndexBytes = 1024

	commitLog, err := log2.NewLog(dir, config)
	if err != nil {
		log.Fatalf("failed to create log: %v", err)
	}

	serverConfig := &server.Config{
		CommitLog: commitLog,
	}

	grpcServer, err := server.NewGRPCServer(serverConfig)
	if err != nil {
		log.Fatalf("failed to create gRPC server: %v", err)
	}

	listener, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("failed to listen on port 9000: %v", err)
	}

	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatalf("failed to serve gRPC server: %v", err)
		}
	}()
	defer grpcServer.Stop()

	conn, err := grpc.Dial(":9000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to dial gRPC server: %v", err)
	}
	defer conn.Close()

	client := api.NewLogClient(conn)

	offset1, err := produce(client, "hello world")
	if err != nil {
		log.Fatalf("failed to produce record: %v", err)
	}
	fmt.Printf("Produced record at offset: %d\n", offset1)

	offset2, err := produce(client, "hello world 2")
	if err != nil {
		log.Fatalf("failed to produce record: %v", err)
	}
	fmt.Printf("Produced record at offset: %d\n", offset2)

	offset3, err := produce(client, "hello world 3")
	if err != nil {
		log.Fatalf("failed to produce record: %v", err)
	}
	fmt.Printf("Produced record at offset: %d\n", offset3)

	record, err := consume(client, offset1)
	if err != nil {
		log.Fatalf("failed to consume record: %v", err)
	}
	fmt.Printf("Consumed record: %s\n", record)

	record, err = consume(client, offset2)
	if err != nil {
		log.Fatalf("failed to consume record: %v", err)
	}
	fmt.Printf("Consumed record: %s\n", record)

	record, err = consume(client, offset3)
	if err != nil {
		log.Fatalf("failed to consume record: %v", err)
	}
	fmt.Printf("Consumed record: %s\n", record)
}

func produce(client api.LogClient, value string) (int64, error) {
	produceResponse, err := client.Produce(context.Background(), &api.ProduceRequest{
		Record: &api.Record{
			Value: []byte(value),
		},
	})
	if err != nil {
		return 0, err
	}
	return int64(produceResponse.Offset), nil
}

func consume(client api.LogClient, offset int64) (string, error) {
	consumeResponse, err := client.Consume(context.Background(), &api.ConsumeRequest{
		Offset: uint64(offset),
	})
	if err != nil {
		return "", err
	}
	return string(consumeResponse.Record.Value), nil
}

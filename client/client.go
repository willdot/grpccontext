package main

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	pb "github.com/willdot/grpccontext/server/proto"
)

const (
	address = "localhost:50051"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithUnaryInterceptor(clientInterceptor))
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewTestClient(conn)

	// make gRPC request to server to DoSomething
	response, err := c.DoSomething(context.Background(), &pb.Input{Id: 1})
	if err != nil {
		log.Fatalf("failed to DoSomething: %v", err)
	}

	log.Printf("response: %s", response.GetResult())

	// make gRPC request to server to RunLongTask
	_, err = c.RunLongTask(context.Background(), &pb.Input{Id: 2})
	if err != nil {
		log.Fatalf("failed to RunLongTask: %v", err)
	}
}

func clientInterceptor(ctx context.Context, method string, req interface{}, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	// create a context with a timeout that can be used to cancel the gRPC request
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	// create the headers
	headers := metadata.Pairs(
		"source", "client service")

	// add the headers to the context
	ctx = metadata.NewOutgoingContext(ctx, headers)

	// make request
	err := invoker(ctx, method, req, reply, cc, opts...)
	if err != nil {
		return err
	}

	return nil
}

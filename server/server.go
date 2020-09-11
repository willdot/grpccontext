package main

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/pkg/errors"
	pb "github.com/willdot/grpccontext/server/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
)

func main() {
	log.Println("server started")
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)
	pb.RegisterTestServer(s, &Server{})

	if err := s.Serve(lis); err != nil {
		log.Fatal(errors.Wrap(err, "failed to serve grpc server"))
	}
}

const (
	addr = "localhost:50051"
)

type Server struct {
}

func (s *Server) DoSomething(ctx context.Context, req *pb.Input) (*pb.Output, error) {
	log.Printf("Id: %v\n", req.Id)

	// get headers out of context
	md, _ := metadata.FromIncomingContext(ctx)
	if len(md["source"]) != 0 {
		source := md.Get("source")[0]

		log.Printf("request source: %s\n", source)
	}

	return &pb.Output{Result: "hello"}, nil
}

func (s *Server) RunLongTask(ctx context.Context, req *pb.Input) (*pb.Empty, error) {
	log.Printf("Id: %v\n", req.Id)

	deadline, ok := ctx.Deadline()
	if ok {
		log.Printf("deadline: %v\n", deadline)
	}

	// run a continuous loop until the context is canceled
	for {
		select {
		case <-ctx.Done():
			// this error won't be received by the client as it would have abandoned it's request and so won't get the response
			return nil, ctx.Err()
		default:
			time.Sleep(time.Second * 1)
		}
	}
}

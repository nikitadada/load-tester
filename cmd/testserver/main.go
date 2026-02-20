package main

import (
	"context"
	"log"
	"net"
	"time"

	pb "github.com/nikitadada/load-tester/internal/proto/gen"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedPingServiceServer
}

func (s *server) Ping(ctx context.Context, req *pb.PingRequest) (*pb.PingResponse, error) {
	time.Sleep(10 * time.Millisecond) // имитация обработки

	return &pb.PingResponse{
		Message: "pong",
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterPingServiceServer(grpcServer, &server{})

	log.Println("gRPC server started on :50051")
	grpcServer.Serve(lis)
}

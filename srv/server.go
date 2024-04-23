package biz

import (
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	pb "launcher/api/v1"
	"log"
	"net"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

type server struct {
	pb.UnimplementedLauncherServiceServer
}

func (s *server) Launch(ctx context.Context, in *pb.LaunchRequest) (*pb.LaunchResponse, error) {
	log.Printf("Received: %v", in.String())
	return &pb.LaunchResponse{}, nil
}

func StartGrpcServer() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterLauncherServiceServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

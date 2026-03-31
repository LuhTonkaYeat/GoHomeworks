package main

import (
	"log"
	"net"

	pb "github.com/LuhTonkaYeat/GoHomeworks/hw2/services/collector/api/proto"
	"github.com/LuhTonkaYeat/GoHomeworks/hw2/services/collector/internal/adapter/github"
	grpcHandler "github.com/LuhTonkaYeat/GoHomeworks/hw2/services/collector/internal/adapter/grpc"
	"github.com/LuhTonkaYeat/GoHomeworks/hw2/services/collector/internal/usecase"
	"google.golang.org/grpc"
)

func main() {
	githubClient := github.NewClient()
	repoUseCase := usecase.NewRepositoryUseCase(githubClient)
	grpcServer := grpcHandler.NewServer(repoUseCase)

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	server := grpc.NewServer()
	pb.RegisterCollectorServiceServer(server, grpcServer)

	log.Println("Collector server started on :50051")
	log.Println("Waiting for gRPC requests...")

	if err := server.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

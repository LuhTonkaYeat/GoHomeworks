package main

import (
	"log"
	"net"
	"os"

	pb "github.com/LuhTonkaYeat/GoHomeworks/services/collector/api/proto/collector"
	"github.com/LuhTonkaYeat/GoHomeworks/services/collector/internal/adapter/github"
	"github.com/LuhTonkaYeat/GoHomeworks/services/collector/internal/adapter/subscribe"
	grpcHandler "github.com/LuhTonkaYeat/GoHomeworks/services/collector/internal/delivery/grpc"
	"github.com/LuhTonkaYeat/GoHomeworks/services/collector/internal/usecase"
	"google.golang.org/grpc"
)

func main() {
	subscribeAddr := os.Getenv("SUBSCRIBE_ADDR")
	if subscribeAddr == "" {
		subscribeAddr = "subscribe:50053"
		log.Printf("SUBSCRIBE_ADDR not set, using default: %s", subscribeAddr)
	}

	subscribeClient, err := subscribe.NewClient(subscribeAddr)
	if err != nil {
		log.Fatalf("Failed to create Subscribe client: %v", err)
	}
	defer subscribeClient.Close()

	githubClient := github.NewClient()

	repoUseCase := usecase.NewRepositoryUseCase(githubClient, subscribeClient)

	grpcServer := grpcHandler.NewServer(repoUseCase)

	port := os.Getenv("PORT")
	if port == "" {
		port = "50051"
	}

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	server := grpc.NewServer()
	pb.RegisterCollectorServiceServer(server, grpcServer)

	log.Printf("Collector server started on port %s", port)
	log.Printf("Subscribe address: %s", subscribeAddr)
	log.Println("Waiting for gRPC requests...")

	if err := server.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

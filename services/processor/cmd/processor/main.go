package main

import (
	"log"
	"net"
	"os"

	pb "github.com/LuhTonkaYeat/GoHomeworks/services/processor/api/proto/processor"
	collectorclient "github.com/LuhTonkaYeat/GoHomeworks/services/processor/internal/adapter/grpc"
	deliverygrpc "github.com/LuhTonkaYeat/GoHomeworks/services/processor/internal/delivery/grpc"
	"github.com/LuhTonkaYeat/GoHomeworks/services/processor/internal/usecase"
	"google.golang.org/grpc"
)

func main() {
	collectorAddr := os.Getenv("COLLECTOR_ADDR")
	if collectorAddr == "" {
		collectorAddr = "collector:50051"
	}

	collectorClient, err := collectorclient.NewClient(collectorAddr)
	if err != nil {
		log.Fatalf("Failed to create collector client: %v", err)
	}
	defer collectorClient.Close()

	repoUseCase := usecase.NewRepositoryUseCase(collectorClient, nil)

	grpcHandler := deliverygrpc.NewHandler(repoUseCase, collectorClient)

	port := os.Getenv("PORT")
	if port == "" {
		port = "50052"
	}

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	server := grpc.NewServer()
	pb.RegisterProcessorServiceServer(server, grpcHandler)

	log.Printf("Processor server started on port %s", port)
	log.Printf("Connected to Collector at %s", collectorAddr)

	if err := server.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

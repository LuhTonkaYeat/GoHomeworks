package main

import (
	"log"
	"net"

	pb "github.com/LuhTonkaYeat/GoHomeworks/hw2/services/processor/api/proto"
	collectorclient "github.com/LuhTonkaYeat/GoHomeworks/hw2/services/processor/internal/adapter/grpc"
	deliverygrpc "github.com/LuhTonkaYeat/GoHomeworks/hw2/services/processor/internal/delivery/grpc"
	"github.com/LuhTonkaYeat/GoHomeworks/hw2/services/processor/internal/usecase"
	"google.golang.org/grpc"
)

func main() {
	collectorAddr := "collector:50051"
	subscriberAddr := "subscriber:50052"

	collectorClient, err := collectorclient.NewClient(collectorAddr)
	if err != nil {
		log.Fatalf("Failed to create collector client: %v", err)
	}
	defer collectorClient.Close()

	subscriberClient, err := collectorclient.NewSubscriberClient(subscriberAddr)
	if err != nil {
		log.Fatalf("Failed to create subscriber client: %v", err)
	}
	defer subscriberClient.Close()

	repoUseCase := usecase.NewRepositoryUseCase(collectorClient, subscriberClient)

	grpcHandler := deliverygrpc.NewHandler(repoUseCase)

	listener, err := net.Listen("tcp", ":50053")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	server := grpc.NewServer()
	pb.RegisterProcessorServiceServer(server, grpcHandler)

	log.Println("Processor server started on :50053")
	log.Printf("Connected to Collector at %s", collectorAddr)
	log.Printf("Connected to Subscriber at %s", subscriberAddr)

	if err := server.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

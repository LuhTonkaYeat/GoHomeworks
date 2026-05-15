package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"strings"

	pb "github.com/LuhTonkaYeat/GoHomeworks/services/processor/api/proto/processor"
	"github.com/LuhTonkaYeat/GoHomeworks/services/processor/internal/adapter/kafka"
	"github.com/LuhTonkaYeat/GoHomeworks/services/processor/internal/adapter/postgres"
	deliverygrpc "github.com/LuhTonkaYeat/GoHomeworks/services/processor/internal/delivery/grpc"
	"github.com/LuhTonkaYeat/GoHomeworks/services/processor/internal/repository"
	"github.com/LuhTonkaYeat/GoHomeworks/services/processor/internal/usecase"
	"google.golang.org/grpc"
)

func main() {
	dbURL := os.Getenv("PROCESSOR_DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5433/processor?sslmode=disable"
	}
	
	dbRepo, err := postgres.NewRepository(dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbRepo.Close()
	
	queries := dbRepo.GetQueries()
	
	kafkaBrokers := strings.Split(os.Getenv("KAFKA_BROKERS"), ",")
	if len(kafkaBrokers) == 0 || kafkaBrokers[0] == "" {
		kafkaBrokers = []string{"localhost:9092"}
	}
	
	requestsTopic := os.Getenv("KAFKA_REQUESTS_TOPIC")
	if requestsTopic == "" {
		requestsTopic = "repository.fetch.requests"
	}
	
	responsesTopic := os.Getenv("KAFKA_RESPONSES_TOPIC")
	if responsesTopic == "" {
		responsesTopic = "repository.fetch.responses"
	}
	
	kafkaProducer := kafka.NewProducer(kafkaBrokers, requestsTopic)
	defer kafkaProducer.Close()
	
	githubClient := NewGitHubClient()
	
	repoUseCase := usecase.NewRepositoryUseCase(queries, kafkaProducer, githubClient)
	
	kafkaConsumer := kafka.NewConsumer(kafkaBrokers, responsesTopic, "processor-group", func(response kafka.RepositoryFetchResponse) error {
		return repoUseCase.HandleKafkaResponse(response)
	})
	
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	kafkaConsumer.Start(ctx)
	defer kafkaConsumer.Close()
	
	grpcHandler := deliverygrpc.NewHandler(repoUseCase)
	
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
	
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	
	go func() {
		log.Printf("Processor server started on port %s", port)
		if err := server.Serve(listener); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()
	
	<-sigChan
	log.Println("Shutting down...")
	server.GracefulStop()
}

type GitHubClient struct{}

func NewGitHubClient() *GitHubClient {
	return &GitHubClient{}
}

func (c *GitHubClient) FetchRepository(ctx context.Context, owner, repo string) (*usecase.GitHubRepository, error) {

	return nil, nil
}
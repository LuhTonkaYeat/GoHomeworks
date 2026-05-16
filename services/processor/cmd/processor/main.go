package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	pb "github.com/LuhTonkaYeat/GoHomeworks/services/processor/api/proto/processor"
	"github.com/LuhTonkaYeat/GoHomeworks/services/processor/internal/adapter/kafka"
	"github.com/LuhTonkaYeat/GoHomeworks/services/processor/internal/adapter/postgres"
	deliverygrpc "github.com/LuhTonkaYeat/GoHomeworks/services/processor/internal/delivery/grpc"
	"github.com/LuhTonkaYeat/GoHomeworks/services/processor/internal/domain"
	"github.com/LuhTonkaYeat/GoHomeworks/services/processor/internal/usecase"
	"google.golang.org/grpc"
)

type GitHubClient struct {
	httpClient *http.Client
	baseURL    string
}

func NewGitHubClient() *GitHubClient {
	return &GitHubClient{
		httpClient: &http.Client{Timeout: 10 * time.Second},
		baseURL:    "https://api.github.com",
	}
}

func (c *GitHubClient) FetchRepository(ctx context.Context, owner, repo string) (*domain.Repository, error) {
	url := fmt.Sprintf("%s/repos/%s/%s", c.baseURL, owner, repo)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			return nil, fmt.Errorf("repository %s/%s not found", owner, repo)
		}
		return nil, fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	var ghRepo struct {
		Name        string `json:"name"`
		FullName    string `json:"full_name"`
		Description string `json:"description"`
		Stars       int    `json:"stargazers_count"`
		Forks       int    `json:"forks_count"`
		CreatedAt   string `json:"created_at"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&ghRepo); err != nil {
		return nil, err
	}

	createdAt, err := time.Parse(time.RFC3339, ghRepo.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &domain.Repository{
		Name:        ghRepo.FullName,
		Description: ghRepo.Description,
		Stars:       ghRepo.Stars,
		Forks:       ghRepo.Forks,
		CreatedAt:   createdAt,
	}, nil
}

func main() {
	dbURL := os.Getenv("PROCESSOR_DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5433/processor?sslmode=disable"
	}
	log.Printf("Connecting to PostgreSQL at: %s", dbURL)

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
		log.Printf("Received Kafka response for %s/%s", response.Owner, response.Repo)
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

	grpcServer := grpc.NewServer()
	pb.RegisterProcessorServiceServer(grpcServer, grpcHandler)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("Processor gRPC server started on port %s", port)
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	<-sigChan
	log.Println("Shutting down...")
	grpcServer.GracefulStop()
	cancel()
}
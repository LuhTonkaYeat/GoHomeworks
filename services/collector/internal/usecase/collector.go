package usecase

import (
	"context"
	"log"
	"time"

	"github.com/LuhTonkaYeat/GoHomeworks/services/collector/internal/adapter/kafka"
	"github.com/LuhTonkaYeat/GoHomeworks/services/collector/internal/domain"
)

type CollectorUseCase interface {
	ProcessFetchRequest(ctx context.Context, request kafka.RepositoryFetchRequest) error
	RefreshSubscriptions(ctx context.Context) error
}

type GitHubClient interface {
	FetchRepository(ctx context.Context, owner, repo string) (*domain.Repository, error)
}

type SubscribeClient interface {
	GetSubscriptions(ctx context.Context) ([]*SubscriptionRepo, error)
}

type SubscriptionRepo struct {
	Owner string
	Name  string
}

type collectorUseCase struct {
	githubClient    GitHubClient
	subscribeClient SubscribeClient
	kafkaProducer   *kafka.Producer
}

func NewCollectorUseCase(
	githubClient GitHubClient,
	subscribeClient SubscribeClient,
	kafkaProducer *kafka.Producer,
) CollectorUseCase {
	return &collectorUseCase{
		githubClient:    githubClient,
		subscribeClient: subscribeClient,
		kafkaProducer:   kafkaProducer,
	}
}

func (uc *collectorUseCase) ProcessFetchRequest(ctx context.Context, request kafka.RepositoryFetchRequest) error {
	log.Printf("Processing fetch request for %s/%s (request_id: %s)", request.Owner, request.Repo, request.RequestID)
	
	repo, err := uc.githubClient.FetchRepository(ctx, request.Owner, request.Repo)
	
	response := kafka.RepositoryFetchResponse{
		RequestID: request.RequestID,
		Owner:     request.Owner,
		Repo:      request.Repo,
	}
	
	if err != nil {
		log.Printf("Error fetching from GitHub: %v", err)
		response.Error = err.Error()
	} else {
		response.FullName = repo.Name
		response.Description = repo.Description
		response.Stars = repo.Stars
		response.Forks = repo.Forks
		response.CreatedAt = repo.CreatedAt.Format(time.RFC3339)
	}
	
	err = uc.kafkaProducer.SendFetchResponse(ctx, response)
	if err != nil {
		log.Printf("Failed to send response to Kafka: %v", err)
		return err
	}
	
	log.Printf("Sent response for %s/%s (request_id: %s)", request.Owner, request.Repo, request.RequestID)
	return nil
}

func (uc *collectorUseCase) RefreshSubscriptions(ctx context.Context) error {
	log.Println("Refreshing subscriptions...")
	
	subscriptions, err := uc.subscribeClient.GetSubscriptions(ctx)
	if err != nil {
		log.Printf("Failed to get subscriptions: %v", err)
		return err
	}
	
	log.Printf("Found %d subscriptions", len(subscriptions))
	
	for _, sub := range subscriptions {
		requestID := generateRequestID()
		request := kafka.RepositoryFetchRequest{
			RequestID: requestID,
			Owner:     sub.Owner,
			Repo:      sub.Name,
		}
		
		err = uc.kafkaProducer.SendFetchRequest(ctx, request)
		if err != nil {
			log.Printf("Failed to send fetch request for %s/%s: %v", sub.Owner, sub.Name, err)
			continue
		}
		
		log.Printf("Sent refresh request for %s/%s", sub.Owner, sub.Name)
	}
	
	return nil
}

func generateRequestID() string {
	b := make([]byte, 8)
	for i := range b {
		b[i] = byte(65 + i%26)
	}
	return string(b)
}
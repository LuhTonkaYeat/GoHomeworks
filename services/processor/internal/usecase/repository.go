package usecase

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"log"
	"time"

	"github.com/LuhTonkaYeat/GoHomeworks/services/processor/internal/adapter/kafka"
	"github.com/LuhTonkaYeat/GoHomeworks/services/processor/internal/domain"
	"github.com/LuhTonkaYeat/GoHomeworks/services/processor/internal/repository"
)

type RepositoryUseCase interface {
	GetRepository(ctx context.Context, owner, repo string) (*domain.Repository, error)
	HandleKafkaResponse(response kafka.RepositoryFetchResponse) error
}

type GitHubClient interface {
	FetchRepository(ctx context.Context, owner, repo string) (*domain.Repository, error)
}

type repositoryUseCase struct {
	dbRepo        *repository.Queries
	kafkaProducer *kafka.Producer
	githubClient  GitHubClient
}

func NewRepositoryUseCase(dbRepo *repository.Queries, kafkaProducer *kafka.Producer, githubClient GitHubClient) RepositoryUseCase {
	return &repositoryUseCase{
		dbRepo:        dbRepo,
		kafkaProducer: kafkaProducer,
		githubClient:  githubClient,
	}
}

func (uc *repositoryUseCase) GetRepository(ctx context.Context, owner, repo string) (*domain.Repository, error) {
	dbRepo, err := uc.dbRepo.GetRepositoryByOwnerAndRepo(ctx, repository.GetRepositoryByOwnerAndRepoParams{
		Owner: owner,
		Repo:  repo,
	})

	if err == nil {
		log.Printf("Cache HIT for %s/%s", owner, repo)
		return &domain.Repository{
			Name:        dbRepo.FullName,
			Description: dbRepo.Description,
			Stars:       int(dbRepo.Stars),
			Forks:       int(dbRepo.Forks),
			CreatedAt:   dbRepo.CreatedAt,
		}, nil
	}

	log.Printf("Cache MISS for %s/%s, sending to Kafka", owner, repo)

	requestID := generateRequestID()
	err = uc.kafkaProducer.SendFetchRequest(ctx, requestID, owner, repo)
	if err != nil {
		log.Printf("Failed to send Kafka request: %v", err)
		return uc.githubClient.FetchRepository(ctx, owner, repo)
	}

	return nil, &NotFoundError{owner: owner, repo: repo}
}

func (uc *repositoryUseCase) HandleKafkaResponse(response kafka.RepositoryFetchResponse) error {
	ctx := context.Background()

	if response.Error != "" {
		log.Printf("Error from Collector for %s/%s: %s", response.Owner, response.Repo, response.Error)
		return nil
	}

	createdAt, err := time.Parse(time.RFC3339, response.CreatedAt)
	if err != nil {
		log.Printf("Failed to parse created_at: %v", err)
		return err
	}

	err = uc.dbRepo.UpsertRepository(ctx, repository.UpsertRepositoryParams{
		Owner:       response.Owner,
		Repo:        response.Repo,
		FullName:    response.FullName,
		Description: response.Description,
		Stars:       int32(response.Stars),
		Forks:       int32(response.Forks),
		CreatedAt:   createdAt,
	})

	if err != nil {
		log.Printf("Failed to upsert repository: %v", err)
		return err
	}

	log.Printf("Saved repository %s/%s to DB", response.Owner, response.Repo)
	return nil
}

type NotFoundError struct {
	owner string
	repo  string
}

func (e *NotFoundError) Error() string {
	return "repository " + e.owner + "/" + e.repo + " not found in cache, request sent to Kafka"
}

func generateRequestID() string {
	b := make([]byte, 8)
	rand.Read(b)
	return hex.EncodeToString(b)
}
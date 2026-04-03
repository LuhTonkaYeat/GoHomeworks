package usecase

import (
	"context"
	"fmt"

	"github.com/LuhTonkaYeat/GoHomeworks/hw2/services/processor/internal/domain"
)

type RepositoryUseCase interface {
	GetRepository(ctx context.Context, owner, repo string) (*domain.Repository, error)
	Ping(ctx context.Context) (string, error)
}

type collectorClient interface {
	GetRepository(ctx context.Context, owner, repo string) (*domain.Repository, error)
}

type subscriberClient interface {
	Ping(ctx context.Context) error
}

type repositoryUseCase struct {
	collectorClient  collectorClient
	subscriberClient subscriberClient
}

func NewRepositoryUseCase(collectorClient collectorClient, subscriberClient subscriberClient) RepositoryUseCase {
	return &repositoryUseCase{
		collectorClient:  collectorClient,
		subscriberClient: subscriberClient,
	}
}

func (uc *repositoryUseCase) GetRepository(ctx context.Context, owner, repo string) (*domain.Repository, error) {
	if owner == "" || repo == "" {
		return nil, fmt.Errorf("owner and repo are required")
	}

	return uc.collectorClient.GetRepository(ctx, owner, repo)
}

func (uc *repositoryUseCase) Ping(ctx context.Context) (string, error) {
	// Проверяем subscriber
	if uc.subscriberClient != nil {
		if err := uc.subscriberClient.Ping(ctx); err != nil {
			return "down", nil
		}
	}
	return "up", nil
}

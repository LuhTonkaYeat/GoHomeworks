package usecase

import (
	"context"
	"fmt"

	"github.com/LuhTonkaYeat/GoHomeworks/hw2/services/gateway/internal/domain"
)

type RepositoryUseCase interface {
	GetRepository(ctx context.Context, owner, repo string) (*domain.Repository, error)
}

type CollectorClient interface {
	GetRepository(ctx context.Context, owner, repo string) (*domain.Repository, error)
}

type repositoryUseCase struct {
	collectorClient CollectorClient
}

func NewRepositoryUseCase(collectorClient CollectorClient) RepositoryUseCase {
	return &repositoryUseCase{
		collectorClient: collectorClient,
	}
}

func (uc *repositoryUseCase) GetRepository(ctx context.Context, owner, repo string) (*domain.Repository, error) {
	if owner == "" || repo == "" {
		return nil, fmt.Errorf("owner and repo are required")
	}

	return uc.collectorClient.GetRepository(ctx, owner, repo)
}

package usecase

import (
	"context"
	"fmt"

	"github.com/LuhTonkaYeat/GoHomeworks/hw2/services/collector/internal/domain"
)

type RepositoryUseCase interface {
	GetRepository(ctx context.Context, owner, repo string) (*domain.Repository, error)
}

type GitHubClient interface {
	FetchRepository(ctx context.Context, owner, repo string) (*domain.Repository, error)
}

type repositoryUseCase struct {
	githubClient GitHubClient
}

func NewRepositoryUseCase(githubClient GitHubClient) RepositoryUseCase {
	return &repositoryUseCase{
		githubClient: githubClient,
	}
}

func (uc *repositoryUseCase) GetRepository(ctx context.Context, owner, repo string) (*domain.Repository, error) {
	if owner == "" || repo == "" {
		return nil, fmt.Errorf("owner and repo are required")
	}

	return uc.githubClient.FetchRepository(ctx, owner, repo)
}

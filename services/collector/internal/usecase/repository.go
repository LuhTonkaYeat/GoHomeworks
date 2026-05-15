package usecase

import (
	"context"
	"fmt"

	"github.com/LuhTonkaYeat/GoHomeworks/services/collector/internal/domain"
)

type RepositoryUseCase interface {
	GetRepository(ctx context.Context, owner, repo string) (*domain.Repository, error)
	GetSubscriptionsInfo(ctx context.Context) ([]*domain.Repository, error)
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

type repositoryUseCase struct {
	githubClient    GitHubClient
	subscribeClient SubscribeClient
}

func NewRepositoryUseCase(githubClient GitHubClient, subscribeClient SubscribeClient) RepositoryUseCase {
	return &repositoryUseCase{
		githubClient:    githubClient,
		subscribeClient: subscribeClient,
	}
}

func (uc *repositoryUseCase) GetRepository(ctx context.Context, owner, repo string) (*domain.Repository, error) {
	if owner == "" || repo == "" {
		return nil, fmt.Errorf("owner and repo are required")
	}

	return uc.githubClient.FetchRepository(ctx, owner, repo)
}

func (uc *repositoryUseCase) GetSubscriptionsInfo(ctx context.Context) ([]*domain.Repository, error) {
	subscriptions, err := uc.subscribeClient.GetSubscriptions(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get subscriptions: %w", err)
	}

	if len(subscriptions) == 0 {
		return []*domain.Repository{}, nil
	}

	var repositories []*domain.Repository
	for _, sub := range subscriptions {
		repo, err := uc.githubClient.FetchRepository(ctx, sub.Owner, sub.Name)
		if err != nil {
			continue
		}
		repositories = append(repositories, repo)
	}

	return repositories, nil
}

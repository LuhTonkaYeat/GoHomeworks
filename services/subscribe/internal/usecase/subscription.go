package usecase

import (
	"context"
	"fmt"

	"github.com/LuhTonkaYeat/GoHomeworks/services/subscribe/internal/adapter/github"
	"github.com/LuhTonkaYeat/GoHomeworks/services/subscribe/internal/repository"
)

type SubscriptionUseCase interface {
	CreateSubscription(ctx context.Context, userID, owner, repo string) error
	DeleteSubscription(ctx context.Context, userID, owner, repo string) error
	GetSubscriptions(ctx context.Context, userID string) ([]repository.GetSubscriptionsRow, error)
	CheckRepoExists(ctx context.Context, owner, repo string) (bool, error)
}

type subscriptionUseCase struct {
	repo   *repository.Queries
	github *github.Client
}

func NewSubscriptionUseCase(repo *repository.Queries, githubClient *github.Client) SubscriptionUseCase {
	return &subscriptionUseCase{
		repo:   repo,
		github: githubClient,
	}
}

func (uc *subscriptionUseCase) CreateSubscription(ctx context.Context, userID, owner, repo string) error {
	exists, err := uc.github.CheckRepoExists(ctx, owner, repo)
	if err != nil {
		return fmt.Errorf("failed to check repository: %w", err)
	}
	if !exists {
		return fmt.Errorf("repository %s/%s does not exist on GitHub", owner, repo)
	}

	err = uc.repo.CreateSubscription(ctx, repository.CreateSubscriptionParams{
		UserID: userID,
		Owner:  owner,
		Repo:   repo,
	})
	if err != nil {
		return fmt.Errorf("failed to create subscription: %w", err)
	}

	return nil
}

func (uc *subscriptionUseCase) DeleteSubscription(ctx context.Context, userID, owner, repo string) error {
	return uc.repo.DeleteSubscription(ctx, repository.DeleteSubscriptionParams{
		UserID: userID,
		Owner:  owner,
		Repo:   repo,
	})
}

func (uc *subscriptionUseCase) GetSubscriptions(ctx context.Context, userID string) ([]repository.GetSubscriptionsRow, error) {
	return uc.repo.GetSubscriptions(ctx, userID)
}

func (uc *subscriptionUseCase) CheckRepoExists(ctx context.Context, owner, repo string) (bool, error) {
	return uc.github.CheckRepoExists(ctx, owner, repo)
}

package usecase

import (
	"context"

	processorPb "github.com/LuhTonkaYeat/GoHomeworks/services/gateway/api/proto/processor"
	subscribePb "github.com/LuhTonkaYeat/GoHomeworks/services/gateway/api/proto/subscribe"
	"github.com/LuhTonkaYeat/GoHomeworks/services/gateway/internal/adapter/grpc"
)

type SubscriptionUseCase interface {
	CreateSubscription(ctx context.Context, owner, repo, userID string) error
	DeleteSubscription(ctx context.Context, owner, repo, userID string) error
	GetSubscriptions(ctx context.Context, userID string) ([]*subscribePb.Repository, error)
	GetSubscriptionsInfo(ctx context.Context, userID string) ([]*processorPb.RepoResponse, error)
}

type subscriptionUseCase struct {
	subscribeClient *grpc.SubscribeClient
	processorClient *grpc.Client
}

func NewSubscriptionUseCase(subscribeClient *grpc.SubscribeClient, processorClient *grpc.Client) SubscriptionUseCase {
	return &subscriptionUseCase{
		subscribeClient: subscribeClient,
		processorClient: processorClient,
	}
}

func (uc *subscriptionUseCase) CreateSubscription(ctx context.Context, owner, repo, userID string) error {
	return uc.subscribeClient.CreateSubscription(ctx, owner, repo, userID)
}

func (uc *subscriptionUseCase) DeleteSubscription(ctx context.Context, owner, repo, userID string) error {
	return uc.subscribeClient.DeleteSubscription(ctx, owner, repo, userID)
}

func (uc *subscriptionUseCase) GetSubscriptions(ctx context.Context, userID string) ([]*subscribePb.Repository, error) {
	return uc.subscribeClient.GetSubscriptions(ctx, userID)
}

func (uc *subscriptionUseCase) GetSubscriptionsInfo(ctx context.Context, userID string) ([]*processorPb.RepoResponse, error) {
	return uc.processorClient.GetSubscriptionsInfo(ctx)
}

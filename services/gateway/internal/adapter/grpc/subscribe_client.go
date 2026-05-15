package grpc

import (
	"context"

	subscribePb "github.com/LuhTonkaYeat/GoHomeworks/services/gateway/api/proto/subscribe"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type SubscribeClient struct {
	client subscribePb.SubscribeServiceClient
	conn   *grpc.ClientConn
}

func NewSubscribeClient(addr string) (*SubscribeClient, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := subscribePb.NewSubscribeServiceClient(conn)

	return &SubscribeClient{
		client: client,
		conn:   conn,
	}, nil
}

func (c *SubscribeClient) Close() error {
	return c.conn.Close()
}

func (c *SubscribeClient) CreateSubscription(ctx context.Context, owner, repo, userID string) error {
	req := &subscribePb.CreateSubscriptionRequest{
		Owner:  owner,
		Repo:   repo,
		UserId: userID,
	}

	_, err := c.client.CreateSubscription(ctx, req)
	return err
}

func (c *SubscribeClient) DeleteSubscription(ctx context.Context, owner, repo, userID string) error {
	req := &subscribePb.DeleteSubscriptionRequest{
		Owner:  owner,
		Repo:   repo,
		UserId: userID,
	}

	_, err := c.client.DeleteSubscription(ctx, req)
	return err
}

func (c *SubscribeClient) GetSubscriptions(ctx context.Context, userID string) ([]*subscribePb.Repository, error) {
	req := &subscribePb.Empty{}

	resp, err := c.client.GetSubscriptions(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp.Subscriptions, nil
}

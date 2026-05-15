package subscribe

import (
	"context"

	subscribePb "github.com/LuhTonkaYeat/GoHomeworks/services/collector/api/proto/subscribe"
	"github.com/LuhTonkaYeat/GoHomeworks/services/collector/internal/usecase"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	grpcClient subscribePb.SubscribeServiceClient
	conn       *grpc.ClientConn
}

func NewClient(addr string) (*Client, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	grpcClient := subscribePb.NewSubscribeServiceClient(conn)

	return &Client{
		grpcClient: grpcClient,
		conn:       conn,
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) GetSubscriptions(ctx context.Context) ([]*usecase.SubscriptionRepo, error) {
	req := &subscribePb.Empty{}

	resp, err := c.grpcClient.GetSubscriptions(ctx, req)
	if err != nil {
		return nil, err
	}

	var subscriptions []*usecase.SubscriptionRepo
	for _, sub := range resp.Subscriptions {
		subscriptions = append(subscriptions, &usecase.SubscriptionRepo{
			Owner: sub.Owner,
			Name:  sub.Name,
		})
	}

	return subscriptions, nil
}

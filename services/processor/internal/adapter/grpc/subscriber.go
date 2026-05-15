package grpc

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type SubscriberClient struct {
	conn *grpc.ClientConn
}

func NewSubscriberClient(addr string) (*SubscriberClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to subscriber: %w", err)
	}

	return &SubscriberClient{
		conn: conn,
	}, nil
}

func (c *SubscriberClient) Close() error {
	return c.conn.Close()
}

func (c *SubscriberClient) Ping(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return nil
}

package grpc

import (
	"context"
	"fmt"
	"time"

	processorPb "github.com/LuhTonkaYeat/GoHomeworks/services/gateway/api/proto/processor"
	"github.com/LuhTonkaYeat/GoHomeworks/services/gateway/internal/domain"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type Client struct {
	conn   *grpc.ClientConn
	client processorPb.ProcessorServiceClient
}

func NewClient(addr string) (*Client, error) {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to processor: %w", err)
	}

	return &Client{
		conn:   conn,
		client: processorPb.NewProcessorServiceClient(conn),
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) GetRepository(ctx context.Context, owner, repo string) (*domain.Repository, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	resp, err := c.client.GetRepository(ctx, &processorPb.RepoRequest{
		Owner: owner,
		Repo:  repo,
	})

	if err != nil {
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.NotFound:
				return nil, fmt.Errorf("repository not found: %s", st.Message())
			case codes.InvalidArgument:
				return nil, fmt.Errorf("invalid request: %s", st.Message())
			default:
				return nil, fmt.Errorf("processor error: %s", st.Message())
			}
		}
		return nil, fmt.Errorf("failed to call processor: %w", err)
	}

	createdAt, err := time.Parse(time.RFC3339, resp.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to parse created_at: %w", err)
	}

	return &domain.Repository{
		Name:        resp.Name,
		Description: resp.Description,
		Stars:       int(resp.Stars),
		Forks:       int(resp.Forks),
		CreatedAt:   createdAt,
	}, nil
}

func (c *Client) Ping(ctx context.Context) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := c.client.Ping(ctx, &processorPb.PingRequest{})
	if err != nil {
		return "down", err
	}
	return resp.Status, nil
}

func (c *Client) GetSubscriptionsInfo(ctx context.Context) ([]*processorPb.RepoResponse, error) {
	req := &processorPb.Empty{}

	resp, err := c.client.GetSubscriptionsInfo(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp.Repositories, nil
}

package client

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/LuhTonkaYeat/GoHomeworks/hw2/api/proto"
)

type CollectorClient struct {
	client pb.CollectorServiceClient
	conn   *grpc.ClientConn
}

func NewCollectorClient(serverAddr string) (*CollectorClient, error) {
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(5*time.Second))
	if err != nil {
		return nil, err
	}

	client := pb.NewCollectorServiceClient(conn)

	return &CollectorClient{
		client: client,
		conn:   conn,
	}, nil
}

func (c *CollectorClient) GetRepository(ctx context.Context, owner, repo string) (*pb.RepoResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	req := &pb.RepoRequest{
		Owner: owner,
		Repo:  repo,
	}

	resp, err := c.client.GetRepository(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *CollectorClient) Close() error {
	return c.conn.Close()
}

func MapGrpcErrorToHTTP(err error) (int, string) {
	st, ok := status.FromError(err)
	if !ok {
		return 500, "Internal server error"
	}

	switch st.Code() {
	case codes.NotFound:
		return 404, st.Message()
	case codes.InvalidArgument:
		return 400, st.Message()
	case codes.DeadlineExceeded:
		return 504, "Request timeout"
	case codes.Unavailable:
		return 503, "Collector service unavailable"
	default:
		return 500, st.Message()
	}
}

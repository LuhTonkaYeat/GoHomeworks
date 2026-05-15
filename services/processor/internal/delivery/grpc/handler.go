package grpc

import (
	"context"

	pb "github.com/LuhTonkaYeat/GoHomeworks/services/processor/api/proto/processor"
	"github.com/LuhTonkaYeat/GoHomeworks/services/processor/internal/adapter/grpc"
	"github.com/LuhTonkaYeat/GoHomeworks/services/processor/internal/usecase"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Handler struct {
	pb.UnimplementedProcessorServiceServer
	repoUseCase     usecase.RepositoryUseCase
	collectorClient *grpc.Client
}

func NewHandler(repoUseCase usecase.RepositoryUseCase, collectorClient *grpc.Client) *Handler {
	return &Handler{
		repoUseCase:     repoUseCase,
		collectorClient: collectorClient,
	}
}

func (h *Handler) GetRepository(ctx context.Context, req *pb.RepoRequest) (*pb.RepoResponse, error) {
	repo, err := h.repoUseCase.GetRepository(ctx, req.Owner, req.Repo)
	if err != nil {
		if err.Error() == "owner and repo are required" {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}

		if len(err.Error()) > 10 && (err.Error()[:10] == "repository" || err.Error()[:10] == "collector ") {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.RepoResponse{
		Name:        repo.Name,
		Description: repo.Description,
		Stars:       int32(repo.Stars),
		Forks:       int32(repo.Forks),
		CreatedAt:   repo.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}, nil
}

func (h *Handler) Ping(ctx context.Context, req *pb.PingRequest) (*pb.PingResponse, error) {
	status, err := h.repoUseCase.Ping(ctx)
	if err != nil {
		return &pb.PingResponse{Status: "down"}, nil
	}
	return &pb.PingResponse{Status: status}, nil
}

func (h *Handler) GetSubscriptionsInfo(ctx context.Context, req *pb.Empty) (*pb.SubscriptionsInfoResponse, error) {
	resp, err := h.collectorClient.GetSubscriptionsInfo(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var repos []*pb.RepoResponse
	for _, repo := range resp.Repositories {
		repos = append(repos, &pb.RepoResponse{
			Name:        repo.Name,
			Description: repo.Description,
			Stars:       repo.Stars,
			Forks:       repo.Forks,
			CreatedAt:   repo.CreatedAt,
		})
	}

	return &pb.SubscriptionsInfoResponse{Repositories: repos}, nil
}

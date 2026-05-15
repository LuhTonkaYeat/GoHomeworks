package grpc

import (
	"context"
	"fmt"
	"time"

	pb "github.com/LuhTonkaYeat/GoHomeworks/services/collector/api/proto/collector"
	"github.com/LuhTonkaYeat/GoHomeworks/services/collector/internal/usecase"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Handler struct {
	pb.UnimplementedCollectorServiceServer
	repoUseCase usecase.RepositoryUseCase
}

func NewServer(repoUseCase usecase.RepositoryUseCase) *Handler {
	return &Handler{
		repoUseCase: repoUseCase,
	}
}

func (h *Handler) GetRepository(ctx context.Context, req *pb.CollectorRepoRequest) (*pb.CollectorRepoResponse, error) {
	if req.Owner == "" || req.Repo == "" {
		return nil, status.Error(codes.InvalidArgument, "owner and repo are required")
	}

	repo, err := h.repoUseCase.GetRepository(ctx, req.Owner, req.Repo)
	if err != nil {
		if err.Error() == fmt.Sprintf("repository %s/%s not found", req.Owner, req.Repo) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.CollectorRepoResponse{
		Name:        repo.Name,
		Description: repo.Description,
		Stars:       int32(repo.Stars),
		Forks:       int32(repo.Forks),
		CreatedAt:   repo.CreatedAt.Format(time.RFC3339),
	}, nil
}

func (h *Handler) GetSubscriptionsInfo(ctx context.Context, req *pb.Empty) (*pb.SubscriptionsInfoResponse, error) {
	repositories, err := h.repoUseCase.GetSubscriptionsInfo(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var repos []*pb.CollectorRepoResponse
	for _, repo := range repositories {
		repos = append(repos, &pb.CollectorRepoResponse{
			Name:        repo.Name,
			Description: repo.Description,
			Stars:       int32(repo.Stars),
			Forks:       int32(repo.Forks),
			CreatedAt:   repo.CreatedAt.Format(time.RFC3339),
		})
	}

	return &pb.SubscriptionsInfoResponse{Repositories: repos}, nil
}

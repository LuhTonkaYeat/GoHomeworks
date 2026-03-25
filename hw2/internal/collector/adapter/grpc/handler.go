package grpc

import (
	"context"
	"fmt"
	"time"

	pb "github.com/LuhTonkaYeat/GoHomeworks/hw2/api/proto"
	"github.com/LuhTonkaYeat/GoHomeworks/hw2/internal/collector/usecase"
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

func (h *Handler) GetRepository(ctx context.Context, req *pb.RepoRequest) (*pb.RepoResponse, error) {
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

	return &pb.RepoResponse{
		Name:        repo.Name,
		Description: repo.Description,
		Stars:       int32(repo.Stars),
		Forks:       int32(repo.Forks),
		CreatedAt:   repo.CreatedAt.Format(time.RFC3339),
	}, nil
}

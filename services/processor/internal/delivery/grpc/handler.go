package grpc

import (
	"context"

	pb "github.com/LuhTonkaYeat/GoHomeworks/services/processor/api/proto/processor"
	"github.com/LuhTonkaYeat/GoHomeworks/services/processor/internal/usecase"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Handler struct {
	pb.UnimplementedProcessorServiceServer
	repoUseCase usecase.RepositoryUseCase
}

func NewHandler(repoUseCase usecase.RepositoryUseCase) *Handler {
	return &Handler{
		repoUseCase: repoUseCase,
	}
}

func (h *Handler) GetRepository(ctx context.Context, req *pb.RepoRequest) (*pb.RepoResponse, error) {
	repo, err := h.repoUseCase.GetRepository(ctx, req.Owner, req.Repo)
	if err != nil {
		if err.Error() == "owner and repo are required" {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		
		if _, ok := err.(*usecase.NotFoundError); ok {
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
	return &pb.PingResponse{Status: "up"}, nil
}

func (h *Handler) GetSubscriptionsInfo(ctx context.Context, req *pb.Empty) (*pb.SubscriptionsInfoResponse, error) {
	repos, err := h.repoUseCase.GetSubscriptionsInfo(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	
	var responses []*pb.RepoResponse
	for _, repo := range repos {
		responses = append(responses, &pb.RepoResponse{
			Name:        repo.Name,
			Description: repo.Description,
			Stars:       int32(repo.Stars),
			Forks:       int32(repo.Forks),
			CreatedAt:   repo.CreatedAt.Format("2006-01-02T15:04:05Z"),
		})
	}
	
	return &pb.SubscriptionsInfoResponse{Repositories: responses}, nil
}
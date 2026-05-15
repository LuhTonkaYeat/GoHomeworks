package grpc

import (
	"context"

	pb "github.com/LuhTonkaYeat/GoHomeworks/services/subscribe/api/proto"
	"github.com/LuhTonkaYeat/GoHomeworks/services/subscribe/internal/usecase"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Handler struct {
	pb.UnimplementedSubscribeServiceServer
	useCase usecase.SubscriptionUseCase
}

func NewHandler(useCase usecase.SubscriptionUseCase) *Handler {
	return &Handler{
		useCase: useCase,
	}
}

func (h *Handler) CreateSubscription(ctx context.Context, req *pb.CreateSubscriptionRequest) (*pb.SubscriptionResponse, error) {
	if req.Owner == "" || req.Repo == "" {
		return nil, status.Error(codes.InvalidArgument, "owner and repo are required")
	}

	userID := req.UserId
	if userID == "" {
		userID = "default"
	}

	err := h.useCase.CreateSubscription(ctx, userID, req.Owner, req.Repo)
	if err != nil {
		return &pb.SubscriptionResponse{Success: false, Message: err.Error()}, nil
	}

	return &pb.SubscriptionResponse{Success: true, Message: "subscription created"}, nil
}

func (h *Handler) DeleteSubscription(ctx context.Context, req *pb.DeleteSubscriptionRequest) (*pb.Empty, error) {
	userID := req.UserId
	if userID == "" {
		userID = "default"
	}

	err := h.useCase.DeleteSubscription(ctx, userID, req.Owner, req.Repo)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.Empty{}, nil
}

func (h *Handler) GetSubscriptions(ctx context.Context, req *pb.Empty) (*pb.SubscriptionsList, error) {
	subscriptions, err := h.useCase.GetSubscriptions(ctx, "default")
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var repos []*pb.Repository
	for _, sub := range subscriptions {
		repos = append(repos, &pb.Repository{
			Owner: sub.Owner,
			Name:  sub.Repo,
		})
	}

	return &pb.SubscriptionsList{Subscriptions: repos}, nil
}

func (h *Handler) CheckRepositoryExists(ctx context.Context, req *pb.CheckRepoRequest) (*pb.CheckRepoResponse, error) {
	exists, err := h.useCase.CheckRepoExists(ctx, req.Owner, req.Repo)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.CheckRepoResponse{Exists: exists}, nil
}

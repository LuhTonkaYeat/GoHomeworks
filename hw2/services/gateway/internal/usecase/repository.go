package usecase

import (
	"context"
	"fmt"

	"github.com/LuhTonkaYeat/GoHomeworks/hw2/services/gateway/internal/domain"
)

type ServiceStatus struct {
	Name   string
	Status string
}

type RepositoryUseCase interface {
	GetRepository(ctx context.Context, owner, repo string) (*domain.Repository, error)
	Ping(ctx context.Context) (string, []ServiceStatus, error)
}

type processorClient interface {
	GetRepository(ctx context.Context, owner, repo string) (*domain.Repository, error)
	Ping(ctx context.Context) (string, error)
}

type repositoryUseCase struct {
	processorClient processorClient
}

func NewRepositoryUseCase(processorClient processorClient) RepositoryUseCase {
	return &repositoryUseCase{
		processorClient: processorClient,
	}
}

func (uc *repositoryUseCase) GetRepository(ctx context.Context, owner, repo string) (*domain.Repository, error) {
	if owner == "" || repo == "" {
		return nil, fmt.Errorf("owner and repo are required")
	}
	return uc.processorClient.GetRepository(ctx, owner, repo)
}

func (uc *repositoryUseCase) Ping(ctx context.Context) (string, []ServiceStatus, error) {
	status, err := uc.processorClient.Ping(ctx)
	if err != nil {
		return "degraded", []ServiceStatus{
			{Name: "processor", Status: "down"},
			{Name: "subscriber", Status: "unknown"},
		}, nil
	}

	services := []ServiceStatus{
		{Name: "processor", Status: status},
		{Name: "subscriber", Status: "up"},
	}

	overallStatus := "ok"
	for _, s := range services {
		if s.Status != "up" {
			overallStatus = "degraded"
			break
		}
	}

	return overallStatus, services, nil
}

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	pb "github.com/LuhTonkaYeat/GoHomeworks/hw2/api/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct {
	pb.UnimplementedCollectorServiceServer
}

type GitHubRepo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Stargazers  int    `json:"stargazers_count"`
	Forks       int    `json:"forks_count"`
	CreatedAt   string `json:"created_at"`
}

func (s *server) GetRepository(ctx context.Context, req *pb.RepoRequest) (*pb.RepoResponse, error) {
	log.Printf("Received request for repo: %s/%s", req.Owner, req.Repo)

	url := fmt.Sprintf("https://api.github.com/repos/%s/%s", req.Owner, req.Repo)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		log.Printf("Network error: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to reach GitHub API: %v", err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		var gitHubRepo GitHubRepo
		if err := json.NewDecoder(resp.Body).Decode(&gitHubRepo); err != nil {
			log.Printf("JSON parse error: %v", err)
			return nil, status.Errorf(codes.Internal, "failed to parse GitHub response: %v", err)
		}

		return &pb.RepoResponse{
			Name:        gitHubRepo.Name,
			Description: gitHubRepo.Description,
			Stars:       int32(gitHubRepo.Stargazers),
			Forks:       int32(gitHubRepo.Forks),
			CreatedAt:   gitHubRepo.CreatedAt,
		}, nil

	case http.StatusNotFound:
		log.Printf("Repository not found: %s/%s", req.Owner, req.Repo)
		return nil, status.Errorf(codes.NotFound, "repository %s/%s not found", req.Owner, req.Repo)

	default:
		log.Printf("GitHub API returned status %d", resp.StatusCode)
		return nil, status.Errorf(codes.Internal, "GitHub API returned status %d", resp.StatusCode)
	}
}

func main() {
	port := ":50051"
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", port, err)
	}
	grpcServer := grpc.NewServer()

	pb.RegisterCollectorServiceServer(grpcServer, &server{})

	log.Printf("Collector server is running on port %s", port)
	log.Println("Waiting for gRPC requests...")

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

package main

import (
	"context"
	"log"
	"net"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"

	pb "github.com/LuhTonkaYeat/GoHomeworks/services/subscribe/api/proto"
	"github.com/LuhTonkaYeat/GoHomeworks/services/subscribe/internal/adapter/github"
	grpcdelivery "github.com/LuhTonkaYeat/GoHomeworks/services/subscribe/internal/delivery/grpc"
	"github.com/LuhTonkaYeat/GoHomeworks/services/subscribe/internal/repository"
	"github.com/LuhTonkaYeat/GoHomeworks/services/subscribe/internal/usecase"
)

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/subscribe?sslmode=disable"
	}

	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer pool.Close()

	queries := repository.New(pool)

	githubClient := github.NewClient()

	subscriptionUC := usecase.NewSubscriptionUseCase(queries, githubClient)

	grpcHandler := grpcdelivery.NewHandler(subscriptionUC)

	port := os.Getenv("PORT")
	if port == "" {
		port = "50053"
	}

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterSubscribeServiceServer(grpcServer, grpcHandler)

	log.Printf("Subscribe service listening on port %s", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

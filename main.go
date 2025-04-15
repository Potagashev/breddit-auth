package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/jackc/pgx/v5"
	"google.golang.org/grpc"

	"github.com/Potagashev/breddit_auth/internal/auth"
	"github.com/Potagashev/breddit_auth/internal/config"
	"github.com/Potagashev/breddit_auth/internal/router"
	"github.com/Potagashev/breddit_auth/internal/users"
	authPb "github.com/Potagashev/breddit_auth/internal/auth/proto"
)

// @title Swagger Example API
// @version 1.0
// @description This is a authentication service.
// @host localhost:8080
// @BasePath /api/v1
func main() {
	cfg, _ := config.LoadConfig()

	conn, err := pgx.Connect(context.Background(), cfg.DbUrl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	user_repository := users.NewUserRepository(conn)
	user_service := users.NewUserService(user_repository)
	auth_service := auth.NewAuthService(user_service, cfg)
	auth_handler := auth.NewAuthHandler(auth_service)

	r := router.NewRouter(auth_handler)

	listener, err := net.Listen("tcp", ":50051")
    if err != nil {
        log.Fatalf("Failed to listen: %v", err)
    }

    grpcServer := grpc.NewServer()

	authPb.RegisterAuthServiceServer(grpcServer, auth_service)

    log.Println("gRPC server is running on port 50051")
    if err := grpcServer.Serve(listener); err != nil {
        log.Fatalf("Failed to serve: %v", err)
    }

	r.Run(fmt.Sprintf(":%s", cfg.AppPort))
}

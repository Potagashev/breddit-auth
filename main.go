package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"

	"github.com/Potagashev/breddit_auth/internal/auth"
	"github.com/Potagashev/breddit_auth/internal/config"
	"github.com/Potagashev/breddit_auth/internal/users"
	"github.com/Potagashev/breddit_auth/internal/router"
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

	r.Run(fmt.Sprintf(":%s", cfg.AppPort))
}

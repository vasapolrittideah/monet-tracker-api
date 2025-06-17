package main

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	userpbv1 "github.com/vasapolrittideah/money-tracker-api/protogen/user/v1"
	grpchandler "github.com/vasapolrittideah/money-tracker-api/services/user/internal/delivery/grpc"
	httphandler "github.com/vasapolrittideah/money-tracker-api/services/user/internal/delivery/http"
	"github.com/vasapolrittideah/money-tracker-api/services/user/internal/repository"
	"github.com/vasapolrittideah/money-tracker-api/services/user/internal/usecase"
	"github.com/vasapolrittideah/money-tracker-api/shared/bootstrap"
	"github.com/vasapolrittideah/money-tracker-api/shared/config"
	"github.com/vasapolrittideah/money-tracker-api/shared/database"
	"github.com/vasapolrittideah/money-tracker-api/shared/validator"
	"google.golang.org/grpc"
)

// @title Money Tracker API
// @version 1.0
// @description	This is a user service for Money Tracker API
// @contact.name Vasapol Rittideah
// @contact.email	vasapol.rittideah@outlook.com
// @license.name MIT
// @license.url https://github.com/vasapolrittideah/money-tracker-api/blob/main/LICENSE
// @host moneytracker.local
// @BasePath /api/v1
func main() {
	cfg := config.Load()
	db := database.Connect(&cfg.Database)
	validator.Init()

	httpAddr := fmt.Sprintf(":%s", cfg.Server.UserServiceHTTPPort)
	grpcAddr := fmt.Sprintf(":%s", cfg.Server.UserServiceGRPCPort)

	app := bootstrap.NewApp("user-service", cfg, db)

	app.RegisterGRPCServiceWithConsul(
		"user-service-1",
		"user-service",
		cfg.Server.UserServiceHost,
		cfg.Server.UserServiceGRPCPort,
		10*time.Second,
		1*time.Minute,
	)

	userRepository := repository.NewUserRepository(db)
	userUsecase := usecase.NewUserUsecase(userRepository, cfg)

	app.AddHTTPServer(httpAddr, func(router fiber.Router) {
		v1 := router.Group("/api/v1")
		userHandler := httphandler.NewUserHTTPHandler(userUsecase, v1, cfg)
		userHandler.RegisterRoutes()
	})

	app.AddGRPCServer(grpcAddr, func(server *grpc.Server) {
		userHandler := grpchandler.NewUserGRPCHandler(userUsecase, cfg)
		userpbv1.RegisterUserServiceServer(server, userHandler)
	})

	app.Run()
}

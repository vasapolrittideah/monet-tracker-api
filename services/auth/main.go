package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	userpbv1 "github.com/vasapolrittideah/money-tracker-api/protogen/user/v1"
	httphandler "github.com/vasapolrittideah/money-tracker-api/services/auth/internal/delivery/http"
	"github.com/vasapolrittideah/money-tracker-api/services/auth/internal/repository"
	"github.com/vasapolrittideah/money-tracker-api/services/auth/internal/usecase"
	"github.com/vasapolrittideah/money-tracker-api/shared/bootstrap"
	"github.com/vasapolrittideah/money-tracker-api/shared/config"
	"github.com/vasapolrittideah/money-tracker-api/shared/database"
	"github.com/vasapolrittideah/money-tracker-api/shared/validator"
)

func main() {
	cfg := config.Load()
	db := database.Connect(&cfg.Database)
	validator.Init()

	httpAddr := fmt.Sprintf(":%s", cfg.Server.AuthServiceHTTPPort)

	app := bootstrap.NewApp("auth-service", cfg, db)

	app.ConnectGRPCClientsFromConsul([]string{"user-service"})

	userClient := userpbv1.NewUserServiceClient(app.GetGRPCClient("user-service"))
	externalAuthRepository := repository.NewExternalAuthRepository(db)
	sessionRepository := repository.NewSessionRepository(db)
	authUsecase := usecase.NewAuthUsecase(userClient, sessionRepository, cfg)
	oauthGoogleUsecase := usecase.NewOAuthGoogleUsecase(userClient, authUsecase, externalAuthRepository, cfg)

	app.AddHTTPServer(httpAddr, func(router fiber.Router) {
		v1 := router.Group("/api/v1")
		authHandler := httphandler.NewAuthHTTPHandler(authUsecase, v1, cfg)
		authHandler.RegisterRoutes()
		oauthGoogleHandler := httphandler.NewOAuthGoogleHTTPHandler(oauthGoogleUsecase, v1, cfg)
		oauthGoogleHandler.RegisterRoutes()
	})

	app.Run()
}

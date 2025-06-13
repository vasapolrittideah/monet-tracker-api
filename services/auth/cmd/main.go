package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"

	"github.com/charmbracelet/log"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	userpbv1 "github.com/vasapolrittideah/money-tracker-api/protogen/user/v1"
	"github.com/vasapolrittideah/money-tracker-api/services/auth/docs"
	"github.com/vasapolrittideah/money-tracker-api/services/auth/internal/controller"
	"github.com/vasapolrittideah/money-tracker-api/services/auth/internal/usecase"
	"github.com/vasapolrittideah/money-tracker-api/shared/bootstrap"
	"github.com/vasapolrittideah/money-tracker-api/shared/consul"
	"github.com/vasapolrittideah/money-tracker-api/shared/middleware"
	"google.golang.org/grpc"
)

// @title Money Tracker API
// @version 1.0
// @description	This is an auth service for Money Tracker API
// @contact.name Vasapol Rittideah
// @contact.email	vasapol.rittideah@outlook.com
// @license.name MIT
// @license.url https://github.com/vasapolrittideah/money-tracker-api/blob/main/LICENSE
// @host moneytracker.local
// @BasePath /api/v1
func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	app := bootstrap.NewApp()
	defer app.Close()

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		if err := startHTTPServer(ctx, &wg, &app); err != nil {
			log.Errorf("failed to start http server: %v", err)
			cancel()
		}
	}()

	wg.Wait()
	log.Info("👋 user service stopped gracefully")
}

func startHTTPServer(ctx context.Context, wg *sync.WaitGroup, app *bootstrap.Application) error {
	defer wg.Done()

	a := fiber.New()
	middleware.RegisterHTTPMiddleware(a)

	a.Get("/swagger/*", swagger.New(swagger.Config{URL: "/doc.json"}))
	a.Get("/doc.json", func(ctx *fiber.Ctx) error {
		ctx.Set("Content-Type", "application/json")
		return ctx.Status(200).SendString(docs.SwaggerInfo.ReadDoc())
	})

	conns, err := createAuthGRPCServiceConnections(app)
	if err != nil {
		return fmt.Errorf("failed to create auth grpc clients: %v", err)
	}

	userClient := userpbv1.NewUserServiceClient(conns["user-service"])

	router := a.Group("/api/v1")

	authUsecase := usecase.NewAuthUsecase(userClient, app.Config)
	authController := controller.NewAuthHTTPController(authUsecase, router, app.Config)
	authController.RegisterRoutes()

	addr := fmt.Sprintf(":%v", app.Config.Server.AuthServiceHTTPPort)

	go func() {
		<-ctx.Done()
		log.Info("🧹 shutting down http server...")
		if err := a.Shutdown(); err != nil {
			log.Errorf("failed to shutdown http server: %v", err)
		}
	}()

	log.Infof("🚀 http server started on %s", addr)

	if err := a.Listen(addr); err != nil {
		return fmt.Errorf("failed to listen on %s: %v", addr, err)
	}

	return nil
}

func createAuthGRPCServiceConnections(app *bootstrap.Application) (map[string]*grpc.ClientConn, error) {
	serviceNames := []string{"user-service"}

	address := fmt.Sprintf("%v:%v", app.Config.Server.ConsulHost, app.Config.Server.ConsulPort)
	consulClient, err := consul.NewConsulClient(address)
	if err != nil {
		return nil, fmt.Errorf("failed to create consul client: %v", err)
	}

	log.Infof("auth service connecting to consul on %s", address)

	conns, err := consulClient.CreateGRPCServiceConnections(serviceNames)
	if err != nil {
		return nil, fmt.Errorf("failed to create grpc clients: %v", err)
	}

	return conns, nil
}

package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"

	"github.com/charmbracelet/log"
	"github.com/gofiber/fiber/v2"
	userpbv1 "github.com/vasapolrittideah/money-tracker-api/protogen/user/v1"
	"github.com/vasapolrittideah/money-tracker-api/services/auth/controller"
	"github.com/vasapolrittideah/money-tracker-api/services/auth/usecase"
	"github.com/vasapolrittideah/money-tracker-api/shared/bootstrap"
	"github.com/vasapolrittideah/money-tracker-api/shared/consul"
	"github.com/vasapolrittideah/money-tracker-api/shared/middleware"
	"google.golang.org/grpc"
)

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

	conns, err := createAuthGRPCClients(app)
	if err != nil {
		return fmt.Errorf("failed to create auth grpc clients: %v", err)
	}

	userClient := userpbv1.NewUserServiceClient(conns["user-service"])

	router := a.Group("/api/v1")

	authUsecase := usecase.NewAuthUsecase(userClient, app.Config)
	authController := controller.NewAuthHTTPController(authUsecase, router, app.Config)
	authController.RegisterRoutes()

	go func() {
		<-ctx.Done()
		log.Info("🧹 shutting down http server...")
		if err := a.Shutdown(); err != nil {
			log.Errorf("failed to shutdown http server: %v", err)
		}
	}()

	addr := fmt.Sprintf(":%v", app.Config.Server.AuthServiceHTTPPort)
	if err := a.Listen(addr); err != nil {
		return fmt.Errorf("failed to listen on %s: %v", addr, err)
	}

	return nil
}

func createAuthGRPCClients(app *bootstrap.Application) (map[string]*grpc.ClientConn, error) {
	serviceNames := []string{"user-service"}

	address := fmt.Sprintf("%v:%v", app.Config.Server.ConsulHost, app.Config.Server.ConsulPort)
	consulClient, err := consul.NewConsulClient(address)
	if err != nil {
		return nil, fmt.Errorf("failed to create consul client: %v", err)
	}

	conns, err := consulClient.CreateGRPCClients(serviceNames)
	if err != nil {
		return nil, fmt.Errorf("failed to create grpc clients: %v", err)
	}

	return conns, nil
}

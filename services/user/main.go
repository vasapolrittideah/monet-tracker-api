package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/charmbracelet/log"
	"github.com/gofiber/fiber/v2"
	userpbv1 "github.com/vasapolrittideah/money-tracker-api/protogen/user/v1"
	"github.com/vasapolrittideah/money-tracker-api/services/user/controller"
	"github.com/vasapolrittideah/money-tracker-api/services/user/repository"
	"github.com/vasapolrittideah/money-tracker-api/services/user/usecase"
	"github.com/vasapolrittideah/money-tracker-api/shared/bootstrap"
	"github.com/vasapolrittideah/money-tracker-api/shared/consul"
	"github.com/vasapolrittideah/money-tracker-api/shared/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	app := bootstrap.NewApp()
	defer app.Close()

	if err := registerUserGRPCService(ctx, &app); err != nil {
		log.Errorf("failed to register user service: %v", err)
		return
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		if err := startHTTPServer(ctx, &wg, &app); err != nil {
			log.Errorf("failed to start http server: %v", err)
			cancel()
		}
	}()

	go func() {
		if err := startGRPCServer(ctx, &wg, &app); err != nil {
			log.Errorf("failed to start grpc server: %v", err)
			cancel()
		}
	}()

	wg.Wait()
	log.Info("👋 user service stopped gracefully")
}

func startGRPCServer(ctx context.Context, wg *sync.WaitGroup, app *bootstrap.Application) error {
	defer wg.Done()

	addr := fmt.Sprintf(":%v", app.Config.Server.UserServiceGRPCPort)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %v", addr, err)
	}

	userRepository := repository.NewUserRepository(app.DB)
	userUsecase := usecase.NewUserUsecase(userRepository, app.Config)
	userController := controller.NewUserGRPCController(userUsecase, app.Config)

	grpcServer := grpc.NewServer()
	userpbv1.RegisterUserServiceServer(grpcServer, userController)

	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)

	go func() {
		<-ctx.Done()
		log.Info("🧹 shutting down grpc server...")
		grpcServer.GracefulStop()
	}()

	log.Infof("🚀 grpc server started on %s", addr)

	if err := grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve grpc server: %v", err)
	}

	return nil
}

func startHTTPServer(ctx context.Context, wg *sync.WaitGroup, app *bootstrap.Application) error {
	defer wg.Done()

	a := fiber.New()
	middleware.RegisterHTTPMiddleware(a)

	router := a.Group("/api/v1")

	userRepository := repository.NewUserRepository(app.DB)
	userUsecase := usecase.NewUserUsecase(userRepository, app.Config)
	userController := controller.NewUserHTTPController(router, userUsecase, app.Config)
	userController.RegisterRoutes()

	addr := fmt.Sprintf(":%v", app.Config.Server.UserServiceHTTPPort)

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

func registerUserGRPCService(ctx context.Context, app *bootstrap.Application) error {
	address := fmt.Sprintf("%v:%v", app.Config.Server.ConsulHost, app.Config.Server.ConsulPort)
	consulClient, err := consul.NewConsulClient(address)
	if err != nil {
		return err
	}

	serviceID := "user-service-1"
	serviceName := "user-service"
	serviceAddress := app.Config.Server.UserServiceHost
	servicePort, _ := strconv.Atoi(app.Config.Server.UserServiceGRPCPort)

	err = consulClient.RegisterGRPCService(
		serviceID,
		serviceName,
		serviceAddress,
		servicePort,
		10*time.Second,
		1*time.Minute,
	)
	if err != nil {
		return err
	}

	log.Info("🎉 user service registered successfully")

	go func() {
		<-ctx.Done()
		if err := consulClient.DeregisterService(serviceID); err != nil {
			log.Error(err)
		} else {
			log.Info("user service deregistered successfully")
		}
	}()

	return nil
}

package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"

	"github.com/charmbracelet/log"
	"github.com/gofiber/fiber/v2"
	userv1 "github.com/vasapolrittideah/money-tracker-api/protogen/user/v1"
	"github.com/vasapolrittideah/money-tracker-api/services/user/controller"
	"github.com/vasapolrittideah/money-tracker-api/services/user/repository"
	"github.com/vasapolrittideah/money-tracker-api/services/user/usecase"
	"github.com/vasapolrittideah/money-tracker-api/shared/bootstrap"
	"github.com/vasapolrittideah/money-tracker-api/shared/middleware"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	app := bootstrap.NewApp()
	defer app.Close()

	group, ctx := errgroup.WithContext(ctx)

	group.Go(grpcServerRunner(ctx, &app))
	group.Go(httpServerRunner(ctx, &app))

	if err := group.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		log.Errorf("failed to run user service: %v", err)
	} else {
		log.Info("👋 user service stopped gracefully")
	}
}

func grpcServerRunner(ctx context.Context, app *bootstrap.Application) func() error {
	return func() error {
		grpcServer := grpc.NewServer()
		reflection.Register(grpcServer)

		userRepository := repository.NewUserRepository(app.DB)
		userUsecase := usecase.NewUserUsecase(userRepository, app.Config)
		userController := controller.NewUserGRPCController(userUsecase, app.Config)
		userv1.RegisterUserServiceServer(grpcServer, userController)

		group, ctx := errgroup.WithContext(ctx)

		group.Go(func() error {
			addr := fmt.Sprintf(":%v", app.Config.Server.UserServiceGRPCPort)
			lis, err := net.Listen("tcp", addr)
			if err != nil {
				return fmt.Errorf("failed to listen on %s: %w", addr, err)
			}

			log.Infof("🚀 grpc server started on %s", addr)

			if err := grpcServer.Serve(lis); err != nil {
				return fmt.Errorf("failed to serve grpc server: %v", err)
			}

			return nil
		})

		group.Go(func() error {
			<-ctx.Done()

			log.Info("🧹 shutting down user service gracefully...")
			grpcServer.GracefulStop()

			return ctx.Err()
		})

		return group.Wait()
	}
}

func httpServerRunner(ctx context.Context, app *bootstrap.Application) func() error {
	return func() error {
		a := fiber.New()
		middleware.RegisterHTTPMiddleware(a)

		router := a.Group("/api/v1")

		userRepository := repository.NewUserRepository(app.DB)
		userUsecase := usecase.NewUserUsecase(userRepository, app.Config)
		userController := controller.NewUserHTTPController(router, userUsecase, app.Config)
		userController.RegisterRoutes()

		group, ctx := errgroup.WithContext(ctx)

		group.Go(func() error {
			addr := fmt.Sprintf(":%v", app.Config.Server.UserServiceHTTPPort)

			log.Infof("🚀 http server started on %s", addr)

			if err := a.Listen(addr); err != nil {
				return fmt.Errorf("failed to listen on %s: %w", addr, err)
			}

			return nil
		})

		group.Go(func() error {
			<-ctx.Done()

			log.Info("🧹 shutting down user service gracefully...")
			if err := a.Shutdown(); err != nil {
				return fmt.Errorf("failed to shutdown http server: %w", err)
			}

			return ctx.Err()
		})

		return group.Wait()
	}
}

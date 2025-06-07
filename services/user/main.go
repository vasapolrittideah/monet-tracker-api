package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"

	"buf.build/go/protovalidate"
	"github.com/charmbracelet/log"
	protovalidate_middleware "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/protovalidate"
	userv1 "github.com/vasapolrittideah/money-tracker-api/protogen/user/v1"
	"github.com/vasapolrittideah/money-tracker-api/services/user/controller"
	"github.com/vasapolrittideah/money-tracker-api/services/user/repository"
	"github.com/vasapolrittideah/money-tracker-api/services/user/usecase"
	"github.com/vasapolrittideah/money-tracker-api/shared/bootstrap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	if err := run(ctx); err != nil && !errors.Is(err, context.Canceled) {
		log.Errorf("failed to run user service: %v", err)
		return
	}

	log.Info("👋 user service stopped gracefully")
}

func run(ctx context.Context) error {
	app := bootstrap.NewApp()
	defer app.Close()

	validator, err := protovalidate.New()
	if err != nil {
		return fmt.Errorf("failed to create protovalidate validator: %w", err)
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(protovalidate_middleware.UnaryServerInterceptor(validator)),
	)

	reflection.Register(grpcServer)

	userRepository := repository.NewUserRepository(app.DB)
	userUsecase := usecase.NewUserUsecase(userRepository, app.Config)
	userController := controller.NewUserController(userUsecase, app.Config)
	userv1.RegisterUserServiceServer(grpcServer, userController)

	errorGroup, ctx := errgroup.WithContext(ctx)

	errorGroup.Go(func() error {
		addr := fmt.Sprintf(":%v", app.Config.Server.UserServicePort)
		lis, err := net.Listen("tcp", addr)
		if err != nil {
			return fmt.Errorf("failed to listen on %s: %w", addr, err)
		}

		log.Infof("🚀 user service listening on %s", addr)

		if err := grpcServer.Serve(lis); err != nil {
			return fmt.Errorf("failed to serve grpc server: %v", err)
		}

		return nil
	})

	errorGroup.Go(func() error {
		<-ctx.Done()

		log.Info("🧹 shutting down user service gracefully...")
		grpcServer.GracefulStop()

		return ctx.Err()
	})

	return errorGroup.Wait()
}

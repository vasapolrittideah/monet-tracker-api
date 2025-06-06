package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/charmbracelet/log"
	"github.com/vasapolrittideah/money-tracker-api/services/user/controller"
	"github.com/vasapolrittideah/money-tracker-api/services/user/repository"
	"github.com/vasapolrittideah/money-tracker-api/services/user/usecase"
	"github.com/vasapolrittideah/money-tracker-api/shared/bootstrap"
	"google.golang.org/grpc"
)

func main() {
	app := bootstrap.NewApp()
	defer app.Close()

	addr := fmt.Sprintf(":%v", app.Config.Server.UserServicePort)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Errorf("failed to listen on %s: %v", addr, err)
		return
	}
	defer func() {
		if err := lis.Close(); err != nil {
			log.Errorf("failed to close listener: %v", err)
		}
	}()

	grpcServer := grpc.NewServer()

	userRepository := repository.NewUserRepository(app.DB)
	userService := usecase.NewUserUsecase(userRepository, app.Config)
	controller.NewUserController(grpcServer, userService, app.Config)

	go func() {
		log.Infof("🚀 user service started on port %v", app.Config.Server.UserServicePort)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatal("failed to serve grpc server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Info("👋 shutdown signal received, stopping server...")

	grpcServer.GracefulStop()

	log.Info("👋 server stopped, see you later")
}

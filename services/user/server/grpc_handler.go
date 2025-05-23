package server

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/vasapolrittideah/money-tracker-api/services/user/handler"
	"github.com/vasapolrittideah/money-tracker-api/services/user/repository"
	"github.com/vasapolrittideah/money-tracker-api/services/user/service"
	"github.com/vasapolrittideah/money-tracker-api/shared/config"
	"github.com/vasapolrittideah/money-tracker-api/shared/logger"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

type userGrpcServer struct {
	cfg *config.Config
	db  *gorm.DB
}

func NewUserGrpcServer(cfg *config.Config, db *gorm.DB) *userGrpcServer {
	return &userGrpcServer{cfg: cfg, db: db}
}

func (s *userGrpcServer) Run() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", s.cfg.Server.UserServiceGrpcPort))
	if err != nil {
		logger.Fatal("USER", "failed to listen: %v", err)
	}

	userGrpcServer := grpc.NewServer()

	userService := service.NewUserService(repository.NewUserRepository(s.db), s.cfg)
	handler.NewUserGrpcHandler(userGrpcServer, userService, s.cfg)

	go func() {
		if err := userGrpcServer.Serve(lis); err != nil {
			logger.Fatal("USER", "failed to serve grpc server: %v", err)
		}
	}()

	logger.Info("USER", "🚀 grpc server started on port %v", s.cfg.Server.UserServiceGrpcPort)

	quit := make(chan os.Signal, 1)
	signal.Notify(
		quit,
		os.Interrupt,
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	<-quit
}

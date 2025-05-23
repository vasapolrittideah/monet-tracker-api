package server

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	fiberlogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	userpb "github.com/vasapolrittideah/money-tracker-api/generated/protobuf/user"
	"github.com/vasapolrittideah/money-tracker-api/services/auth/handler"
	"github.com/vasapolrittideah/money-tracker-api/services/auth/service"
	"github.com/vasapolrittideah/money-tracker-api/shared/config"
	"github.com/vasapolrittideah/money-tracker-api/shared/logger"
	"github.com/vasapolrittideah/money-tracker-api/shared/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type httpServer struct {
	cfg *config.Config
}

func NewAuthHttpServer(cfg *config.Config) *httpServer {
	return &httpServer{cfg: cfg}
}

func (s *httpServer) Run() {
	app := fiber.New()

	loggerConfig := fiberlogger.Config{
		TimeFormat: time.RFC1123Z,
		TimeZone:   "Asia/Bangkok",
	}

	corsConfig := cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
		AllowMethods: strings.Join([]string{
			fiber.MethodGet,
			fiber.MethodPost,
			fiber.MethodPut,
			fiber.MethodDelete,
		}, ","),
	}

	app.Use(
		recover.New(),
		fiberlogger.New(loggerConfig),
		cors.New(corsConfig),
	)

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).SendString("OK")
	})

	conn := newUserClient(s.cfg)
	defer func() {
		cerr := conn.Close()
		if cerr != nil {
			logger.Error("AUTH", "failed to close connection: %v", cerr)
		}
	}()

	userClient := userpb.NewUserServiceClient(conn)
	router := app.Group("/api")

	authService := service.NewAuthService(userClient, s.cfg)
	coreMiddleware := middleware.NewCoreMiddleware(s.cfg)
	authHandler := handler.NewAuthHttpHandler(authService, coreMiddleware, router, s.cfg)
	authHandler.RegisterRouter()

	go func() {
		if err := app.Listen(":" + s.cfg.Server.AuthServiceHttpPort); err != nil {
			logger.Fatal("AUTH", "failed to listen and serve application: %v", err)
		}
	}()

	logger.Info("AUTH", "🚀 http server started on port %v", s.cfg.Server.AuthServiceHttpPort)

	quit := make(chan os.Signal, 1)
	signal.Notify(
		quit,
		os.Interrupt,
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	<-quit
}

func newUserClient(cfg *config.Config) *grpc.ClientConn {
	conn, err := grpc.NewClient(
		fmt.Sprintf("%v:%v", cfg.Server.UserServiceGrpcConnectionHost, cfg.Server.UserServiceGrpcPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		logger.Fatal("AUTH", "failed to connect to user grpc server: %v", err)
	}

	logger.Info("AUTH", "🎉 connected to user grpc server")

	return conn
}

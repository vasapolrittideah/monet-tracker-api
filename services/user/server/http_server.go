package server

import (
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	fiberlogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/vasapolrittideah/money-tracker-api/services/user/handler"
	"github.com/vasapolrittideah/money-tracker-api/services/user/repository"
	"github.com/vasapolrittideah/money-tracker-api/services/user/service"
	"github.com/vasapolrittideah/money-tracker-api/shared/config"
	"github.com/vasapolrittideah/money-tracker-api/shared/logger"
	"gorm.io/gorm"
)

type userHttpServer struct {
	cfg *config.Config
	db  *gorm.DB
}

func NewUserHttpServer(cfg *config.Config, db *gorm.DB) *userHttpServer {
	return &userHttpServer{cfg, db}
}

func (s *userHttpServer) Run() {
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
		return c.Status(fiber.StatusOK).SendString("User service is healthy")
	})

	router := app.Group("/api")
	userService := service.NewUserService(repository.NewUserRepository(s.db), s.cfg)
	userHttpHandler := handler.NewUserHttpHandler(router, userService, s.cfg)
	userHttpHandler.RegisterRouter()

	go func() {
		if err := app.Listen(":" + s.cfg.Server.UserHttpPort); err != nil {
			logger.Fatal("USER", "failed to serve http server: %v", err)
		}
	}()

	logger.Info("USER", "🚀 http server started on port %v", s.cfg.Server.UserHttpPort)

	quit := make(chan os.Signal, 1)
	signal.Notify(
		quit,
		os.Interrupt,
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	<-quit
}

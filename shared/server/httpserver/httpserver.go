package httpserver

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/log"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

type RouteRegistrar func(router fiber.Router)

type HTTPServer struct {
	Addr           string
	RouteRegistrar RouteRegistrar
	app            *fiber.App
}

func (s *HTTPServer) Start() error {
	s.app = fiber.New()

	loggerConfig := logger.Config{
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

	s.app.Use(
		recover.New(),
		logger.New(loggerConfig),
		cors.New(corsConfig),
	)

	if s.RouteRegistrar != nil {
		s.RouteRegistrar(s.app)
	}

	log.Infof("🚀 HTTP server started on %s", s.Addr)

	if err := s.app.Listen(s.Addr); err != nil {
		return fmt.Errorf("failed to listen on %s: %v", s.Addr, err)
	}

	return nil
}

func (s *HTTPServer) Stop() error {
	return s.app.Shutdown()
}

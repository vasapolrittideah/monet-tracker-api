package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/charmbracelet/log"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"github.com/vasapolrittideah/money-tracker-api/swagger/docs"
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

	a := fiber.New()

	a.Get("/swagger/*", swagger.New(swagger.Config{URL: "/doc.json"}))
	a.Get("/doc.json", func(ctx *fiber.Ctx) error {
		ctx.Set("Content-Type", "application/json")
		return ctx.Status(200).SendString(docs.SwaggerInfo.ReadDoc())
	})

	go func() {
		<-ctx.Done()
		log.Info("🧹 shutting down http server...")
		if err := a.Shutdown(); err != nil {
			log.Errorf("failed to shutdown http server: %v", err)
		}
	}()

	addr := ":10000"

	log.Infof("🚀 http server started on %s", addr)

	if err := a.Listen(addr); err != nil {
		log.Errorf("failed to listen on %s: %v", addr, err)
	}
}

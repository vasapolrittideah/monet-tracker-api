package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/vasapolrittideah/money-tracker-api/services/auth/service"
	"github.com/vasapolrittideah/money-tracker-api/shared/config"
	"github.com/vasapolrittideah/money-tracker-api/shared/model/response"
)

type OAuthGoogleHandler interface {
	RegisterRouter()
	GoogleLogin(c *fiber.Ctx) error
	GoogleCallback(c *fiber.Ctx) error
}

type oauthGoogleHandler struct {
	service service.OAuthGoogleService
	router  fiber.Router
	cfg     *config.Config
}

func NewOAuthGoogleHandler(
	service service.OAuthGoogleService,
	router fiber.Router,
	cfg *config.Config,
) OAuthGoogleHandler {
	return &oauthGoogleHandler{service, router, cfg}
}

func (h *oauthGoogleHandler) RegisterRouter() {
	router := h.router.Group("/auth")

	router.Get("/google/login", h.GoogleLogin)
	router.Get("/google/callback", h.GoogleCallback)
}

func (h *oauthGoogleHandler) GoogleLogin(c *fiber.Ctx) error {
	url := h.service.GetGoogleLoginUrl("state-token")
	return c.Redirect(url, fiber.StatusTemporaryRedirect)
}

func (h *oauthGoogleHandler) GoogleCallback(c *fiber.Ctx) error {
	code := c.Query("code")
	user, err := h.service.HandleGoogleCallback(code)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(
			response.Error(err.Code, err.Error()),
		)
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(user))
}

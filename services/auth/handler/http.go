package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/vasapolrittideah/money-tracker-api/services/auth/model"
	"github.com/vasapolrittideah/money-tracker-api/services/auth/service"
	"github.com/vasapolrittideah/money-tracker-api/shared/config"
	"github.com/vasapolrittideah/money-tracker-api/shared/middleware"
	"github.com/vasapolrittideah/money-tracker-api/shared/model/response"
	"google.golang.org/grpc/codes"
)

type AuthHttpHandler interface {
	RegisterRouter()
	SignUp(c *fiber.Ctx) error
	SignIn(c *fiber.Ctx) error
}

type authHttpHandler struct {
	service    service.AuthService
	middleware middleware.CoreMiddleware
	router     fiber.Router
	cfg        *config.Config
}

func NewAuthHttpHandler(
	service service.AuthService,
	middleware middleware.CoreMiddleware,
	router fiber.Router,
	cfg *config.Config,
) AuthHttpHandler {
	return &authHttpHandler{service, middleware, router, cfg}
}

func (h *authHttpHandler) RegisterRouter() {
	router := h.router.Group("/auth")

	router.Post("/sign-up", h.SignUp)
	router.Post("/sign-in", h.SignIn)
}

func (h *authHttpHandler) SignUp(c *fiber.Ctx) error {
	payload := new(model.SignUpRequest)

	if err := c.BodyParser(payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			response.Error(codes.InvalidArgument, err.Error()),
		)
	}

	res, err := h.service.SignUp(payload)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			response.Error(err.Code, err.Error()),
		)
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(res))
}

func (h *authHttpHandler) SignIn(c *fiber.Ctx) error {
	payload := new(model.SignInRequest)

	if err := c.BodyParser(payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			response.Error(codes.InvalidArgument, err.Error()),
		)
	}

	res, err := h.service.SignIn(payload)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			response.Error(err.Code, err.Error()),
		)
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(res))
}

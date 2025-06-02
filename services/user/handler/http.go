package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/vasapolrittideah/money-tracker-api/services/user/service"
	"github.com/vasapolrittideah/money-tracker-api/shared/config"
	"github.com/vasapolrittideah/money-tracker-api/shared/model/response"
	"google.golang.org/grpc/codes"
)

type UserHttpHandler interface {
	RegisterRouter()
	GetAllUsers(c *fiber.Ctx) error
	GetUserById(c *fiber.Ctx) error
	GetUserByEmail(c *fiber.Ctx) error
}

type userHttpHandler struct {
	router  fiber.Router
	service service.UserService
	cfg     *config.Config
}

func NewUserHttpHandler(
	router fiber.Router,
	service service.UserService,
	cfg *config.Config,
) UserHttpHandler {
	return &userHttpHandler{router, service, cfg}
}

func (h *userHttpHandler) RegisterRouter() {
	router := h.router.Group("/users")

	router.Get("/", h.GetAllUsers)
	router.Get("/:id", h.GetUserById)
	router.Get("/email/:email", h.GetUserByEmail)
}

func (h *userHttpHandler) GetAllUsers(c *fiber.Ctx) error {
	users, err := h.service.GetAllUsers()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			response.Error(err.Code, err.Error()),
		)
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(users))
}

func (h *userHttpHandler) GetUserById(c *fiber.Ctx) error {
	var id uuid.UUID
	if err := c.ParamsParser(&id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			response.Error(codes.InvalidArgument, err.Error()),
		)
	}

	user, err := h.service.GetUserById(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			response.Error(err.Code, err.Error()),
		)
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(user))
}

func (h *userHttpHandler) GetUserByEmail(c *fiber.Ctx) error {
	var email string
	if err := c.ParamsParser(&email); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			response.Error(codes.InvalidArgument, err.Error()),
		)
	}

	user, err := h.service.GetUserByEmail(email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			response.Error(err.Code, err.Error()),
		)
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(user))
}

package controller

import (
	"net/http"
	"net/mail"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/vasapolrittideah/money-tracker-api/shared/config"
	"github.com/vasapolrittideah/money-tracker-api/shared/domain"
	"github.com/vasapolrittideah/money-tracker-api/shared/httperror"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type userHTTPController struct {
	router  fiber.Router
	usecase domain.UserUsecase
	config  *config.Config
}

func NewUserHTTPController(
	router fiber.Router,
	usecase domain.UserUsecase,
	config *config.Config,
) *userHTTPController {
	return &userHTTPController{
		router:  router,
		usecase: usecase,
		config:  config,
	}
}

func (c *userHTTPController) RegisterRoutes() {
	router := c.router.Group("/users")

	router.Get("/", c.GetAllUsers)
	router.Get("/:id", c.GetUserByID)
	router.Get("/email/:email", c.GetUserByEmail)
	router.Post("/", c.CreateUser)
	router.Put("/:id", c.UpdateUser)
	router.Delete("/:id", c.DeleteUser)
}

func (c *userHTTPController) GetAllUsers(ctx *fiber.Ctx) error {
	users, err := c.usecase.GetAllUsers()
	if err != nil {
		st := status.Convert(err)
		return ctx.Status(httperror.HTTPStatusFromCode(st.Code())).JSON(
			httperror.NewHTTPError(st.Code(), st.Message()),
		)
	}

	return ctx.Status(http.StatusOK).JSON(users)
}

func (c *userHTTPController) GetUserByID(ctx *fiber.Ctx) error {
	idParam := ctx.Params("id")

	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(
			httperror.NewHTTPError(codes.InvalidArgument, "invalid user id format"),
		)
	}

	user, err := c.usecase.GetUserByID(id)
	if err != nil {
		st := status.Convert(err)
		return ctx.Status(httperror.HTTPStatusFromCode(st.Code())).JSON(
			httperror.NewHTTPError(st.Code(), st.Message()),
		)
	}

	return ctx.Status(http.StatusOK).JSON(user)
}

func (c *userHTTPController) GetUserByEmail(ctx *fiber.Ctx) error {
	email := ctx.Params("email")

	_, err := mail.ParseAddress(email)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(
			httperror.NewHTTPError(codes.InvalidArgument, "invalid email format"),
		)
	}

	user, err := c.usecase.GetUserByEmail(email)
	if err != nil {
		st := status.Convert(err)
		return ctx.Status(httperror.HTTPStatusFromCode(st.Code())).JSON(
			httperror.NewHTTPError(st.Code(), st.Message()),
		)
	}

	return ctx.Status(http.StatusOK).JSON(user)
}

func (c *userHTTPController) CreateUser(ctx *fiber.Ctx) error {
	var user domain.User
	if err := ctx.BodyParser(&user); err != nil {
		st := status.Convert(err)
		return ctx.Status(http.StatusBadGateway).JSON(
			httperror.NewHTTPError(codes.InvalidArgument, st.Message()),
		)
	}

	createdUser, err := c.usecase.CreateUser(&user)
	if err != nil {
		st := status.Convert(err)
		return ctx.Status(httperror.HTTPStatusFromCode(st.Code())).JSON(
			httperror.NewHTTPError(st.Code(), st.Message()),
		)
	}

	return ctx.Status(http.StatusOK).JSON(createdUser)
}

func (c *userHTTPController) UpdateUser(ctx *fiber.Ctx) error {
	idParam := ctx.Params("id")

	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(
			httperror.NewHTTPError(codes.InvalidArgument, "invalid user id format"),
		)
	}

	var user domain.User
	if err := ctx.BodyParser(&user); err != nil {
		st := status.Convert(err)
		return ctx.Status(http.StatusBadGateway).JSON(
			httperror.NewHTTPError(codes.InvalidArgument, st.Message()),
		)
	}

	updatedUser, err := c.usecase.UpdateUser(id, &user)
	if err != nil {
		st := status.Convert(err)
		return ctx.Status(httperror.HTTPStatusFromCode(st.Code())).JSON(
			httperror.NewHTTPError(st.Code(), st.Message()),
		)
	}

	return ctx.Status(http.StatusOK).JSON(updatedUser)
}

func (c *userHTTPController) DeleteUser(ctx *fiber.Ctx) error {
	idParam := ctx.Params("id")

	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(
			httperror.NewHTTPError(codes.InvalidArgument, "invalid user id format"),
		)
	}

	deletedUser, err := c.usecase.DeleteUser(id)
	if err != nil {
		st := status.Convert(err)
		return ctx.Status(httperror.HTTPStatusFromCode(st.Code())).JSON(
			httperror.NewHTTPError(st.Code(), st.Message()),
		)
	}

	return ctx.Status(http.StatusOK).JSON(deletedUser)
}

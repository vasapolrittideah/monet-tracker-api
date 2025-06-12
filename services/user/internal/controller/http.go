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
	usecase domain.UserUsecase
	router  fiber.Router
	config  *config.Config
}

func NewUserHTTPController(
	usecase domain.UserUsecase,
	router fiber.Router,
	config *config.Config,
) *userHTTPController {
	return &userHTTPController{
		usecase: usecase,
		router:  router,
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

// GetAllUsers godoc
// @Summary Get all users
// @Description get a list of all users
// @Tags User
// @Acceopt json
// @Produce json
// @Success 200 {array} domain.User "OK"
// @Failure 404 {object} httperror.HTTPError "Not Found"
// @Failure 500 {object} httperror.HTTPError "Internal Server Error"
// @Router /users [get]
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

// GetUserByID godoc
// @Summary Get user by id
// @Description get a user by id
// @Tags User
// @Acceopt json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} domain.User "OK"
// @Failure 400 {object} httperror.HTTPError "Bad Request"
// @Failure 404 {object} httperror.HTTPError "Not Found"
// @Failure 500 {object} httperror.HTTPError "Internal Server Error"
// @Router /users/{id} [get]
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

// GetUserByEmail godoc
// @Summary Get user by email
// @Description get a user by email
// @Tags User
// @Acceopt json
// @Produce json
// @Param email path string true "User Email"
// @Success 200 {object} domain.User "OK"
// @Failure 400 {object} httperror.HTTPError "Bad Request"
// @Failure 404 {object} httperror.HTTPError "Not Found"
// @Failure 500 {object} httperror.HTTPError "Internal Server Error"
// @Router /users/email/{email} [get]
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

// CreateUser godoc
// @Summary Create user
// @Description create a new user
// @Tags User
// @Acceopt json
// @Produce json
// @Param user body domain.CreateUserRequest true "User to create"
// @Success 200 {object} domain.User "OK"
// @Failure 400 {object} httperror.HTTPValidationError "Bad Request"
// @Failure 409 {object} httperror.HTTPError "Conflict"
// @Failure 500 {object} httperror.HTTPError "Internal Server Error"
// @Router /users [post]
func (c *userHTTPController) CreateUser(ctx *fiber.Ctx) error {
	var req domain.CreateUserRequest

	if err := ctx.BodyParser(&req); err != nil {
		st := status.Convert(err)
		return ctx.Status(http.StatusBadGateway).JSON(
			httperror.NewHTTPError(codes.InvalidArgument, st.Message()),
		)
	}

	user := domain.User{
		FullName: req.FullName,
		Email:    req.Email,
		Password: req.Password,
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

// UpdateUser godoc
// @Summary Update user
// @Description update a user
// @Tags User
// @Acceopt json
// @Produce json
// @Param id path string true "User ID"
// @Param user body domain.UpdateUserRequest true "User to update"
// @Success 200 {object} domain.User "OK"
// @Failure 400 {object} httperror.HTTPValidationError "Bad Request"
// @Failure 404 {object} httperror.HTTPError "Not Found"
// @Failure 500 {object} httperror.HTTPError "Internal Server Error"
// @Router /users/{id} [put]
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

// DeleteUser godoc
// @Summary Delete user
// @Description delete a user
// @Tags User
// @Acceopt json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} domain.User "OK"
// @Failure 400 {object} httperror.HTTPError "Bad Request"
// @Failure 404 {object} httperror.HTTPError "Not Found"
// @Failure 500 {object} httperror.HTTPError "Internal Server Error"
// @Router /users/{id} [delete]
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

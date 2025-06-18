package httphandler

import (
	"net/http"
	"net/mail"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/vasapolrittideah/money-tracker-api/shared/config"
	"github.com/vasapolrittideah/money-tracker-api/shared/domain"
	"github.com/vasapolrittideah/money-tracker-api/shared/errors/apperror"
	"github.com/vasapolrittideah/money-tracker-api/shared/errors/httperror"
	"github.com/vasapolrittideah/money-tracker-api/shared/validator"
)

type userHTTPHandler struct {
	usecase domain.UserUsecase
	router  fiber.Router
	config  *config.Config
}

func NewUserHTTPHandler(
	usecase domain.UserUsecase,
	router fiber.Router,
	config *config.Config,
) *userHTTPHandler {
	return &userHTTPHandler{
		usecase: usecase,
		router:  router,
		config:  config,
	}
}

func (h *userHTTPHandler) RegisterRoutes() {
	router := h.router.Group("/users")

	router.Get("/", h.GetAllUsers)
	router.Get("/:id", h.GetUserByID)
	router.Get("/email/:email", h.GetUserByEmail)
	router.Post("/", h.CreateUser)
	router.Put("/:id", h.UpdateUser)
	router.Delete("/:id", h.DeleteUser)
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
func (h *userHTTPHandler) GetAllUsers(c *fiber.Ctx) error {
	users, err := h.usecase.GetAllUsers()
	if err != nil {
		return httperror.FromAppError(c, err.(*apperror.AppError))
	}

	return c.Status(http.StatusOK).JSON(users)
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
func (h *userHTTPHandler) GetUserByID(c *fiber.Ctx) error {
	idParam := c.Params("id")

	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		return httperror.NewBadRequestError(c, "invalid user id format")
	}

	user, err := h.usecase.GetUserByID(id)
	if err != nil {
		return httperror.FromAppError(c, err.(*apperror.AppError))
	}

	return c.Status(http.StatusOK).JSON(user)
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
func (h *userHTTPHandler) GetUserByEmail(c *fiber.Ctx) error {
	email := c.Params("email")

	_, err := mail.ParseAddress(email)
	if err != nil {
		return httperror.NewBadRequestError(c, "invalid email format")
	}

	user, err := h.usecase.GetUserByEmail(email)
	if err != nil {
		return httperror.FromAppError(c, err.(*apperror.AppError))
	}

	return c.Status(http.StatusOK).JSON(user)
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
func (h *userHTTPHandler) CreateUser(c *fiber.Ctx) error {
	var req domain.CreateUserRequest

	if err := c.BodyParser(&req); err != nil {
		return httperror.NewBadRequestError(c, err.Error())
	}

	if err := validator.ValidateInput(c.Context(), req); err != nil {
		return httperror.NewValidationError(c, err)
	}

	user := domain.User{
		FullName: req.FullName,
		Email:    req.Email,
		Password: req.Password,
	}

	created, err := h.usecase.CreateUser(&user)
	if err != nil {
		return httperror.FromAppError(c, err.(*apperror.AppError))
	}

	return c.Status(http.StatusOK).JSON(created)
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
func (h *userHTTPHandler) UpdateUser(c *fiber.Ctx) error {
	idParam := c.Params("id")

	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		return httperror.NewBadRequestError(c, "invalid user id format")
	}

	var req domain.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return httperror.NewBadRequestError(c, err.Error())
	}

	if err := validator.ValidateInput(c.Context(), req); err != nil {
		return httperror.NewValidationError(c, err)
	}

	user, err := h.usecase.GetUserByID(id)
	if err != nil {
		return httperror.FromAppError(c, err.(*apperror.AppError))
	}

	if req.FullName != nil {
		user.FullName = *req.FullName
	}
	if req.Email != nil {
		user.Email = *req.Email
	}
	if req.Password != nil {
		user.Password = *req.Password
	}
	if req.Verified != nil {
		user.Verified = *req.Verified
	}
	if req.Registered != nil {
		user.Registered = *req.Registered
	}
	if req.RefreshToken != nil {
		user.RefreshToken = *req.RefreshToken
	}

	updated, err := h.usecase.UpdateUser(user)
	if err != nil {
		return httperror.FromAppError(c, err.(*apperror.AppError))
	}

	return c.Status(http.StatusOK).JSON(updated)
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
func (h *userHTTPHandler) DeleteUser(c *fiber.Ctx) error {
	idParam := c.Params("id")

	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		return httperror.NewBadRequestError(c, "invalid user id format")
	}

	deleted, err := h.usecase.DeleteUser(id)
	if err != nil {
		return httperror.FromAppError(c, err.(*apperror.AppError))
	}

	return c.Status(http.StatusOK).JSON(deleted)
}

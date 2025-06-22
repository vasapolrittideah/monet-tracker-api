package httphandler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	auth "github.com/vasapolrittideah/money-tracker-api/services/auth/internal"
	"github.com/vasapolrittideah/money-tracker-api/shared/config"
	"github.com/vasapolrittideah/money-tracker-api/shared/errors/apperror"
	"github.com/vasapolrittideah/money-tracker-api/shared/errors/httperror"
	"github.com/vasapolrittideah/money-tracker-api/shared/validator"
)

type authHTTPHandler struct {
	usecase auth.AuthUsecase
	router  fiber.Router
	config  *config.Config
}

func NewAuthHTTPHandler(usecase auth.AuthUsecase, router fiber.Router, config *config.Config) *authHTTPHandler {
	return &authHTTPHandler{
		usecase: usecase,
		router:  router,
		config:  config,
	}
}

func (c *authHTTPHandler) RegisterRoutes() {
	router := c.router.Group("/auth")

	router.Post("/sign-up", c.SignUp)
	router.Post("/sign-in", c.SignIn)
}

// SignUp godoc
// @Summary Sign Up
// @Description register a new user
// @Tags Auth
// @Acceopt json
// @Produce json
// @Param user body domain.SignUpRequest true "User to register"
// @Success 200 {object} domain.User "OK"
// @Failure 400 {object} httperror.HTTPValidationError "Bad Request"
// @Failure 409 {object} httperror.HTTPError "Conflict"
// @Failure 500 {object} httperror.HTTPError "Internal Server Error"
// @Router /auth/sign-up [post]
func (h *authHTTPHandler) SignUp(c *fiber.Ctx) error {
	req := new(auth.SignUpRequest)

	if err := c.BodyParser(req); err != nil {
		return httperror.NewBadRequestError(c, err.Error())
	}

	if err := validator.ValidateInput(c.Context(), req); err != nil {
		return httperror.NewValidationError(c, err)
	}

	user, err := h.usecase.SignUp(c.Context(), req)
	if err != nil {
		return httperror.FromAppError(c, err.(*apperror.AppError))
	}

	return c.Status(http.StatusOK).JSON(user)
}

// SignIn godoc
// @Summary Sign In
// @Description sign in a user
// @Tags Auth
// @Acceopt json
// @Produce json
// @Param user body domain.SignInRequest true "User to sign in"
// @Success 200 {object} domain.Token "OK"
// @Failure 400 {object} httperror.HTTPValidationError "Bad Request"
// @Failure 401 {object} httperror.HTTPError "Unauthorized"
// @Failure 500 {object} httperror.HTTPError "Internal Server Error"
// @Router /auth/sign-in [post]
func (h *authHTTPHandler) SignIn(c *fiber.Ctx) error {
	req := new(auth.SignInRequest)

	if err := c.BodyParser(req); err != nil {
		return httperror.NewBadRequestError(c, err.Error())
	}

	if err := validator.ValidateInput(c.Context(), req); err != nil {
		return httperror.NewValidationError(c, err)
	}

	token, err := h.usecase.SignIn(c.Context(), req)
	if err != nil {
		return httperror.FromAppError(c, err.(*apperror.AppError))
	}

	return c.Status(http.StatusOK).JSON(token)
}

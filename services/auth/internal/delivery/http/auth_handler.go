package httphandler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/vasapolrittideah/money-tracker-api/shared/config"
	"github.com/vasapolrittideah/money-tracker-api/shared/domain"
	"github.com/vasapolrittideah/money-tracker-api/shared/httperror"
	"github.com/vasapolrittideah/money-tracker-api/shared/validator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type authHTTPHandler struct {
	usecase domain.AuthUsecase
	router  fiber.Router
	config  *config.Config
}

func NewAuthHTTPHandler(usecase domain.AuthUsecase, router fiber.Router, config *config.Config) *authHTTPHandler {
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
func (c *authHTTPHandler) SignUp(ctx *fiber.Ctx) error {
	req := new(domain.SignUpRequest)

	if err := ctx.BodyParser(req); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(
			httperror.NewHTTPError(codes.InvalidArgument, err.Error()),
		)
	}

	if err := validator.ValidateInput(ctx.Context(), req); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(
			httperror.NewValidationError(err.Details),
		)
	}

	user, err := c.usecase.SignUp(req)
	if err != nil {
		st := status.Convert(err)
		return ctx.Status(httperror.HTTPStatusFromCode(st.Code())).JSON(
			httperror.NewHTTPError(st.Code(), st.Message()),
		)
	}

	return ctx.Status(http.StatusOK).JSON(user)
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
func (c *authHTTPHandler) SignIn(ctx *fiber.Ctx) error {
	req := new(domain.SignInRequest)

	if err := ctx.BodyParser(req); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(
			httperror.NewHTTPError(codes.InvalidArgument, err.Error()),
		)
	}

	if err := validator.ValidateInput(ctx.Context(), req); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(
			httperror.NewValidationError(err.Details),
		)
	}

	token, err := c.usecase.SignIn(req)
	if err != nil {
		st := status.Convert(err)
		return ctx.Status(httperror.HTTPStatusFromCode(st.Code())).JSON(
			httperror.NewHTTPError(st.Code(), st.Message()),
		)
	}

	return ctx.Status(http.StatusOK).JSON(token)
}

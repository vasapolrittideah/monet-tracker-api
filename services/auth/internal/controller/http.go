package controller

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/vasapolrittideah/money-tracker-api/shared/config"
	"github.com/vasapolrittideah/money-tracker-api/shared/domain"
	"github.com/vasapolrittideah/money-tracker-api/shared/httperror"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type authHTTPController struct {
	usecase domain.AuthUsecase
	router  fiber.Router
	config  *config.Config
}

func NewAuthHTTPController(usecase domain.AuthUsecase, router fiber.Router, config *config.Config) *authHTTPController {
	return &authHTTPController{
		usecase: usecase,
		router:  router,
		config:  config,
	}
}

func (c *authHTTPController) RegisterRoutes() {
	router := c.router.Group("/auth")

	router.Post("/sign-up", c.SignUp)
	router.Post("/sign-in", c.SignIn)
}

func (c *authHTTPController) SignUp(ctx *fiber.Ctx) error {
	req := new(domain.SignUpRequest)

	if err := ctx.BodyParser(req); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(
			httperror.NewHTTPError(codes.InvalidArgument, err.Error()),
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

func (c *authHTTPController) SignIn(ctx *fiber.Ctx) error {
	req := new(domain.SignInRequest)

	if err := ctx.BodyParser(req); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(
			httperror.NewHTTPError(codes.InvalidArgument, err.Error()),
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

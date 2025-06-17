package httphandler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/vasapolrittideah/money-tracker-api/shared/config"
	"github.com/vasapolrittideah/money-tracker-api/shared/domain"
	"github.com/vasapolrittideah/money-tracker-api/shared/httperror"
	"google.golang.org/grpc/codes"
)

type oauthGoogleHTTPHandler struct {
	usecase domain.OAuthGoogleUsecase
	router  fiber.Router
	config  *config.Config
}

func NewOAuthGoogleHTTPHandler(
	usecase domain.OAuthGoogleUsecase,
	router fiber.Router,
	config *config.Config,
) *oauthGoogleHTTPHandler {
	return &oauthGoogleHTTPHandler{
		usecase: usecase,
		router:  router,
		config:  config,
	}
}

func (c *oauthGoogleHTTPHandler) RegisterRoutes() {
	router := c.router.Group("/auth")

	router.Get("/google", c.SignInWithGoogle)
	router.Get("/google/callback", c.HandleGoogleCallback)
}

func (c *oauthGoogleHTTPHandler) SignInWithGoogle(ctx *fiber.Ctx) error {
	url := c.usecase.GetSignInWithGoogleURL("state-token")
	return ctx.Redirect(url, fiber.StatusTemporaryRedirect)
}

func (c *oauthGoogleHTTPHandler) HandleGoogleCallback(ctx *fiber.Ctx) error {
	code := ctx.Query("code")
	token, err := c.usecase.HandleGoogleCallback(code)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(
			httperror.NewHTTPError(codes.Internal, err.Error()),
		)
	}

	return ctx.Status(http.StatusOK).JSON(token)
}

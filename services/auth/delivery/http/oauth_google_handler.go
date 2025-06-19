package httphandler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/vasapolrittideah/money-tracker-api/shared/config"
	"github.com/vasapolrittideah/money-tracker-api/shared/domain"
	"github.com/vasapolrittideah/money-tracker-api/shared/errors/apperror"
	"github.com/vasapolrittideah/money-tracker-api/shared/errors/httperror"
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

func (h *oauthGoogleHTTPHandler) SignInWithGoogle(c *fiber.Ctx) error {
	url := h.usecase.GetSignInWithGoogleURL("state-token")
	return c.Redirect(url, fiber.StatusTemporaryRedirect)
}

func (h *oauthGoogleHTTPHandler) HandleGoogleCallback(c *fiber.Ctx) error {
	code := c.Query("code")
	token, err := h.usecase.HandleGoogleCallback(code)
	if err != nil {
		return httperror.FromAppError(c, err.(*apperror.AppError))
	}

	return c.Status(http.StatusOK).JSON(token)
}

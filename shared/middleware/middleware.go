package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/vasapolrittideah/money-tracker-api/shared/config"
	"github.com/vasapolrittideah/money-tracker-api/shared/domain/response"
	"github.com/vasapolrittideah/money-tracker-api/shared/utils/jwtutil"
	"google.golang.org/grpc/codes"
)

type CoreMiddleware interface {
	Authenticate(tokenType TokenType) fiber.Handler
}

type coreMiddleware struct {
	cfg *config.Config
}

func NewCoreMiddleware(cfg *config.Config) CoreMiddleware {
	return &coreMiddleware{cfg}
}

type TokenType int

const (
	AccessToken TokenType = iota
	RefreshToken
)

func (m coreMiddleware) Authenticate(tokenType TokenType) fiber.Handler {
	return func(c *fiber.Ctx) error {
		const bearer = "Bearer"
		token := c.Get("Authorization")
		if token == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(
				response.Error(codes.Unauthenticated, "No Authorization header found"),
			)
		}

		headerParts := strings.Split(token, " ")
		if len(headerParts) != 2 || headerParts[0] != bearer {
			return c.Status(fiber.StatusUnauthorized).JSON(
				response.Error(codes.Unauthenticated, "Malformed Authorization header"),
			)
		}

		var secretKey string
		switch tokenType {
		case AccessToken:
			secretKey = m.cfg.Jwt.AccessTokenSecretKey
		case RefreshToken:
			secretKey = m.cfg.Jwt.RefreshTokenSecretKey
		}

		claims, err := jwtutil.ParseToken(headerParts[1], secretKey)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(
				response.Error(codes.Unauthenticated, err.Error()),
			)
		}

		c.Locals("token", headerParts[1])
		c.Locals("sub", (*claims)["sub"])

		return c.Next()
	}
}

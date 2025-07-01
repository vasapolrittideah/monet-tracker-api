package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/vasapolrittideah/money-tracker-api/shared/errors/httperror"
	"github.com/vasapolrittideah/money-tracker-api/shared/utils/tokenutil"
)

func ValidateToken(secret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			return httperror.NewUnauthorizedError(c, "missing access token")
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := tokenutil.ValidateToken(tokenStr, secret)
		if err != nil {
			return httperror.NewUnauthorizedError(c, "invalid access token")
		}

		c.Locals("user_id", claims.UserID)
		c.Locals("session_id", claims.SessionID)

		return c.Next()
	}
}

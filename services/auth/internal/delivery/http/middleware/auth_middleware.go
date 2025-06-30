package middleware

import (
	"errors"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	auth "github.com/vasapolrittideah/money-tracker-api/services/auth/internal"
	"github.com/vasapolrittideah/money-tracker-api/shared/config"
	"github.com/vasapolrittideah/money-tracker-api/shared/errors/httperror"
	"github.com/vasapolrittideah/money-tracker-api/shared/utils/hashutil"
	"github.com/vasapolrittideah/money-tracker-api/shared/utils/tokenutil"
	"gorm.io/gorm"
)

type AuthMiddleware struct {
	sessionRepo auth.SessionRepository
	jwtConfig   *config.JWTConfig
}

func NewAuthMiddleware(sessionRepo auth.SessionRepository, jwtConfig *config.JWTConfig) *AuthMiddleware {
	return &AuthMiddleware{
		sessionRepo: sessionRepo,
		jwtConfig:   jwtConfig,
	}
}

func (m *AuthMiddleware) ValidateRefreshToken() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			return httperror.NewUnauthorizedError(c, "missing refresh token")
		}

		refreshToken := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := tokenutil.ValidateToken(refreshToken, m.jwtConfig.RefreshTokenSecretKey)
		if err != nil {
			return httperror.NewUnauthorizedError(c, "invalid refresh token")
		}

		session, err := m.sessionRepo.GetSessionByID(c.Context(), claims.SessionID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return httperror.NewUnauthorizedError(c, "session not found")
			}

			return httperror.NewUnauthorizedError(c, "failed to get session")
		}

		if ok, err := hashutil.Verify(refreshToken, session.Token); err != nil || !ok {
			return httperror.NewUnauthorizedError(c, "invalid refresh token")
		}

		if session.Revoked || session.ExpiresAt.Before(time.Now()) {
			return httperror.NewUnauthorizedError(c, "session expired or revoked")
		}

		c.Locals("session", session)

		return c.Next()
	}
}

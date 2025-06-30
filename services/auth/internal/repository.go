package auth

import (
	"context"

	"github.com/vasapolrittideah/money-tracker-api/shared/domain"
)

type ExternalAuthRepository interface {
	GetExternalAuthByProvider(ctx context.Context, provider string, providerID string) (*domain.ExternalAuth, error)
	CreateExternalAuth(ctx context.Context, externalAuth *domain.ExternalAuth) (*domain.ExternalAuth, error)
	DeleteExternalAuth(ctx context.Context, id uint64) (*domain.ExternalAuth, error)
}

type SessionRepository interface {
	GetSessionByID(ctx context.Context, id uint64) (*domain.Session, error)
	GetSessionByToken(ctx context.Context, token string) (*domain.Session, error)
	CreateSession(ctx context.Context, session *domain.Session) (*domain.Session, error)
	UpdateSession(ctx context.Context, id uint64, updates map[string]any) (*domain.Session, error)
	DeleteSessionByID(ctx context.Context, id uint64) (*domain.Session, error)
	DeleteSessionByUserID(ctx context.Context, userID uint64) (*domain.Session, error)
	RevokeSession(ctx context.Context, id uint64) (*domain.Session, error)
}

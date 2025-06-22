package auth

import (
	"context"

	"github.com/vasapolrittideah/money-tracker-api/shared/domain"
)

type AuthRepository interface {
	GetExternalAuthByProviderID(ctx context.Context, providerID string) (*domain.ExternalAuth, error)
	CreateExternalAuth(ctx context.Context, externalAuth *domain.ExternalAuth) (*domain.ExternalAuth, error)
	UpdateExternalAuth(ctx context.Context, id uint64, externalAuth *domain.ExternalAuth) (*domain.ExternalAuth, error)
	DeleteExternalAuth(ctx context.Context, id uint64) (*domain.ExternalAuth, error)
}

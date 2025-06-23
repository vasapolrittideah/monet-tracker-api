package repository

import (
	"context"

	auth "github.com/vasapolrittideah/money-tracker-api/services/auth/internal"
	"github.com/vasapolrittideah/money-tracker-api/shared/domain"
	"gorm.io/gorm"
)

type externalAuthRepository struct {
	db *gorm.DB
}

func NewExternalAuthRepository(db *gorm.DB) auth.ExternalAuthRepository {
	return &externalAuthRepository{db}
}

func (r *externalAuthRepository) GetExternalAuthByProvider(
	ctx context.Context,
	provider string,
	providerID string,
) (*domain.ExternalAuth, error) {
	var auth domain.ExternalAuth
	if err := r.db.WithContext(ctx).First(&auth, "provider = ? AND provider_id = ?", provider, providerID).Error; err != nil {
		return nil, err
	}

	return &auth, nil
}

func (r *externalAuthRepository) CreateExternalAuth(
	ctx context.Context,
	externalAuth *domain.ExternalAuth,
) (*domain.ExternalAuth, error) {
	if err := r.db.WithContext(ctx).Create(externalAuth).Error; err != nil {
		return nil, err
	}

	return externalAuth, nil
}

func (r *externalAuthRepository) DeleteExternalAuth(ctx context.Context, id uint64) (*domain.ExternalAuth, error) {
	var auth domain.ExternalAuth
	if err := r.db.WithContext(ctx).First(&auth, "id = ?", id).Error; err != nil {
		return nil, err
	}
	if err := r.db.WithContext(ctx).Delete(&auth).Error; err != nil {
		return nil, err
	}

	return &auth, nil
}

package repository

import (
	"context"

	"github.com/vasapolrittideah/money-tracker-api/shared/domain"
	"gorm.io/gorm"
)

type authRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) domain.AuthRepository {
	return &authRepository{db}
}

func (r *authRepository) GetExternalAuthByProviderID(
	ctx context.Context,
	providerID string,
) (*domain.ExternalAuth, error) {
	var auth domain.ExternalAuth
	if err := r.db.WithContext(ctx).First(&auth, "provider_id = ?", providerID).Error; err != nil {
		return nil, err
	}

	return &auth, nil
}

func (r *authRepository) CreateExternalAuth(
	ctx context.Context,
	externalAuth *domain.ExternalAuth,
) (*domain.ExternalAuth, error) {
	if err := r.db.WithContext(ctx).Create(externalAuth).Error; err != nil {
		return nil, err
	}

	return externalAuth, nil
}

func (r *authRepository) UpdateExternalAuth(
	ctx context.Context,
	id uint64,
	externalAuth *domain.ExternalAuth,
) (*domain.ExternalAuth, error) {
	var auth domain.ExternalAuth
	if err := r.db.WithContext(ctx).First(&auth, "id = ?", id).Error; err != nil {
		return nil, err
	}
	if err := r.db.WithContext(ctx).Model(&auth).Where("id = ?", id).Updates(externalAuth).Error; err != nil {
		return nil, err
	}
	if err := r.db.WithContext(ctx).First(&auth, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return &auth, nil
}

func (r *authRepository) DeleteExternalAuth(ctx context.Context, id uint64) (*domain.ExternalAuth, error) {
	var auth domain.ExternalAuth
	if err := r.db.WithContext(ctx).First(&auth, "id = ?", id).Error; err != nil {
		return nil, err
	}
	if err := r.db.WithContext(ctx).Delete(&auth).Error; err != nil {
		return nil, err
	}

	return &auth, nil
}

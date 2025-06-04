package repository

import (
	"github.com/vasapolrittideah/money-tracker-api/shared/model/apperror"
	"github.com/vasapolrittideah/money-tracker-api/shared/model/domain"
	"github.com/vasapolrittideah/money-tracker-api/shared/utils/errorutil"
	"gorm.io/gorm"
)

type AuthRepository interface {
	GetExternalLoginByProviderId(providerId string) (*domain.ExternalLogin, *apperror.Error)
	CreateExternalLogin(externalLogin *domain.ExternalLogin) (*domain.ExternalLogin, *apperror.Error)
}

type authRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &authRepository{db}
}

func (r *authRepository) GetExternalLoginByProviderId(providerId string) (*domain.ExternalLogin, *apperror.Error) {
	var externalLogin *domain.ExternalLogin
	if err := r.db.First(&externalLogin, "provider_id = ?", providerId).Error; err != nil {
		return nil, errorutil.HandleRecordNotFoundError(err)
	}

	return externalLogin, nil
}

func (r *authRepository) CreateExternalLogin(
	externalLogin *domain.ExternalLogin,
) (*domain.ExternalLogin, *apperror.Error) {
	if err := r.db.Create(&externalLogin).Error; err != nil {
		return nil, errorutil.HandleUnqiueConstraintError(err)
	}

	return externalLogin, nil
}

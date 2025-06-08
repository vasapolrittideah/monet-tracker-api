package domain

import "gorm.io/gorm"

type ExternalAuth struct {
	gorm.Model
	Provider   string `gorm:"not null"`
	ProviderID string `gorm:"not null"`
	UserID     uint
}

type Token struct {
	AccessToken  string
	RefreshToken string
}

type AuthRepository interface {
	GetExternalAuthByProviderID(providerID string) (*ExternalAuth, error)
	CreateExternalAuth(externalAuth *ExternalAuth) (*ExternalAuth, error)
	UpdateExternalAuth(id uint64, externalAuth *ExternalAuth) (*ExternalAuth, error)
	DeleteExternalAuth(id uint64) (*ExternalAuth, error)
}

type AuthUsecase interface {
	SignUp(req *SignUpRequest) (*User, error)
	SignIn(req *SignInRequest) (*Token, error)
}

type SignUpRequest struct {
	FullName string
	Email    string
	Password string
}

type SignInRequest struct {
	Email    string
	Password string
}

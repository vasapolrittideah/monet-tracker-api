package domain

import (
	"context"
	"time"
)

type ExternalAuth struct {
	ID         uint64 `json:"id"         gorm:"primaryKey;autoIncrement"`
	Provider   string `                  gorm:"not null"`
	ProviderID string `                  gorm:"not null"`
	UserID     uint64
	CreatedAt  time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type Token struct {
	AccessToken  string `json:"access_token"  extensions:"x-order=1"`
	RefreshToken string `json:"refresh_token" extensions:"x-order=2"`
}

// @swaggerignore
type SignUpRequest struct {
	FullName string `json:"full_name" example:"John Doe"         extensions:"x-order=1"`
	Email    string `json:"email"     example:"john@example.com" extensions:"x-order=2"`
	Password string `json:"password"  example:"password"         extensions:"x-order=3"`
}

// @swaggerignore
type SignInRequest struct {
	Email    string `json:"email"    example:"john@example.com" extensions:"x-order=2"`
	Password string `json:"password" example:"password"         extensions:"x-order=3"`
}

type AuthRepository interface {
	GetExternalAuthByProviderID(ctx context.Context, providerID string) (*ExternalAuth, error)
	CreateExternalAuth(ctx context.Context, externalAuth *ExternalAuth) (*ExternalAuth, error)
	UpdateExternalAuth(ctx context.Context, id uint64, externalAuth *ExternalAuth) (*ExternalAuth, error)
	DeleteExternalAuth(ctx context.Context, id uint64) (*ExternalAuth, error)
}

type AuthUsecase interface {
	SignUp(ctx context.Context, req *SignUpRequest) (*User, error)
	SignIn(ctx context.Context, req *SignInRequest) (*Token, error)
}

type OAuthGoogleUsecase interface {
	GetSignInWithGoogleURL(state string) string
	HandleGoogleCallback(ctx context.Context, code string) (*Token, error)
}

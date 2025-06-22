package auth

import (
	"context"

	"github.com/vasapolrittideah/money-tracker-api/shared/domain"
)

type AuthUsecase interface {
	SignUp(ctx context.Context, req *SignUpRequest) (*domain.User, error)
	SignIn(ctx context.Context, req *SignInRequest) (*TokenResponse, error)
}

type OAuthGoogleUsecase interface {
	GetSignInWithGoogleURL(state string) string
	HandleGoogleCallback(ctx context.Context, code string) (*TokenResponse, error)
}

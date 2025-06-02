package model

import "github.com/vasapolrittideah/money-tracker-api/shared/model/domain"

type SignUpRequest struct {
	FullName string `json:"full_name" validate:"required"`
	Email    string `json:"email"     validate:"required,email"`
	Password string `json:"password"  validate:"required"`
}

type SignUpResponse struct {
	*domain.User
}

type SignInRequest struct {
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type SignInResponse struct {
	*domain.Jwt
}

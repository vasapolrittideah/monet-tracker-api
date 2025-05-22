package model

import "github.com/vasapolrittideah/money-tracker-api/shared/domain/entity"

type SignUpRequest struct {
	FullName string `json:"full_name" validate:"required"`
	Email    string `json:"email"     validate:"required,email"`
	Password string `json:"password"  validate:"required"`
}

type SignUpResponse struct {
	*entity.User
}

type SignInRequest struct {
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type SignInResponse struct {
	*entity.Jwt
}

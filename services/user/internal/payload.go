package user

import (
	userpbv1 "github.com/vasapolrittideah/money-tracker-api/protogen/user/v1"
	"github.com/vasapolrittideah/money-tracker-api/shared/utils/protoutil"
)

type UpdateUserRequest struct {
	FullName   *string `json:"full_name"  example:"John Doe"         extensions:"x-order=1"`
	Email      *string `json:"email"      example:"john@example.com" extensions:"x-order=2" validate:"omitempty,email"`
	Password   *string `json:"password"   example:"securepassword"   extensions:"x-order=3"`
	Verified   *bool   `json:"verified"   example:"true"             extensions:"x-order=4"`
	Registered *bool   `json:"registered" example:"true"             extensions:"x-order=5"`
}

func NewUpdateUserRequestFromProto(req *userpbv1.UpdateUserRequest) *UpdateUserRequest {
	return &UpdateUserRequest{
		FullName:   protoutil.UnwrapString(req.FullName),
		Email:      protoutil.UnwrapString(req.Email),
		Password:   protoutil.UnwrapString(req.Password),
		Verified:   protoutil.UnwrapBool(req.Verified),
		Registered: protoutil.UnwrapBool(req.Registered),
	}
}

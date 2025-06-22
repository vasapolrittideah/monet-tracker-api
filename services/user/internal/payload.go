package user

type UpdateUserRequest struct {
	FullName     *string `json:"full_name"     example:"John Doe"         extensions:"x-order=1"`
	Email        *string `json:"email"         example:"john@example.com" extensions:"x-order=2" validate:"omitempty,email"`
	Password     *string `json:"password"      example:"securepassword"   extensions:"x-order=3"`
	Verified     *bool   `json:"verified"      example:"true"             extensions:"x-order=4"`
	Registered   *bool   `json:"registered"    example:"true"             extensions:"x-order=5"`
	RefreshToken *string `json:"refresh_token" example:"refresh_token"    extensions:"x-order=6"`
}

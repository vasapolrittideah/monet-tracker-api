package auth

type SignUpRequest struct {
	FullName string `json:"full_name" example:"John Doe"         extensions:"x-order=1"`
	Email    string `json:"email"     example:"john@example.com" extensions:"x-order=2"`
	Password string `json:"password"  example:"password"         extensions:"x-order=3"`
}

type SignInRequest struct {
	Email    string `json:"email"    example:"john@example.com" extensions:"x-order=2"`
	Password string `json:"password" example:"password"         extensions:"x-order=3"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"  extensions:"x-order=1"`
	RefreshToken string `json:"refresh_token" extensions:"x-order=2"`
}

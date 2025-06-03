package config

import (
	"time"

	"github.com/caarlos0/env/v11"
)

type JwtConfig struct {
	AccessTokenSecretKey  string        `env:"ACCESS_TOKEN_SECRET_KEY"`
	AccessTokenExpiresIn  time.Duration `env:"ACCESS_TOKEN_EXPIRES_IN"`
	RefreshTokenSecretKey string        `env:"REFRESH_TOKEN_SECRET_KEY"`
	RefreshTokenExpiresIn time.Duration `env:"REFRESH_TOKEN_EXPIRES_IN"`
}

type DatabaseConfig struct {
	Host     string `env:"POSTGRES_HOST"`
	Port     string `env:"POSTGRES_PORT"`
	User     string `env:"POSTGRES_USER"`
	Password string `env:"POSTGRES_PASSWORD"`
	Name     string `env:"POSTGRES_DB"`
}

type ServerConfig struct {
	UserServiceHttpPort           string `env:"USER_SERVICE_HTTP_PORT"`
	UserServiceGrpcPort           string `env:"USER_SERVICE_GRPC_PORT"`
	UserServiceGrpcConnectionHost string `env:"USER_SERVICE_GRPC_CONNECTION_HOST"`
	AuthServiceHttpPort           string `env:"AUTH_SERVICE_HTTP_PORT"`
}

type OAuthGoogleConfig struct {
	ClientId     string `env:"OAUTH_GOOGLE_CLIENT_ID"`
	ClientSecret string `env:"OAUTH_GOOGLE_CLIENT_SECRET"`
	RedirectUrl  string `env:"OAUTH_GOOGLE_REDIRECT_URL"`
}

type Config struct {
	Environment string `env:"ENVIRONMENT"`
	Jwt         JwtConfig
	Server      ServerConfig
	Database    DatabaseConfig
	OAuthGoogle OAuthGoogleConfig
}

func Load() (*Config, error) {
	var config Config
	if err := env.Parse(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

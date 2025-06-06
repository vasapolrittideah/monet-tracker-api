package config

import (
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/charmbracelet/log"
)

type JWTConfig struct {
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
	GatewayPort     string `env:"GATEWAY_PORT"`
	UserServiceHost string `env:"USER_SERVICE_HOST"`
	UserServicePort string `env:"USER_SERVICE_PORT"`
	AuthServiceHost string `env:"AUTH_SERVICE_HOST"`
	AuthServicePort string `env:"AUTH_SERVICE_PORT"`
}

type OAuthGoogleConfig struct {
	ClientID     string `env:"OAUTH_GOOGLE_CLIENT_ID"`
	ClientSecret string `env:"OAUTH_GOOGLE_CLIENT_SECRET"`
	RedirectURL  string `env:"OAUTH_GOOGLE_REDIRECT_URL"`
}

type Config struct {
	Environment string `env:"ENVIRONMENT"`
	JWT         JWTConfig
	Server      ServerConfig
	Database    DatabaseConfig
	OAuthGoogle OAuthGoogleConfig
}

func Load() *Config {
	var config Config
	if err := env.Parse(&config); err != nil {
		log.Fatal("failed to parse environment variables: %v", err)
	}

	return &config
}

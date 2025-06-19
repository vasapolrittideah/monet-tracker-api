package config

import (
	"time"

	"github.com/bazelbuild/rules_go/go/runfiles"
	"github.com/caarlos0/env/v11"
	"github.com/charmbracelet/log"
	"github.com/joho/godotenv"
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
	ConsulHost          string `env:"CONSUL_HOST"`
	ConsulPort          string `env:"CONSUL_PORT"`
	UserServiceHost     string `env:"USER_SERVICE_HOST"`
	UserServiceHTTPPort string `env:"USER_SERVICE_HTTP_PORT"`
	UserServiceGRPCPort string `env:"USER_SERVICE_GRPC_PORT"`
	AuthServiceHost     string `env:"AUTH_SERVICE_HOST"`
	AuthServiceHTTPPort string `env:"AUTH_SERVICE_HTTP_PORT"`
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
	rf, err := runfiles.New()
	if err != nil {
		log.Fatalf("failed to create runfiles object: %v", err)
	}

	envPath, err := rf.Rlocation("money_tracker_api/.env")
	if err != nil {
		log.Fatalf("failed to locate .env file: %v", err)
	}

	if err := godotenv.Load(envPath); err != nil {
		log.Fatalf("failed to load .env file: %v", err)
	}

	var config Config
	if err := env.Parse(&config); err != nil {
		log.Fatal("failed to parse environment variables: %v", err)
	}

	return &config
}

package config

import (
	"time"

	"github.com/spf13/viper"
	"github.com/vasapolrittideah/money-tracker-api/shared/logger"
	"github.com/vasapolrittideah/money-tracker-api/shared/utils/pathutil"
)

type JwtConfig struct {
	AccessTokenSecretKey  string        `mapstructure:"ACCESS_TOKEN_SECRET_KEY"`
	AccessTokenExpiresIn  time.Duration `mapstructure:"ACCESS_TOKEN_EXPIRES_IN"`
	RefreshTokenSecretKey string        `mapstructure:"REFRESH_TOKEN_SECRET_KEY"`
	RefreshTokenExpiresIn time.Duration `mapstructure:"REFRESH_TOKEN_EXPIRES_IN"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"POSTGRES_HOST"`
	Port     string `mapstructure:"POSTGRES_PORT"`
	User     string `mapstructure:"POSTGRES_USER"`
	Password string `mapstructure:"POSTGRES_PASSWORD"`
	Name     string `mapstructure:"POSTGRES_DB"`
}

type ServerConfig struct {
	AuthHttpPort string `mapstructure:"AUTH_SERVICE_HTTP_PORT"`
	UserHttpPort string `mapstructure:"USER_SERVICE_HTTP_PORT"`
	UserGrpcPort string `mapstructure:"USER_SERVICE_GRPC_PORT"`
}

type Config struct {
	Environment string         `mapstructure:"ENVIRONMENT"`
	Jwt         JwtConfig      `mapstructure:",squash"`
	Server      ServerConfig   `mapstructure:",squash"`
	Database    DatabaseConfig `mapstructure:",squash"`
}

func Load() (config *Config, err error) {
	rootDir, err := pathutil.GetProjectRoot()
	if err != nil || rootDir == "" {
		logger.Fatal("CORE", "failed to get project root: %v", err)
	}

	viper.SetConfigType("env")
	viper.SetConfigName(".env")
	viper.AddConfigPath(rootDir)

	if err := viper.ReadInConfig(); err != nil {
		logger.Fatal("CORE", "unable to read config file: %v", err)
	}

	if err := viper.Unmarshal(&config); err != nil {
		logger.Fatal("CORE", "unable to decode into struct: %v", err)
	}

	return
}

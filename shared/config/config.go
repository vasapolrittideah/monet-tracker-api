package config

import (
	"time"

	"github.com/spf13/viper"
	"vasapolrittideah/money-tracker-api/shared/logger"
	"vasapolrittideah/money-tracker-api/shared/utils/pathutil"
)

type JwtConfig struct {
	AccessTokenPrivateKey  string        `mapstructure:"ACCESS_TOKEN_PRIVATE_KEY"`
	AccessTokenPublicKey   string        `mapstructure:"ACCESS_TOKEN_PUBLIC_KEY"`
	AccessTokenExpiresIn   time.Duration `mapstructure:"ACCESS_TOKEN_EXPIRES_IN"`
	AccessTokenMaxAge      string        `mapstructure:"ACCESS_TOKEN_MAX_AGE"`
	RefreshTokenPrivateKey string        `mapstructure:"REFRESH_TOKEN_PRIVATE_KEY"`
	RefreshTokenPublicKey  string        `mapstructure:"REFRESH_TOKEN_PUBLIC_KEY"`
	RefreshTokenExpiresIn  time.Duration `mapstructure:"REFRESH_TOKEN_EXPIRES_IN"`
	RefreshTokenMaxAge     string        `mapstructure:"REFRESH_TOKEN_MAX_AGE"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"DB_HOST"`
	Port     string `mapstructure:"DB_PORT"`
	User     string `mapstructure:"DB_USER"`
	Password string `mapstructure:"DB_PASSWORD"`
	Name     string `mapstructure:"DB_NAME"`
}

type ServerConfig struct {
	AuthHttpPort string `mapstructure:"AUTH_HTTP_PORT"`
	UserHttpPort string `mapstructure:"USER_HTTP_PORT"`
	UserGrpcPort string `mapstructure:"USER_GRPC_PORT"`
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

package jwtutil

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/matthewhartstonge/argon2"
	"github.com/vasapolrittideah/money-tracker-api/shared/model/apperror"
	"google.golang.org/grpc/codes"
)

func GenerateJwt(ttl time.Duration, secretKey string, userId uuid.UUID) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"sub": userId.String(),
		"exp": now.Add(ttl).Unix(),
		"iat": now.Unix(),
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secretKey))
	if err != nil {
		return "", apperror.New(codes.Internal, fmt.Errorf("unable to sign token: %v", err.Error()))
	}

	return token, nil
}

func ValidateJwt(token string, secretKey string) (*jwt.Token, error) {
	parsed, err := jwt.Parse(token, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, apperror.New(
				codes.Internal,
				fmt.Errorf("unexpected signing method: %v", t.Header["alg"]),
			)
		}

		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, apperror.New(codes.Internal, fmt.Errorf("unable to parse token: %v", err.Error()))
	}

	return parsed, nil
}

func ParseToken(tokenString string, secretKey string) (*jwt.MapClaims, error) {
	token, err := ValidateJwt(tokenString, secretKey)
	if err != nil {
		return nil, apperror.New(codes.Internal, fmt.Errorf("token is invalid or has been expired"))
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, apperror.New(codes.Unauthenticated, fmt.Errorf("token is invalid"))
	}

	return &claims, nil
}

func HashRefreshToken(refreshToken string) (string, error) {
	argon := argon2.DefaultConfig()

	encoded, err := argon.HashEncoded([]byte(refreshToken))
	if err != nil {
		return "", apperror.New(codes.Internal, fmt.Errorf("unable to hash refresh token: %v", err.Error()))
	}

	return string(encoded), nil
}

func VerifyRefreshToken(encoded string, refreshToken string) (bool, error) {
	return argon2.VerifyEncoded([]byte(refreshToken), []byte(encoded))
}

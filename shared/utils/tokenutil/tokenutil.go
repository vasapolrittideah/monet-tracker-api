package tokenutil

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GenerateToken(ttl time.Duration, secretKey string, userId uuid.UUID) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"sub": userId.String(),
		"exp": now.Add(ttl).Unix(),
		"iat": now.Unix(),
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secretKey))
	if err != nil {
		return "", status.Errorf(codes.Internal, "unable to generate token: %v", err)
	}

	return token, nil
}

func ValidateToken(token string, secretKey string) (*jwt.Token, error) {
	parsed, err := jwt.Parse(token, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, status.Errorf(codes.Unauthenticated, "unexpected signing method: %v", t.Header["alg"])
		}

		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "unable to parse token: %v", err.Error())
	}

	return parsed, nil
}

func ParseToken(tokenString string, secretKey string) (*jwt.MapClaims, error) {
	token, err := ValidateToken(tokenString, secretKey)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "unable to validate token: %v", err.Error())
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token")
	}

	return &claims, nil
}

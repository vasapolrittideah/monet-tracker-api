package tokenutil

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func GenerateToken(ttl time.Duration, secretKey string, userID uint64) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": now.Add(ttl).Unix(),
		"iat": now.Unix(),
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secretKey))
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %v", err)
	}

	return token, nil
}

func ValidateToken(token string, secretKey string) (*jwt.Token, error) {
	parsed, err := jwt.Parse(token, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %v", err.Error())
	}

	return parsed, nil
}

func ParseToken(tokenString string, secretKey string) (*jwt.MapClaims, error) {
	token, err := ValidateToken(tokenString, secretKey)
	if err != nil {
		return nil, fmt.Errorf("failed to validate token: %v", err.Error())
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return &claims, nil
}

package passwordutil

import (
	"github.com/matthewhartstonge/argon2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func HashPassword(password string) (string, error) {
	argon := argon2.DefaultConfig()
	encoded, err := argon.HashEncoded([]byte(password))
	if err != nil {
		return "", status.Errorf(codes.Internal, "unable to hash password: %v", err.Error())
	}

	return string(encoded), nil
}

func VerifyPassword(encoded string, password string) (bool, error) {
	return argon2.VerifyEncoded([]byte(password), []byte(encoded))
}

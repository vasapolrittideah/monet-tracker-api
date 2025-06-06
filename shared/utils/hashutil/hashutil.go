package hashutil

import (
	"github.com/matthewhartstonge/argon2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Hash(input string) (string, error) {
	argon := argon2.DefaultConfig()
	encoded, err := argon.HashEncoded([]byte(input))
	if err != nil {
		return "", status.Errorf(codes.Internal, "unable to hash input: %v", err.Error())
	}

	return string(encoded), nil
}

func Verify(encoded string, password string) (bool, error) {
	return argon2.VerifyEncoded([]byte(password), []byte(encoded))
}

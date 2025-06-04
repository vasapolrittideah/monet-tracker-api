package passwordutil

import (
	"fmt"

	"github.com/matthewhartstonge/argon2"
	"github.com/vasapolrittideah/money-tracker-api/shared/constants/errorcode"
	"github.com/vasapolrittideah/money-tracker-api/shared/model/apperror"
)

func HashPassword(password string) (string, *apperror.Error) {
	argon := argon2.DefaultConfig()
	encoded, err := argon.HashEncoded([]byte(password))
	if err != nil {
		return "", apperror.New(errorcode.Internal, fmt.Errorf("unable to hash password: %v", err.Error()))
	}

	return string(encoded), nil
}

func VerifyPassword(encoded string, password string) (bool, error) {
	return argon2.VerifyEncoded([]byte(password), []byte(encoded))
}

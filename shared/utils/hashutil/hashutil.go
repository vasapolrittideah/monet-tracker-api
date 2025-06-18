package hashutil

import (
	"fmt"

	"github.com/matthewhartstonge/argon2"
)

func Hash(input string) (string, error) {
	argon := argon2.DefaultConfig()
	encoded, err := argon.HashEncoded([]byte(input))
	if err != nil {
		return "", fmt.Errorf("failed to hash input: %v", err.Error())
	}

	return string(encoded), nil
}

func Verify(password string, encoded string) (bool, error) {
	return argon2.VerifyEncoded([]byte(password), []byte(encoded))
}

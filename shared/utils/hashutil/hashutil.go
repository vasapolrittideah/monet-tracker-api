package hashutil

import (
	"github.com/matthewhartstonge/argon2"
)

func Hash(input string) (string, error) {
	argon := argon2.DefaultConfig()
	encoded, err := argon.HashEncoded([]byte(input))
	if err != nil {
		return "", err
	}

	return string(encoded), nil
}

func Verify(password string, encoded string) (bool, error) {
	return argon2.VerifyEncoded([]byte(password), []byte(encoded))
}

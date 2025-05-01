package auth

import (
	"crypto/rand"
	"fmt"
)

func GenerateRefreshToken() (string, error) {
	buf := make([]byte, 64)

	_, err := rand.Read(buf)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", buf), nil
}

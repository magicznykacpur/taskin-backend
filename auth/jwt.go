package auth

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWTToken(userID string, secret string, expiresIn time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "taskin-backend",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
		Subject:   userID,
	})
	return token.SignedString([]byte(secret))
}

func ValidateJWTToken(tokenString, secret string) (string, error) {
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		},
	)

	if err != nil {
		return "", err
	}

	userID, err := claims.GetSubject()
	if err != nil {
		return "", err
	}

	return userID, nil
}

func GetAuthTokensFromHeaders(header http.Header) (string, string, error) {
	refreshToken := header.Get("RefreshToken")
	bearerToken := header.Get("Authorization")
	if bearerToken == "" || refreshToken == "" {
		return "", "", fmt.Errorf("missing authorization")
	}

	parts := strings.Split(bearerToken, " ")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("token malformed")
	}

	return refreshToken, parts[1], nil
}

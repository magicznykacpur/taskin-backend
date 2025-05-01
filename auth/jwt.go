package auth

import (
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

package auth

import (
	"testing"
	"time"

	"github.com/magicznykacpur/taskin-backend/auth"
	"github.com/stretchr/testify/assert"
)

func TestJWTToken(t *testing.T) {
	userID := "12345"
	secret := "super-secret"

	tokenString, err := auth.GenerateJWTToken(userID, secret, time.Millisecond * 3)

	assert.NoError(t, err)
	assert.NotEqual(t, "", tokenString)
}

func TestValidateJWTToken(t *testing.T) {
	userID := "12345"
	secret := "super-secret"

	tokenString, err := auth.GenerateJWTToken(userID, secret, time.Second * 20)
	
	assert.NoError(t, err)
	assert.NotEqual(t, "", tokenString)

	validatedUserId, err := auth.ValidateJWTToken(tokenString, secret)

	assert.NoError(t, err)
	assert.Equal(t, userID, validatedUserId)
}

func TestExpiredJWTToken(t *testing.T) {
	userID := "12345"
	secret := "super-secret"

	tokenString, err := auth.GenerateJWTToken(userID, secret, time.Microsecond * 1)
	
	assert.NoError(t, err)
	assert.NotEqual(t, "", tokenString)

	_, err = auth.ValidateJWTToken(tokenString, secret)

	assert.Error(t, err)
}
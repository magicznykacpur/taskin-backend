package auth

import (
	"testing"

	"github.com/magicznykacpur/taskin-backend/auth"
	"github.com/stretchr/testify/assert"
)

func TestHashPassword(t *testing.T) {
	pass := "secure_password"
	hash, err := auth.HashPassword(pass)

	assert.NoError(t, err)

	err = auth.ComparePassword(string(hash), pass)

	assert.NoError(t, err)
}
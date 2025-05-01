package auth

import (
	"testing"

	"github.com/magicznykacpur/taskin-backend/auth"
	"github.com/stretchr/testify/assert"
)

func TestRefreshToken(t *testing.T) {
	token, err := auth.GenerateRefreshToken()

	assert.NoError(t, err)
	// 64 bytes converted to hexadecimal == 128
	assert.Equal(t, 128, len(token))
}
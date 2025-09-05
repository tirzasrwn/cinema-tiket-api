package utils

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGenerateAndVerifyToken(t *testing.T) {
	// Set test secret
	os.Setenv("JWT_SECRET", "test_secret_key")
	os.Setenv("JWT_EXPIRATION_HOURS", "24")

	userID := 1
	email := "test@example.com"

	// Test token generation
	token, err := GenerateToken(userID, email)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// Test token verification
	claims, err := VerifyToken(token)
	assert.NoError(t, err)
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, email, claims.Email)
	assert.WithinDuration(t, time.Now().Add(24*time.Hour), claims.ExpiresAt.Time, time.Second)
}

func TestVerifyToken_Invalid(t *testing.T) {
	os.Setenv("JWT_SECRET", "test_secret_key")

	// Test invalid token
	_, err := VerifyToken("invalid.token.here")
	assert.Error(t, err)

	// Test expired token
	expiredToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJlbWFpbCI6InRlc3RAZXhhbXBsZS5jb20iLCJleHAiOjE2MzAwMDAwMDB9.invalid_signature"
	_, err = VerifyToken(expiredToken)
	assert.Error(t, err)
}

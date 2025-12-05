package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashPassword(t *testing.T) {
	password := "rahasia123"

	hash, err := HashPassword(password)

	assert.NoError(t, err, "HashPassword should not return an error")
	assert.NotEmpty(t, hash, "Hash should not be empty")
	assert.NotEqual(t, password, hash, "Hash should not be equal to the original password")
}

func TestCheckPasswordHash(t *testing.T) {
	password := "rahasia123"
	hash, _ := HashPassword(password)

	match := CheckPasswordHash(password, hash)
	assert.True(t, match, "CheckPasswordHash should return true for matching password and hash")

	matchSalah := CheckPasswordHash("wrongpassword", hash)
	assert.False(t, matchSalah, "CheckPasswordHash should return false for non-matching password and hash")
}

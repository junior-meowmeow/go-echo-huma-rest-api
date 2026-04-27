package auth_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/utility/auth"
)

func TestHashUtility(t *testing.T) {
	t.Run("Should success", func(t *testing.T) {
		password := "my-secure-password-123"

		hash, err := auth.HashPassword(password)
		require.NoError(t, err)
		assert.NotEmpty(t, hash)

		assert.True(t, auth.CheckPasswordHash(password, hash), "Password should match its hash")
	})

	t.Run("Should fails when password is wrong", func(t *testing.T) {
		password := "correct-password"
		wrongPassword := "incorrect-password"

		hash, err := auth.HashPassword(password)
		require.NoError(t, err)

		assert.False(t, auth.CheckPasswordHash(wrongPassword, hash), "Wrong password should not match")
	})

	t.Run("Should success with long password (bcrypt limit)", func(t *testing.T) {
		// This utility use SHA-256 pre-hashing to bypass bcrypt's 72 bytes limit.
		longPassword := strings.Repeat("a", 100)

		hash, err := auth.HashPassword(longPassword)
		require.NoError(t, err, "Should not return bcrypt length error")

		assert.True(t, auth.CheckPasswordHash(longPassword, hash), "100-char password should match")

		modifiedLongPassword := longPassword + "b"
		assert.False(t, auth.CheckPasswordHash(modifiedLongPassword, hash))
	})

	t.Run("Should success with empty password", func(t *testing.T) {
		password := ""
		hash, err := auth.HashPassword(password)
		require.NoError(t, err)
		assert.True(t, auth.CheckPasswordHash(password, hash))
	})

	t.Run("Should generate different hashes", func(t *testing.T) {
		password := "constant-password"

		hash1, _ := auth.HashPassword(password)
		hash2, _ := auth.HashPassword(password)

		assert.NotEqual(t, hash1, hash2, "Bcrypt should apply unique salts to every hash")
		assert.True(t, auth.CheckPasswordHash(password, hash1))
		assert.True(t, auth.CheckPasswordHash(password, hash2))
	})
}

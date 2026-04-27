package auth_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/utility/auth"
)

func TestJWTUtility(t *testing.T) {
	secret := "test-secret-key"
	expiration := 15 * time.Minute
	jwtUtil := auth.NewJWTUtility(secret, expiration)

	t.Run("Generate and Parse Token", func(t *testing.T) {
		userID := "user-uuid-123"
		role := "admin"

		// Generate token
		token, err := jwtUtil.GenerateToken(userID, role)
		require.NoError(t, err)
		assert.NotEmpty(t, token)

		// Parse token
		authCtx, err := jwtUtil.ParseToken(token)
		require.NoError(t, err)

		// Verify content
		assert.Equal(t, userID, authCtx.UserID)
		assert.Equal(t, role, authCtx.Role)
	})

	t.Run("Expired Token", func(t *testing.T) {
		expiredUtil := auth.NewJWTUtility(secret, -1*time.Hour)

		token, err := expiredUtil.GenerateToken("user-1", "user")
		require.NoError(t, err)

		_, err = jwtUtil.ParseToken(token)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "token is expired")
	})

	t.Run("Invalid Signature", func(t *testing.T) {
		userID := "user-1"
		token, err := jwtUtil.GenerateToken(userID, "user")
		require.NoError(t, err)

		wrongUtil := auth.NewJWTUtility("different-secret", expiration)

		_, err = wrongUtil.ParseToken(token)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "signature is invalid")
	})
}

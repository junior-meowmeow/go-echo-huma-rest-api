package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/domain/entity"
)

type jwtUtility struct {
	secretKey       string
	tokenExpiration time.Duration
}

//revive:disable:unexported-return // Intentionally returns an unexported struct to enforce dependency on the interface in other layers.
func NewJWTUtility(secretKey string, tokenExpiration time.Duration) *jwtUtility {
	return &jwtUtility{
		secretKey:       secretKey,
		tokenExpiration: tokenExpiration,
	}
}

//revive:enable:unexported-return

type CustomClaims struct {
	jwt.RegisteredClaims
	Role string `json:"role"`
}

func (u *jwtUtility) GenerateToken(userID string, role string) (string, error) {
	claims := CustomClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(u.tokenExpiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		Role: role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(u.secretKey))
}

func (u *jwtUtility) ParseToken(tokenString string) (AuthContext, error) {
	var authContext AuthContext

	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(t *jwt.Token) (any, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method: %w", entity.ErrInvalidCredentials)
		}
		return []byte(u.secretKey), nil
	})

	if err != nil {
		return authContext, err
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok || !token.Valid {
		return authContext, jwt.ErrTokenInvalidClaims
	}

	authContext = AuthContext{
		UserID: claims.Subject,
		Role:   claims.Role,
	}

	return authContext, nil
}

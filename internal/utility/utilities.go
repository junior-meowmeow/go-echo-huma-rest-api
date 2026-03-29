package utility

import (
	"time"

	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/utility/auth"
)

type Utilities struct {
	Token auth.TokenUtility
}

func NewUtilities(jwtSecret string, tokenExpiration time.Duration) *Utilities {
	return &Utilities{
		Token: auth.NewJWTUtility(jwtSecret, tokenExpiration),
	}
}

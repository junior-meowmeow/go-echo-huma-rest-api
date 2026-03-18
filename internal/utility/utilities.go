package utility

import (
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/utility/auth"
)

type Utilities struct {
	Token auth.TokenUtility
}

func NewUtilities(jwtSecret string) *Utilities {
	return &Utilities{
		Token: auth.NewJWTUtility(jwtSecret),
	}
}

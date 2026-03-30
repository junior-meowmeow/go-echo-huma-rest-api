package auth

type contextKey string

const AuthContextKey contextKey = "auth"

type AuthContext struct {
	UserID string
	Role   string
}

type TokenUtility interface {
	GenerateToken(userID string, role string) (string, error)
	ParseToken(token string) (AuthContext, error)
}

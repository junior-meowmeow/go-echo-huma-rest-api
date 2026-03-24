package schema

import (
	"time"
)

type UsernameSchema struct {
	Username string `json:"username" minLength:"3" maxLength:"50" doc:"User name" example:"new_user"`
}

type PasswordSchema struct {
	Password string `json:"password" minLength:"8" maxLength:"100" format:"password" doc:"User password" example:"pass_word" writeOnly:"true"`
}

type RoleSchema struct {
	Role string `json:"role,omitempty" enum:"user,admin" doc:"User role." example:"user"`
}

type User struct {
	ID string `json:"id" doc:"User ID" readOnly:"true"`
	UsernameSchema
	RoleSchema
	CreatedAt time.Time `json:"createdAt" doc:"Timestamp when the user was created" readOnly:"true"`
}

type RegisterUserRequest struct {
	Body struct {
		UsernameSchema
		PasswordSchema
		RoleSchema
	}
}

type RegisterUserResponse struct {
	Body struct {
		ID string `json:"id" doc:"Created User ID"`
	}
}

type LoginUserRequest struct {
	Body struct {
		UsernameSchema
		PasswordSchema
	}
}

type LoginUserResponse struct {
	Body struct {
		Token string `json:"token" doc:"JWT Bearer token for authentication"`
	}
}

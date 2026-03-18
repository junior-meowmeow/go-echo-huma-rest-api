package entity

import (
	"time"
)

type User struct {
	ID string

	Username string
	Password string
	Role     string

	CreatedAt time.Time
	UpdatedAt time.Time
}

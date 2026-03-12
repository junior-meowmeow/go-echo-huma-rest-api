package entity

import (
	"time"
)

type Review struct {
	ID string

	Author  string
	Rating  int
	Message string

	CreatedAt time.Time
}

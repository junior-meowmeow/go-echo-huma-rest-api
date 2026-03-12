package entity

import (
	"time"
)

type Book struct {
	ID string

	Name             string
	Description      string
	Metadata         BookMetadata
	CoverImageFileID string

	CreatedAt  time.Time
	ModifiedAt time.Time
}

type BookMetadata struct {
	Author string
	ISBN   string
	Genre  string
}

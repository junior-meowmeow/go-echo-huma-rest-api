package entities

import (
	"time"
)

type BookMetadata struct {
	Author string
	ISBN   string
	Genre  string
}

type Book struct {
	ID string

	Name             string
	Description      string
	Metadata         BookMetadata
	CoverImageFileID string

	CreatedAt  time.Time
	ModifiedAt time.Time
}

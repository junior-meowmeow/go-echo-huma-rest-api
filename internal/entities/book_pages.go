package entities

import (
	"time"
)

type BookPageMetadata struct {
	IsBookmarked bool
	Highlight    string
}

type BookPage struct {
	ID     string
	BookID string

	PageNumber          int64
	Content             string
	Metadata            BookPageMetadata
	AttachedImageFileID string

	CreatedAt  time.Time
	ModifiedAt time.Time
}

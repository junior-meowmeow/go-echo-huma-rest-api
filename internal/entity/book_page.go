package entity

import (
	"time"
)

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

type BookPageMetadata struct {
	IsBookmarked bool
	Highlight    string
}

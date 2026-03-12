package entities

import (
	"time"
)

type FileMetadata struct {
	ID string

	Filename    string
	Size        int64
	ContentType string
	S3Key       string

	CreatedAt  time.Time
	ModifiedAt time.Time
}

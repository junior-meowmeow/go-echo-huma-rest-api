package entity

import (
	"time"
)

type FileRecord struct {
	ID string

	FileName    string
	Size        int64
	ContentType string
	S3Key       string

	CreatedAt  time.Time
	ModifiedAt time.Time
}

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

	CreatedAt time.Time
	UpdatedAt time.Time
}

type FileDownloadInfo struct {
	FileName       string
	DownloadURL    string
	ExpirationTime time.Time
}

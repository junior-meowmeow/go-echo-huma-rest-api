package storage

import (
	"context"
	"io"
	"time"
)

type FileStorage interface {
	UploadFile(ctx context.Context, key string, file io.Reader, size int64, contentType string) error
	GetPresignedDownloadURL(ctx context.Context, key string, filename string, duration time.Duration) (string, error)
	CheckFileExists(ctx context.Context, key string) (bool, error)
	ListFiles(ctx context.Context, maxKeys int) ([]string, error)
}

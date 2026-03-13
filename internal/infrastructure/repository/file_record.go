package repository

import (
	"context"

	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/entity"
)

type FileRecordRepository interface {
	CreateFileRecord(ctx context.Context, fileRecord *entity.FileRecord) (string, error)
	GetFileRecordByID(ctx context.Context, fileID string) (entity.FileRecord, error)
}

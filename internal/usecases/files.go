package usecases

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"time"

	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/entities"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/repositories/mongo_repositories"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/repositories/s3_repositories"

	"github.com/google/uuid"
)

type FilesUseCase interface {
	UploadFile(ctx context.Context, fileStream io.Reader, filename string, size int64, contentType string, baseKey string) (string, error)
	GetFileDownloadLink(ctx context.Context, fileID string) (url string, expiresAt time.Time, filename string, err error)
	ListS3Files(ctx context.Context) ([]string, error)
}

type filesUseCase struct {
	FileMetadataRepository mongo_repositories.FileMetadataRepository
	ObjectStorage          s3_repositories.ObjectStorage
}

func NewFilesUseCase(fileMetadataRepo mongo_repositories.FileMetadataRepository, objectStorage s3_repositories.ObjectStorage) *filesUseCase {
	return &filesUseCase{
		FileMetadataRepository: fileMetadataRepo,
		ObjectStorage:          objectStorage,
	}
}

func (u *filesUseCase) UploadFile(ctx context.Context, fileStream io.Reader, filename string, size int64, contentType string, baseKey string) (string, error) {
	ext := filepath.Ext(filename)
	var objectKey string

	// Generate Unique Key
	maxRetries := 5
	for i := range maxRetries {
		objectKey = fmt.Sprintf("%s%s%s", baseKey, uuid.New().String(), ext)

		exists, err := u.ObjectStorage.CheckFileExists(ctx, objectKey)
		if err != nil {
			return "", fmt.Errorf("failed to check file existence in S3: %w", err)
		}

		if !exists {
			break
		}

		if i == maxRetries-1 {
			return "", fmt.Errorf("failed to generate unique S3 key after %d attempts", maxRetries)
		}
	}

	err := u.ObjectStorage.UploadFile(ctx, objectKey, fileStream, size, contentType)
	if err != nil {
		return "", fmt.Errorf("failed to upload to S3: %w", err)
	}

	currentTime := time.Now()

	record := &entities.FileMetadata{
		Filename:    filename,
		Size:        size,
		ContentType: contentType,
		S3Key:       objectKey,
		CreatedAt:   currentTime,
		ModifiedAt:  currentTime,
	}

	id, err := u.FileMetadataRepository.SaveFileMetadata(ctx, record)
	if err != nil {
		return "", fmt.Errorf("failed to save file metadata: %w", err)
	}

	return id, nil
}

func (u *filesUseCase) GetFileDownloadLink(ctx context.Context, fileID string) (string, time.Time, string, error) {
	record, err := u.FileMetadataRepository.GetFileMetadataByID(ctx, fileID)
	if err != nil {
		return "", time.Time{}, "", fmt.Errorf("file not found: %w", err)
	}

	duration := 15 * time.Minute
	expirationTime := time.Now().Add(duration)

	url, err := u.ObjectStorage.GetPresignedDownloadURL(ctx, record.S3Key, record.Filename, duration)
	if err != nil {
		return "", time.Time{}, "", fmt.Errorf("failed to presign url: %w", err)
	}

	return url, expirationTime, record.Filename, nil
}

func (u *filesUseCase) ListS3Files(ctx context.Context) ([]string, error) {
	fileKeys, err := u.ObjectStorage.ListFiles(ctx, 20)
	if err != nil {
		return nil, fmt.Errorf("failed to list S3 files: %w", err)
	}

	return fileKeys, nil
}

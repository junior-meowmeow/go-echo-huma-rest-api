package usecase

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"time"

	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/entity"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/infrastructure/repository"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/infrastructure/storage"

	"github.com/google/uuid"
)

type FileUseCase interface {
	UploadFile(ctx context.Context, fileStream io.Reader, filename string, size int64, contentType string, baseKey string) (string, error)
	GetFileDownloadLink(ctx context.Context, fileID string) (url string, expiresAt time.Time, fileName string, err error)
	GetS3FileList(ctx context.Context) ([]string, error)
}

type fileUseCase struct {
	FileRecordRepository repository.FileRecordRepository
	ObjectStorage        storage.ObjectStorage
}

func NewFileUseCase(fileRecordRepository repository.FileRecordRepository, objectStorage storage.ObjectStorage) *fileUseCase {
	return &fileUseCase{
		FileRecordRepository: fileRecordRepository,
		ObjectStorage:        objectStorage,
	}
}

func (u *fileUseCase) UploadFile(ctx context.Context, fileStream io.Reader, fileName string, size int64, contentType string, baseKey string) (string, error) {
	ext := filepath.Ext(fileName)
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

	fileRecord := &entity.FileRecord{
		FileName:    fileName,
		Size:        size,
		ContentType: contentType,
		S3Key:       objectKey,
		CreatedAt:   currentTime,
		ModifiedAt:  currentTime,
	}

	id, err := u.FileRecordRepository.CreateFileRecord(ctx, fileRecord)
	if err != nil {
		return "", fmt.Errorf("failed to save file record: %w", err)
	}

	return id, nil
}

func (u *fileUseCase) GetFileDownloadLink(ctx context.Context, fileID string) (string, time.Time, string, error) {
	fileRecord, err := u.FileRecordRepository.GetFileRecordByID(ctx, fileID)
	if err != nil {
		return "", time.Time{}, "", fmt.Errorf("file not found: %w", err)
	}

	duration := 15 * time.Minute
	expirationTime := time.Now().Add(duration)

	url, err := u.ObjectStorage.GetPresignedDownloadURL(ctx, fileRecord.S3Key, fileRecord.FileName, duration)
	if err != nil {
		return "", time.Time{}, "", fmt.Errorf("failed to presign url: %w", err)
	}

	return url, expirationTime, fileRecord.FileName, nil
}

func (u *fileUseCase) GetS3FileList(ctx context.Context) ([]string, error) {
	fileKeys, err := u.ObjectStorage.ListFiles(ctx, 20)
	if err != nil {
		return nil, fmt.Errorf("failed to list S3 files: %w", err)
	}

	return fileKeys, nil
}

package usecase

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"time"

	"github.com/google/uuid"

	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/domain/entity"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/domain/port"
)

type FileUseCase interface {
	UploadFile(ctx context.Context, fileStream io.Reader, filename string, size int64, contentType string, baseKey string) (string, error)
	GetFileDownloadLink(ctx context.Context, fileID string) (entity.FileDownloadInfo, error)
	GetS3FileList(ctx context.Context) ([]string, error)
}

type fileUseCase struct {
	FileRecordRepository port.FileRecordRepository
	FileStorage          port.FileStorage
}

//revive:disable:unexported-return // Intentionally returns an unexported struct to enforce dependency on the interface in other layers.
func NewFileUseCase(fileRecordRepository port.FileRecordRepository, fileStorage port.FileStorage) *fileUseCase {
	return &fileUseCase{
		FileRecordRepository: fileRecordRepository,
		FileStorage:          fileStorage,
	}
}

//revive:enable:unexported-return

func (u *fileUseCase) UploadFile(
	ctx context.Context,
	fileStream io.Reader,
	fileName string,
	size int64,
	contentType string,
	baseKey string,
) (string, error) {
	ext := filepath.Ext(fileName)
	var objectKey string

	// Generate Unique Key
	maxRetries := 5
	for i := range maxRetries {
		objectKey = fmt.Sprintf("%s%s%s", baseKey, uuid.New().String(), ext)

		exists, err := u.FileStorage.CheckFileExists(ctx, objectKey)
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

	err := u.FileStorage.UploadFile(ctx, objectKey, fileStream, size, contentType)
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
		UpdatedAt:   currentTime,
	}

	id, err := u.FileRecordRepository.CreateFileRecord(ctx, fileRecord)
	if err != nil {
		return "", fmt.Errorf("failed to save file record: %w", err)
	}

	return id, nil
}

func (u *fileUseCase) GetFileDownloadLink(ctx context.Context, fileID string) (entity.FileDownloadInfo, error) {
	fileRecord, err := u.FileRecordRepository.GetFileRecordByID(ctx, fileID)
	if err != nil {
		return entity.FileDownloadInfo{}, fmt.Errorf("file not found: %w", err)
	}

	const duration = 15 * time.Minute
	expirationTime := time.Now().Add(duration)

	url, err := u.FileStorage.GetPresignedDownloadURL(ctx, fileRecord.S3Key, fileRecord.FileName, duration)
	if err != nil {
		return entity.FileDownloadInfo{}, fmt.Errorf("failed to presign url: %w", err)
	}

	fileDownloadInfo := entity.FileDownloadInfo{
		DownloadURL:    url,
		ExpirationTime: expirationTime,
		FileName:       fileRecord.FileName,
	}

	return fileDownloadInfo, nil
}

func (u *fileUseCase) GetS3FileList(ctx context.Context) ([]string, error) {
	const filesLimit = 20
	fileKeys, err := u.FileStorage.ListFiles(ctx, filesLimit)
	if err != nil {
		return nil, fmt.Errorf("failed to list S3 files: %w", err)
	}

	return fileKeys, nil
}

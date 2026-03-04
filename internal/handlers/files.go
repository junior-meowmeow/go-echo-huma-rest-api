package handlers

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/models"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/repository"

	"github.com/google/uuid"
)

func (h *Handler) UploadFile(ctx context.Context, input *models.UploadFileInput) (*models.UploadFileOutput, error) {
	formData := input.RawBody.Data()
	uploadedFile := formData.File
	ext := filepath.Ext(uploadedFile.Filename)

	var objectKey string

	// Generate Unique Key
	maxRetries := 5
	for i := range maxRetries {
		objectKey = fmt.Sprintf("%s%s%s", formData.ObjectBaseKey, uuid.New().String(), ext)

		exists, err := h.ObjectStorage.CheckFileExists(ctx, objectKey)
		if err != nil {
			return nil, fmt.Errorf("failed to check file existence in S3: %w", err)
		}

		if !exists {
			break
		}

		if i == maxRetries-1 {
			return nil, fmt.Errorf("failed to generate unique S3 key after %d attempts", maxRetries)
		}
	}

	err := h.ObjectStorage.UploadFile(ctx, objectKey, uploadedFile, uploadedFile.Size, uploadedFile.ContentType)
	if err != nil {
		return nil, fmt.Errorf("failed to upload to S3: %w", err)
	}

	currentTime := time.Now()

	record := &repository.FileRecord{
		Filename:    uploadedFile.Filename,
		Size:        uploadedFile.Size,
		ContentType: uploadedFile.ContentType,
		S3Key:       objectKey,
		CreatedAt:   currentTime,
		ModifiedAt:  currentTime,
	}

	id, err := h.Files.SaveFileMetadata(ctx, record)
	if err != nil {
		return nil, fmt.Errorf("failed to save metadata: %w", err)
	}

	resp := &models.UploadFileOutput{
		Body: models.FileMetadata{
			FileID:      id,
			Filename:    uploadedFile.Filename,
			Size:        uploadedFile.Size,
			ContentType: uploadedFile.ContentType,
			CreatedAt:   record.CreatedAt,
		},
	}

	return resp, nil
}

func (h *Handler) GetFileDownloadLink(ctx context.Context, input *models.GetFileDownloadLinkInput) (*models.GetFileDownloadLinkOutput, error) {
	record, err := h.Files.GetFileMetadataByID(ctx, input.FileID)
	if err != nil {
		return nil, fmt.Errorf("file not found: %w", err)
	}

	duration := 15 * time.Minute
	expirationTime := time.Now().Add(duration)

	url, err := h.ObjectStorage.GetPresignedDownloadURL(ctx, record.S3Key, record.Filename, duration)
	if err != nil {
		return nil, fmt.Errorf("failed to sign url: %w", err)
	}

	resp := &models.GetFileDownloadLinkOutput{
		Body: models.DownloadFileBody{
			Filename:    record.Filename,
			DownloadURL: url,
			ExpiresAt:   expirationTime,
		},
	}

	return resp, nil
}

func (h *Handler) ListS3Files(ctx context.Context, _ *struct{}) (*models.ListS3FilesOutput, error) {
	fileKeys, err := h.ObjectStorage.ListFiles(ctx, 20)
	if err != nil {
		return nil, fmt.Errorf("failed to list S3 files: %w", err)
	}

	resp := &models.ListS3FilesOutput{
		Body: models.ListS3FilesBody{
			Files: fileKeys,
			Count: len(fileKeys),
		},
	}

	return resp, nil
}

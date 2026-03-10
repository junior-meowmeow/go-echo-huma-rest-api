package handlers

import (
	"context"

	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/controllers/restapi/models"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/usecases"
)

type FilesHandler interface {
	UploadFile(ctx context.Context, input *models.UploadFileInput) (*models.UploadFileOutput, error)
	GetFileDownloadLink(ctx context.Context, input *models.GetFileDownloadLinkInput) (*models.GetFileDownloadLinkOutput, error)
	ListS3Files(ctx context.Context, input *struct{}) (*models.ListS3FilesOutput, error)
}

type filesHandler struct {
	FilesUseCase usecases.FilesUseCase
}

func NewFilesHandler(filesUseCase usecases.FilesUseCase) *filesHandler {
	return &filesHandler{
		FilesUseCase: filesUseCase,
	}
}

func (h *filesHandler) UploadFile(ctx context.Context, input *models.UploadFileInput) (*models.UploadFileOutput, error) {
	formData := input.RawBody.Data()
	uploadedFile := formData.File

	id, err := h.FilesUseCase.UploadFile(
		ctx,
		uploadedFile,
		uploadedFile.Filename,
		uploadedFile.Size,
		uploadedFile.ContentType,
		formData.ObjectBaseKey,
	)
	if err != nil {
		return nil, err
	}

	resp := &models.UploadFileOutput{
		Body: models.FileMetadata{
			FileID: id,
		},
	}

	return resp, nil
}

func (h *filesHandler) GetFileDownloadLink(ctx context.Context, input *models.GetFileDownloadLinkInput) (*models.GetFileDownloadLinkOutput, error) {
	url, expiresAt, filename, err := h.FilesUseCase.GetFileDownloadLink(ctx, input.FileID)
	if err != nil {
		return nil, err
	}

	resp := &models.GetFileDownloadLinkOutput{
		Body: models.DownloadFileBody{
			Filename:    filename,
			DownloadURL: url,
			ExpiresAt:   expiresAt,
		},
	}

	return resp, nil
}

func (h *filesHandler) ListS3Files(ctx context.Context, _ *struct{}) (*models.ListS3FilesOutput, error) {
	fileKeys, err := h.FilesUseCase.ListS3Files(ctx)
	if err != nil {
		return nil, err
	}

	resp := &models.ListS3FilesOutput{
		Body: models.ListS3FilesBody{
			Files: fileKeys,
			Count: len(fileKeys),
		},
	}

	return resp, nil
}

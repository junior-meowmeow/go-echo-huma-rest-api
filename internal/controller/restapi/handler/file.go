package handler

import (
	"context"

	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/controller/restapi/schema"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/usecase"
)

type FileHandler interface {
	UploadFile(ctx context.Context, input *schema.UploadFileInput) (*schema.UploadFileOutput, error)
	GetFileDownloadLink(ctx context.Context, input *schema.GetFileDownloadLinkInput) (*schema.GetFileDownloadLinkOutput, error)
	ListS3Files(ctx context.Context, input *struct{}) (*schema.ListS3FilesOutput, error)
}

type fileHandler struct {
	FileUseCase usecase.FileUseCase
}

func NewFileHandler(fileUseCase usecase.FileUseCase) *fileHandler {
	return &fileHandler{
		FileUseCase: fileUseCase,
	}
}

func (h *fileHandler) UploadFile(ctx context.Context, input *schema.UploadFileInput) (*schema.UploadFileOutput, error) {
	formData := input.RawBody.Data()
	uploadedFile := formData.File

	id, err := h.FileUseCase.UploadFile(
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

	resp := schema.UploadFileOutput{
		Body: schema.FileRecord{
			FileID: id,
		},
	}

	return &resp, nil
}

func (h *fileHandler) GetFileDownloadLink(ctx context.Context, input *schema.GetFileDownloadLinkInput) (*schema.GetFileDownloadLinkOutput, error) {
	url, expiresAt, filename, err := h.FileUseCase.GetFileDownloadLink(ctx, input.FileID)
	if err != nil {
		return nil, err
	}

	resp := schema.GetFileDownloadLinkOutput{
		Body: schema.DownloadFileBody{
			Filename:    filename,
			DownloadURL: url,
			ExpiresAt:   expiresAt,
		},
	}

	return &resp, nil
}

func (h *fileHandler) ListS3Files(ctx context.Context, _ *struct{}) (*schema.ListS3FilesOutput, error) {
	fileKeys, err := h.FileUseCase.ListS3Files(ctx)
	if err != nil {
		return nil, err
	}

	resp := schema.ListS3FilesOutput{
		Body: schema.ListS3FilesBody{
			Files: fileKeys,
			Count: len(fileKeys),
		},
	}

	return &resp, nil
}

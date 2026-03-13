package handler

import (
	"context"

	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/controller/restapi/schema"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/usecase"
)

type FileHandler interface {
	UploadFile(ctx context.Context, request *schema.UploadFileRequest) (*schema.UploadFileResponse, error)
	GetFileDownloadLink(ctx context.Context, request *schema.GetFileDownloadLinkRequest) (*schema.GetFileDownloadLinkResponse, error)
	GetS3FileList(ctx context.Context, request *schema.GetS3FileListRequest) (*schema.GetS3FileListResponse, error)
}

type fileHandler struct {
	FileUseCase usecase.FileUseCase
}

func NewFileHandler(fileUseCase usecase.FileUseCase) *fileHandler {
	return &fileHandler{
		FileUseCase: fileUseCase,
	}
}

func (h *fileHandler) UploadFile(ctx context.Context, input *schema.UploadFileRequest) (*schema.UploadFileResponse, error) {
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

	resp := schema.UploadFileResponse{}
	resp.Body.ID = id

	return &resp, nil
}

func (h *fileHandler) GetFileDownloadLink(ctx context.Context, request *schema.GetFileDownloadLinkRequest) (*schema.GetFileDownloadLinkResponse, error) {
	url, expiresAt, fileName, err := h.FileUseCase.GetFileDownloadLink(ctx, request.ID)
	if err != nil {
		return nil, err
	}

	resp := schema.GetFileDownloadLinkResponse{
		Body: schema.FileDownloadInfo{
			FileName:    fileName,
			DownloadURL: url,
			ExpiresAt:   expiresAt,
		},
	}

	return &resp, nil
}

func (h *fileHandler) GetS3FileList(ctx context.Context, _ *schema.GetS3FileListRequest) (*schema.GetS3FileListResponse, error) {
	fileKeys, err := h.FileUseCase.GetS3FileList(ctx)
	if err != nil {
		return nil, err
	}

	resp := schema.GetS3FileListResponse{}
	resp.Body.Files = fileKeys
	resp.Body.Count = len(fileKeys)

	return &resp, nil
}

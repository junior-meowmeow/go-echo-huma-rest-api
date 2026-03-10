package models

import (
	"time"

	"github.com/danielgtaylor/huma/v2"
)

type UploadFileInput struct {
	RawBody huma.MultipartFormFiles[struct {
		File          huma.FormFile `form:"file" required:"true" doc:"The file content to upload"`
		ObjectBaseKey string        `form:"objectBaseKey" doc:"Base object key in object storage"`
	}]
}

type FileMetadata struct {
	FileID string `json:"fileid" doc:"File ID"`
}

type UploadFileOutput struct {
	Body FileMetadata `json:"body"`
}

type GetFileDownloadLinkInput struct {
	FileID string `query:"id" example:"123" doc:"file id"`
}

type DownloadFileBody struct {
	Filename    string    `json:"filename"`
	DownloadURL string    `json:"downloadUrl"`
	ExpiresAt   time.Time `json:"expiresAt"`
}

type GetFileDownloadLinkOutput struct {
	Body DownloadFileBody `json:"body"`
}

type ListS3FilesBody struct {
	Files []string `json:"files" doc:"List of file keys found in the S3 bucket"`
	Count int      `json:"count" doc:"Total number of files found"`
}

type ListS3FilesOutput struct {
	Body ListS3FilesBody `json:"body"`
}

package schema

import (
	"time"

	"github.com/danielgtaylor/huma/v2"
)

type FileDownloadInfo struct {
	FileName    string    `json:"fileName" doc:"File name"`
	DownloadURL string    `json:"downloadUrl" doc:"Temporary download URL"`
	ExpiresAt   time.Time `json:"expiresAt" doc:"Timestamp when the download URL expires"`
}

type UploadFileRequest struct {
	RawBody huma.MultipartFormFiles[struct {
		File          huma.FormFile `form:"file" required:"true" doc:"File content to upload"`
		ObjectBaseKey string        `form:"objectBaseKey" doc:"Base object key in object storage"`
	}]
}

type UploadFileResponse struct {
	Body struct {
		ID string `json:"id" doc:"Uploaded File ID"`
	}
}

type GetFileDownloadLinkRequest struct {
	ID string `query:"id" pattern:"^[a-fA-F0-9]{24}$" doc:"File ID"`
}

type GetFileDownloadLinkResponse struct {
	Body FileDownloadInfo
}

type GetS3FileListRequest struct{}

type GetS3FileListResponse struct {
	Body struct {
		Files []string `json:"files" doc:"List of file keys found in the S3 bucket"`
		Count int      `json:"count" doc:"Total number of files found"`
	}
}

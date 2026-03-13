package storage

import (
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/infrastructure/storage/s3api"
)

type Storages struct {
	ObjectStorage ObjectStorage
}

func NewStorages(s3Client *s3.Client, bucketName string) *Storages {
	return &Storages{
		ObjectStorage: s3api.NewS3Storage(s3Client, bucketName),
	}
}

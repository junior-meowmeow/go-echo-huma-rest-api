package s3_repositories

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type ObjectStorage interface {
	UploadFile(ctx context.Context, key string, file multipart.File, size int64, contentType string) error
	GetPresignedDownloadURL(ctx context.Context, key string, filename string, duration time.Duration) (string, error)
	CheckFileExists(ctx context.Context, key string) (bool, error)
	ListFiles(ctx context.Context, maxKeys int) ([]string, error)
}

type S3Repository struct {
	Client        *s3.Client
	PresignClient *s3.PresignClient
	BucketName    string
}

func NewS3Repository(client *s3.Client, bucketName string) *S3Repository {
	presignClient := s3.NewPresignClient(client)
	return &S3Repository{
		Client:        client,
		PresignClient: presignClient,
		BucketName:    bucketName,
	}
}

func (r *S3Repository) UploadFile(ctx context.Context, key string, file multipart.File, size int64, contentType string) error {
	_, err := r.Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:        aws.String(r.BucketName),
		Key:           aws.String(key),
		Body:          file,
		ContentLength: aws.Int64(size),
		ContentType:   aws.String(contentType),
	})
	return err
}

func (r *S3Repository) GetPresignedDownloadURL(ctx context.Context, key string, filename string, duration time.Duration) (string, error) {
	req, err := r.PresignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket:                     aws.String(r.BucketName),
		Key:                        aws.String(key),
		ResponseContentDisposition: aws.String(fmt.Sprintf("attachment; filename=\"%s\"", filename)),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = duration
	})
	if err != nil {
		return "", err
	}
	return req.URL, nil
}

func (r *S3Repository) CheckFileExists(ctx context.Context, key string) (bool, error) {
	_, err := r.Client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(r.BucketName),
		Key:    aws.String(key),
	})

	if err != nil {
		// Check if the error is specifically "Not Found"
		var notFound *types.NotFound
		if errors.As(err, &notFound) {
			return false, nil
		}

		// Return the actual error
		return false, fmt.Errorf("failed to check file existence: %w", err)
	}

	// If no error, the file exists
	return true, nil
}

func (r *S3Repository) ListFiles(ctx context.Context, maxKeys int) ([]string, error) {
	output, err := r.Client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket:  aws.String(r.BucketName),
		MaxKeys: aws.Int32(int32(maxKeys)),
	})
	if err != nil {
		return nil, err
	}

	var keys []string
	for _, item := range output.Contents {
		keys = append(keys, *item.Key)
	}
	return keys, nil
}

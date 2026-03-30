package app

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/aws"
	sdkConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/config"
)

func newS3Client(ctx context.Context, cfg config.S3Config) (*s3.Client, error) {
	sdkConfig, err := sdkConfig.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config: %w", err)
	}

	client := s3.NewFromConfig(sdkConfig, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(cfg.Endpoint)
		o.UsePathStyle = true
	})

	slog.InfoContext(ctx, fmt.Sprintf("Created a new S3 clients and connected to %s", cfg.Endpoint))

	return client, nil
}

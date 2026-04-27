package s3api_test

import (
	"context"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	smithyendpoints "github.com/aws/smithy-go/endpoints"
	"github.com/docker/go-connections/nat"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	accesskey = "a"
	secretkey = "b"
	token     = "c"
	region    = "us-east-1"
)

type resolverV2 struct {
	// you could inject additional application context here as well
}

func (*resolverV2) ResolveEndpoint(ctx context.Context, params s3.EndpointParameters) (
	smithyendpoints.Endpoint, error,
) {
	// delegate back to the default v2 resolver otherwise
	return s3.NewDefaultEndpointResolverV2().ResolveEndpoint(ctx, params)
}

func s3Client(ctx context.Context, c testcontainers.Container) (*s3.Client, error) {
	mappedPort, err := c.MappedPort(ctx, nat.Port("4566/tcp"))
	if err != nil {
		return nil, err
	}

	provider, err := testcontainers.NewDockerProvider()
	if err != nil {
		return nil, err
	}
	defer provider.Close()

	host, err := provider.DaemonHost(ctx)
	if err != nil {
		return nil, err
	}

	awsCfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accesskey, secretkey, token)),
	)
	if err != nil {
		return nil, err
	}

	// reference: https://aws.github.io/aws-sdk-go-v2/docs/configuring-sdk/endpoints/#with-both
	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String("http://" + host + ":" + mappedPort.Port())
		o.EndpointResolverV2 = &resolverV2{}
		o.UsePathStyle = true
	})

	return client, nil
}

func setupS3Client(t *testing.T) *s3.Client {
	t.Helper()

	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "hectorvent/floci:latest",
		ExposedPorts: []string{"4566/tcp"},
		WaitingFor:   wait.ForLog("Ready").WithStartupTimeout(30 * time.Second),
	}

	ctr, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	require.NoError(t, err, "failed to start container")

	t.Cleanup(func() {
		if err := ctr.Terminate(ctx); err != nil {
			t.Logf("failed to terminate container: %s", err)
		}
	})

	s3Client, err := s3Client(ctx, ctr)
	require.NoError(t, err, "failed to create S3 Client")

	return s3Client
}

//revive:disable-next-line:context-as-argument // Let testing.T be the first argument.
func cleanBucket(t *testing.T, ctx context.Context, client *s3.Client, bucket string) {
	t.Helper()

	out, err := client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
	})
	require.NoError(t, err)

	if len(out.Contents) == 0 {
		return
	}

	for _, obj := range out.Contents {
		_, err := client.DeleteObject(ctx, &s3.DeleteObjectInput{
			Bucket: aws.String(bucket),
			Key:    obj.Key,
		})
		require.NoError(t, err)
	}
}

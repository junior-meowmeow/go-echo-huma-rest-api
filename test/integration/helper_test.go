package integration_test

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	smithyendpoints "github.com/aws/smithy-go/endpoints"
	"github.com/docker/go-connections/nat"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/localstack"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func setupMongoDB(t *testing.T) *mongo.Database {
	t.Helper()

	ctx := context.Background()

	mongodbContainer, err := mongodb.Run(ctx, "mongo:8.0")
	require.NoError(t, err, "failed to start container")

	t.Cleanup(func() {
		if err := mongodbContainer.Terminate(ctx); err != nil {
			t.Logf("failed to terminate container: %s", err)
		}
	})

	endpoint, err := mongodbContainer.ConnectionString(ctx)
	require.NoError(t, err, "failed to get connection string")

	client, err := mongo.Connect(options.Client().ApplyURI(endpoint))
	require.NoError(t, err, "failed to connect to MongoDB")

	err = client.Ping(ctx, nil)
	require.NoError(t, err, "failed to ping MongoDB")

	// Return a random DB for isolation
	return client.Database("test_db_" + t.Name())
}

func cleanCollection(t *testing.T, coll *mongo.Collection) {
	t.Helper()
	_, err := coll.DeleteMany(context.Background(), bson.D{})
	require.NoError(t, err, "failed to clear collection")
}

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

func s3Client(ctx context.Context, l *localstack.LocalStackContainer) (*s3.Client, error) {
	mappedPort, err := l.MappedPort(ctx, nat.Port("4566/tcp"))
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

	ctr, err := localstack.Run(ctx, "localstack/localstack:s3-latest")
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

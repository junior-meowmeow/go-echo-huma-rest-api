package app

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"github.com/google/uuid"

	"github.com/junior-meowmeow/go-echo-huma-rest-api/internal/config"
)

//nolint:gochecknoglobals // Respect the original code.
var uuidType = reflect.TypeFor[uuid.UUID]()

func newMongoDBClient(ctx context.Context, cfg config.MongoConfig) (*mongo.Client, error) {
	hostPort := net.JoinHostPort(cfg.Host, cfg.Port)
	mongoURI := fmt.Sprintf("mongodb://%s:%s@%s/%s", cfg.DBUser, cfg.DBPass, hostPort, cfg.DBName)

	// Source: https://gist.github.com/SupaHam/3afe982dc75039356723600ccc91ff77
	registry := bson.NewRegistry()
	registry.RegisterTypeEncoder(
		uuidType,
		bson.ValueEncoderFunc(uuidEncoder),
	)
	registry.RegisterTypeDecoder(
		uuidType,
		bson.ValueDecoderFunc(uuidDecoder),
	)

	opts := options.Client().ApplyURI(mongoURI).SetRegistry(registry)

	client, err := mongo.Connect(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to create mongo client: %w", err)
	}

	slog.InfoContext(ctx, fmt.Sprintf("Created a new MongoDB client and connected to %s", hostPort))

	if err := pingMongoDB(ctx, client); err != nil {
		return nil, fmt.Errorf("failed to ping mongoDB: %w", err)
	}

	return client, nil
}

func pingMongoDB(ctx context.Context, client *mongo.Client) error {
	const pingTimeout = 5 * time.Second
	pingCtx, cancel := context.WithTimeout(ctx, pingTimeout)
	defer cancel()

	return client.Ping(pingCtx, nil)
}

func disconnectMongoDB(ctx context.Context, client *mongo.Client) error {
	if client == nil {
		slog.DebugContext(ctx, "MongoDB Client is nil.")
		return nil
	}

	err := client.Disconnect(ctx)
	if err != nil {
		return err
	}
	slog.InfoContext(ctx, "MongoDB Client disconnected.")
	return nil
}

func uuidEncoder(
	_ bson.EncodeContext,
	vw bson.ValueWriter,
	val reflect.Value,
) error {
	if !val.IsValid() || val.Type() != uuidType {
		return bson.ValueEncoderError{
			Name:     "uuidEncoder",
			Types:    []reflect.Type{uuidType},
			Received: val,
		}
	}

	bytes, _ := val.Interface().(uuid.UUID)

	return vw.WriteBinaryWithSubtype(
		bytes[:],
		bson.TypeBinaryUUID,
	)
}

func uuidDecoder(
	_ bson.DecodeContext,
	vr bson.ValueReader,
	val reflect.Value,
) error {
	if !val.IsValid() || !val.CanSet() || val.Type() != uuidType {
		return bson.ValueDecoderError{
			Name:     "uuidDecoder",
			Types:    []reflect.Type{uuidType},
			Received: val,
		}
	}

	//nolint:exhaustive // The default case handles all other unsupported BSON types.
	switch vrType := vr.Type(); vrType {
	case bson.TypeBinary:
		bytes, bytesType, err := vr.ReadBinary()
		if err != nil {
			return err
		}
		if bytesType != bson.TypeBinaryUUID {
			return fmt.Errorf(
				"unsupported binary subtype %v for UUID",
				bytesType,
			)
		}

		//revive:disable-next-line:var-naming // Respect the original code.
		uuid_, err := uuid.FromBytes(bytes)
		if err != nil {
			return err
		}

		val.Set(reflect.ValueOf(uuid_))
		return nil

	case bson.TypeNull:
		return vr.ReadNull()

	case bson.TypeUndefined:
		return vr.ReadUndefined()

	default:
		return fmt.Errorf(
			"cannot decode %v into UUID",
			vrType,
		)
	}
}

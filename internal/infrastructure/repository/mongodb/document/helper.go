package document

import (
	"fmt"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func IDToString(id any) (string, error) {
	switch value := id.(type) {
	case bson.ObjectID:
		return value.Hex(), nil
	case uuid.UUID:
		return value.String(), nil
	case bson.Binary: // UUID
		outID, err := uuid.FromBytes(value.Data)
		if err != nil {
			return "", fmt.Errorf("failed to decode UUID: %w", err)
		}
		return outID.String(), nil
	case string:
		return value, nil
	default:
		return "", fmt.Errorf("unsupported ID type returned from database: %T", id)
	}
}

func StringToObjectID(id string) (bson.ObjectID, error) {
	if id == "" {
		return bson.NilObjectID, nil
	}

	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return bson.NilObjectID, err
	}

	return oid, nil
}

func StringToUUID(id string) (uuid.UUID, error) {
	if id == "" {
		return uuid.Nil, nil
	}

	outID, err := uuid.Parse(id)
	if err != nil {
		return uuid.Nil, err
	}

	return outID, nil
}

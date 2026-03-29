package document

import (
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func IDToString(id any) (string, error) {
	switch value := id.(type) {
	case bson.ObjectID:
		return value.Hex(), nil
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

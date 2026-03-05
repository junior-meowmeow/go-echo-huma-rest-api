package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Review struct {
	ID        bson.ObjectID `bson:"_id,omitempty"`
	Author    string        `bson:"author"`
	Rating    int           `bson:"rating"`
	Message   string        `bson:"message"`
	CreatedAt time.Time     `bson:"createdAt"`
}

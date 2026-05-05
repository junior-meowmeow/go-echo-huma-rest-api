package database

import (
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type MongoAdapter struct {
	db *mongo.Database
}

func NewMongoAdapter(db *mongo.Database) *MongoAdapter {
	return &MongoAdapter{db: db}
}

package models

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

var client *mongo.Client
var ctx context.Context
var footballerCollection *mongo.Collection

// SetMongoConnection for initializing a database connection.
func SetMongoConnection() {
	client = ConnectMongo()
	ctx = context.TODO()
	footballerCollection = client.Database("footballer").Collection("footballer")
}

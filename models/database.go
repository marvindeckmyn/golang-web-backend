package models

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ConnectMongo to connect to the mongo database.
func ConnectMongo() *mongo.Client {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017/")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	return client
}

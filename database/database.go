package database

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

/*
    Simple setup to connect to the database
*/
func GetMongoDbCollection() (*mongo.Client, error) {
    // Connects to the database
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}

    // Pings the database to check if everything
    // is alright.
	err = client.Ping(context.Background(), readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}

	return client, nil
}

/*
    Extended setup for databases and collections
*/
func getMongoDbCollection(DbName string, CollectionName string) (*mongo.Collection, error) {
    // Calls the database connection function
	client, err := GetMongoDbCollection()
	if err != nil {
		return nil, err
	}

    // Establishes a simpler method to work with collections
	collection := client.Database(DbName).Collection(CollectionName)

	return collection, nil
}

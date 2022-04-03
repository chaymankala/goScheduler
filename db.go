package main

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var dbClient *mongo.Client

func initDB() {
	DbURI := getEnvVariable(ENV_DB_URI)
	client, err := mongo.NewClient(options.Client().ApplyURI(DbURI))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), dbConnectTimeout)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to Database!")
	dbClient = client
}

func getCollection(collectionName string) *mongo.Collection {
	dbName := getEnvVariable(ENV_DB_NAME)
	collectionJobs := dbClient.Database(dbName).Collection(collectionName)
	return collectionJobs
}

func getJobsCollectionName() string {
	return getEnvVariable(ENV_COLLECTION_NAME)
}

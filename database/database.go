package database

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// connect establishes a connection to MongoDB and returns the client and collections
func Connect() (*mongo.Client, *mongo.Collection, *mongo.Collection) {

	// connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// MongoDB Connection String
	connectionString := os.Getenv("MONGO_URI")

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionString))
	if err != nil {
		log.Fatal(err)
	}

	// Initialize Collections
	db := client.Database("GoGoNotes")
	userCollection := db.Collection("users")
	noteCollection := db.Collection("notes")

	// Check connection by Running a Query
	err = userCollection.FindOne(ctx, bson.M{}).Err()
	if err != nil && err != mongo.ErrNoDocuments {
		log.Fatalf("Failed to query user collection: %v", err)
	}

	log.Println("Successfully Connected to the MongoDB database ;)")
	return client, userCollection, noteCollection

}

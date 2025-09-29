// +build ignore

package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Connect to MongoDB
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017" // Default URI
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}
	defer client.Disconnect(context.TODO())

	// Get database and collection
	db := client.Database("fms_audit") // Adjust database name if different
	collection := db.Collection("orders")

	// Update all orders that don't have table_number field
	filter := bson.M{"table_number": bson.M{"$exists": false}}
	update := bson.M{"$set": bson.M{"table_number": 1}} // Set default table number to 1

	result, err := collection.UpdateMany(context.TODO(), filter, update)
	if err != nil {
		log.Fatal("Failed to update orders:", err)
	}

	fmt.Printf("Successfully updated %d orders with default table_number = 1\n", result.ModifiedCount)

	// Verify the update
	count, err := collection.CountDocuments(context.TODO(), bson.M{"table_number": bson.M{"$exists": true}})
	if err != nil {
		log.Fatal("Failed to count updated orders:", err)
	}

	fmt.Printf("Total orders with table_number field: %d\n", count)
}

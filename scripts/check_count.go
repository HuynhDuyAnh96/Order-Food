//go:build ignore

package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	uri := os.Getenv("MONGODB_URI")
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}
	count, err := client.Database("fms_audit").Collection("dishes").CountDocuments(context.Background(), map[string]interface{}{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Số documents trong dishes trên Atlas: %d\n", count)
}

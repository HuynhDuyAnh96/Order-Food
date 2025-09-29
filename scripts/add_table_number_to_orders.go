package main

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Kết nối MongoDB
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(context.Background())

	// Kiểm tra kết nối
	if err := client.Ping(context.Background(), nil); err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}

	db := client.Database("fms_audit")
	collection := db.Collection("orders")

	// Tìm tất cả orders không có table_number hoặc có table_number = 0
	filter := bson.M{
		"$or": []bson.M{
			{"table_number": bson.M{"$exists": false}}, // Không có field table_number
			{"table_number": 0},                        // Hoặc table_number = 0
		},
	}

	// Đếm số orders cần update
	count, err := collection.CountDocuments(context.Background(), filter)
	if err != nil {
		log.Fatalf("Failed to count documents: %v", err)
	}

	fmt.Printf("Found %d orders without table_number\n", count)

	if count == 0 {
		fmt.Println("No orders need to be updated")
		return
	}

	// Update tất cả orders không có table_number, set default = 1
	update := bson.M{
		"$set": bson.M{
			"table_number": 1, // Default table number
		},
	}

	result, err := collection.UpdateMany(context.Background(), filter, update)
	if err != nil {
		log.Fatalf("Failed to update orders: %v", err)
	}

	fmt.Printf("Successfully updated %d orders with default table_number = 1\n", result.ModifiedCount)

	// Verify the update
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		log.Fatalf("Failed to verify update: %v", err)
	}
	defer cursor.Close(context.Background())

	fmt.Println("\nVerifying updated orders:")
	var orders []bson.M
	if err = cursor.All(context.Background(), &orders); err != nil {
		log.Fatalf("Failed to decode orders: %v", err)
	}

	for _, order := range orders {
		fmt.Printf("Order ID: %s, Table Number: %v\n", order["_id"], order["table_number"])
	}
}

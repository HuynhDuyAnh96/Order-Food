// +build ignore

package main

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Dish struct {
	ID            string  `bson:"_id"`
	Name          string  `bson:"name"`
	Price         float64 `bson:"price"`
	Description   string  `bson:"description"`
	ImageURL      string  `bson:"image_url"`
	CookingMethod string  `bson:"cooking_method"`
	Rating        float64 `bson:"rating"`
	Featured      bool    `bson:"featured"`
}

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
	collection := db.Collection("dishes")

	// Xóa dữ liệu cũ
	collection.Drop(context.Background())

	// Dữ liệu mẫu
	dishes := []interface{}{
		Dish{
			ID:            "1",
			Name:          "Sò điệp nướng phô mai",
			Description:   "Sò điệp tươi nướng với phô mai béo ngậy, thơm lừng",
			Price:         85000,
			CookingMethod: "grilled",
			Rating:        4.8,
			Featured:      true,
			ImageURL:      "http://localhost:8080/images/h6.jpg",
		},
		Dish{
			ID:            "2",
			Name:          "Càng ghẹ rang me",
			Description:   "Càng ghẹ tươi rang với mắm ruốc đậm đà, cay nồng",
			Price:         90000,
			CookingMethod: "stir-fried",
			Rating:        4.7,
			Featured:      true,
			ImageURL:      "http://localhost:8080/images/h7.jpg",
		},
		Dish{
			ID:            "3",
			Name:          "Sò lông nướng mở hành",
			Description:   "Sò lông tươi nướng với mở hành thơm phức",
			Price:         150000,
			CookingMethod: "grilled",
			Rating:        4.9,
			Featured:      true,
			ImageURL:      "http://localhost:8080/images/h2.jpg",
		},
		Dish{
			ID:            "4",
			Name:          "Ốc hương xào bắp",
			Description:   "Ốc hương tươi xào với bắp non giòn ngọt",
			Price:         120000,
			CookingMethod: "stir-fried",
			Rating:        4.5,
			Featured:      false,
			ImageURL:      "http://localhost:8080/images/h1.jpg",
		},
		Dish{
			ID:            "5",
			Name:          "Ốc bươu nướng tiêu",
			Description:   "Ốc bươu tươi nướng với tiêu đen thơm cay",
			Price:         75000,
			CookingMethod: "grilled",
			Rating:        4.6,
			Featured:      true,
			ImageURL:      "http://localhost:8080/images/h4.jpg",
		},
		Dish{
			ID:            "6",
			Name:          "Ốc len xào dừa",
			Description:   "Ốc len tươi xào với nước dừa thơm ngọt",
			Price:         95000,
			CookingMethod: "steamed",
			Rating:        4.4,
			Featured:      false,
			ImageURL:      "http://localhost:8080/images/h5.jpg",
		},
		Dish{
			ID:            "7",
			Name:          "Mực xào mì",
			Description:   "Mực tươi xào với mì tôm đậm đà",
			Price:         110000,
			CookingMethod: "stir-fried",
			Rating:        4.3,
			Featured:      false,
			ImageURL:      "http://localhost:8080/images/h8.jpg",
		},
		Dish{
			ID:            "8",
			Name:          "Hào nướng phô mai",
			Description:   "Hào tươi nướng với phô mai béo ngậy",
			Price:         65000,
			CookingMethod: "grilled",
			Rating:        4.5,
			Featured:      true,
			ImageURL:      "http://localhost:8080/images/h11.jpg",
		},
		Dish{
			ID:            "9",
			Name:          "Nghêu hấp xả",
			Description:   "Nghêu tươi hấp với sả thơm, thanh mát",
			Price:         55000,
			CookingMethod: "steamed",
			Rating:        4.2,
			Featured:      false,
			ImageURL:      "http://localhost:8080/images/h9.jpg",
		},
		Dish{
			ID:            "10",
			Name:          "Tôm hùm nướng bơ tỏi",
			Description:   "Tôm hùm tươi nướng với bơ tỏi thơm phức",
			Price:         250000,
			CookingMethod: "grilled",
			Rating:        4.9,
			Featured:      true,
			ImageURL:      "http://localhost:8080/images/h3.jpg",
		},
	}

	// Insert dữ liệu
	result, err := collection.InsertMany(context.Background(), dishes)
	if err != nil {
		log.Fatalf("Failed to insert dishes: %v", err)
	}

	log.Printf("✅ Successfully inserted %d dishes into MongoDB", len(result.InsertedIDs))
	log.Println("🎉 Seed data completed!")
}

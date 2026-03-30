// +build ignore

package main

import (
	"context"
	"fmt"
	"log"
	"os"

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
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}
	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(context.Background())

	if err := client.Ping(context.Background(), nil); err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}

	db := client.Database("fms_audit")
	collection := db.Collection("dishes")

	collection.Drop(context.Background())

	dishes := []interface{}{
		// ── NGHÊU ──────────────────────────────────────────────────
		Dish{
			ID: "1", Name: "Nghêu Hấp Xả",
			Description: "Nghêu tươi hấp với sả thơm, thanh mát",
			Price: 40000, CookingMethod: "steamed", Rating: 4.3, Featured: false,
			ImageURL: "http://localhost:8080/images/h9.jpg",
		},
		Dish{
			ID: "2", Name: "Nghêu Hấp Thái",
			Description: "Nghêu tươi hấp kiểu Thái, chua cay đậm đà",
			Price: 40000, CookingMethod: "steamed", Rating: 4.4, Featured: false,
			ImageURL: "http://localhost:8080/images/h9.jpg",
		},

		// ── MÓNG TAY ───────────────────────────────────────────────
		Dish{
			ID: "3", Name: "Móng Tay Xào Mì",
			Description: "Móng tay xào với mì sợi đậm đà",
			Price: 40000, CookingMethod: "stir-fried", Rating: 4.2, Featured: false,
			ImageURL: "http://localhost:8080/images/h1.jpg",
		},
		Dish{
			ID: "4", Name: "Móng Tay Xào Rau Muống",
			Description: "Móng tay xào cùng rau muống giòn ngon",
			Price: 40000, CookingMethod: "stir-fried", Rating: 4.2, Featured: false,
			ImageURL: "http://localhost:8080/images/h1.jpg",
		},
		Dish{
			ID: "5", Name: "Móng Tay Bơ Bắp",
			Description: "Móng tay xào bơ bắp ngọt béo",
			Price: 40000, CookingMethod: "stir-fried", Rating: 4.3, Featured: false,
			ImageURL: "http://localhost:8080/images/h1.jpg",
		},
		Dish{
			ID: "6", Name: "Móng Tay Xào Me",
			Description: "Móng tay xào me chua ngọt hấp dẫn",
			Price: 40000, CookingMethod: "stir-fried", Rating: 4.4, Featured: false,
			ImageURL: "http://localhost:8080/images/h1.jpg",
		},

		// ── ỐC LÁC ────────────────────────────────────────────────
		Dish{
			ID: "7", Name: "Ốc Lác Xào Sa Tế",
			Description: "Ốc lác xào sa tế cay thơm đặc trưng",
			Price: 40000, CookingMethod: "stir-fried", Rating: 4.5, Featured: true,
			ImageURL: "http://localhost:8080/images/h4.jpg",
		},
		Dish{
			ID: "8", Name: "Ốc Lác Hấp Xả",
			Description: "Ốc lác hấp sả thơm tươi ngon",
			Price: 40000, CookingMethod: "steamed", Rating: 4.3, Featured: false,
			ImageURL: "http://localhost:8080/images/h4.jpg",
		},
		Dish{
			ID: "9", Name: "Ốc Lác Hấp Thái",
			Description: "Ốc lác hấp kiểu Thái chua cay lạ miệng",
			Price: 40000, CookingMethod: "steamed", Rating: 4.3, Featured: false,
			ImageURL: "http://localhost:8080/images/h4.jpg",
		},

		// ── MỰC ───────────────────────────────────────────────────
		Dish{
			ID: "10", Name: "Mực Xào Mì",
			Description: "Mực tươi xào với mì sợi đậm đà",
			Price: 80000, CookingMethod: "stir-fried", Rating: 4.3, Featured: false,
			ImageURL: "http://localhost:8080/images/h8.jpg",
		},
		Dish{
			ID: "11", Name: "Mực Xào Rau Muống",
			Description: "Mực tươi xào rau muống giòn ngon",
			Price: 80000, CookingMethod: "stir-fried", Rating: 4.3, Featured: false,
			ImageURL: "http://localhost:8080/images/h8.jpg",
		},
		Dish{
			ID: "12", Name: "Mực Bơ Bắp",
			Description: "Mực tươi xào bơ bắp béo ngậy",
			Price: 80000, CookingMethod: "stir-fried", Rating: 4.4, Featured: false,
			ImageURL: "http://localhost:8080/images/h8.jpg",
		},
		Dish{
			ID: "13", Name: "Mực Sa Tế",
			Description: "Mực tươi xào sa tế cay nồng",
			Price: 80000, CookingMethod: "stir-fried", Rating: 4.5, Featured: true,
			ImageURL: "http://localhost:8080/images/h8.jpg",
		},

		// ── ỐC HƯƠNG ──────────────────────────────────────────────
		Dish{
			ID: "14", Name: "Ốc Hương Bơ Bắp",
			Description: "Ốc hương tươi xào bơ bắp ngọt béo",
			Price: 50000, CookingMethod: "stir-fried", Rating: 4.5, Featured: false,
			ImageURL: "http://localhost:8080/images/h1.jpg",
		},
		Dish{
			ID: "15", Name: "Ốc Hương Bơ Tỏi",
			Description: "Ốc hương tươi xào bơ tỏi thơm phức",
			Price: 50000, CookingMethod: "stir-fried", Rating: 4.6, Featured: true,
			ImageURL: "http://localhost:8080/images/h1.jpg",
		},
		Dish{
			ID: "16", Name: "Ốc Hương Cháy Tỏi",
			Description: "Ốc hương tươi rang cháy tỏi vàng giòn",
			Price: 50000, CookingMethod: "stir-fried", Rating: 4.7, Featured: true,
			ImageURL: "http://localhost:8080/images/h1.jpg",
		},
		Dish{
			ID: "17", Name: "Ốc Hương Me",
			Description: "Ốc hương tươi xào me chua ngọt",
			Price: 50000, CookingMethod: "stir-fried", Rating: 4.5, Featured: false,
			ImageURL: "http://localhost:8080/images/h1.jpg",
		},
		Dish{
			ID: "18", Name: "Ốc Hương Rang Muối",
			Description: "Ốc hương tươi rang muối thơm giòn",
			Price: 50000, CookingMethod: "stir-fried", Rating: 4.4, Featured: false,
			ImageURL: "http://localhost:8080/images/h1.jpg",
		},

		// ── SÒ HUYẾT ──────────────────────────────────────────────
		Dish{
			ID: "19", Name: "Sò Huyết Xào Me",
			Description: "Sò huyết tươi xào me chua cay đậm đà",
			Price: 50000, CookingMethod: "stir-fried", Rating: 4.5, Featured: false,
			ImageURL: "http://localhost:8080/images/h5.jpg",
		},
		Dish{
			ID: "20", Name: "Sò Huyết Cháy Tỏi",
			Description: "Sò huyết tươi rang cháy tỏi vàng giòn",
			Price: 50000, CookingMethod: "stir-fried", Rating: 4.6, Featured: true,
			ImageURL: "http://localhost:8080/images/h5.jpg",
		},
		Dish{
			ID: "21", Name: "Sò Huyết Bơ Tỏi",
			Description: "Sò huyết tươi xào bơ tỏi thơm béo",
			Price: 50000, CookingMethod: "stir-fried", Rating: 4.5, Featured: false,
			ImageURL: "http://localhost:8080/images/h5.jpg",
		},

		// ── CÀNG GHẸ ──────────────────────────────────────────────
		Dish{
			ID: "22", Name: "Càng Ghẹ Cháy Tỏi",
			Description: "Càng ghẹ tươi rang cháy tỏi vàng giòn",
			Price: 50000, CookingMethod: "stir-fried", Rating: 4.7, Featured: true,
			ImageURL: "http://localhost:8080/images/h7.jpg",
		},
		Dish{
			ID: "23", Name: "Càng Ghẹ Rang Muối",
			Description: "Càng ghẹ tươi rang muối thơm đậm đà",
			Price: 50000, CookingMethod: "stir-fried", Rating: 4.6, Featured: false,
			ImageURL: "http://localhost:8080/images/h7.jpg",
		},
		Dish{
			ID: "24", Name: "Càng Ghẹ Xào Me",
			Description: "Càng ghẹ tươi xào me chua ngọt hấp dẫn",
			Price: 50000, CookingMethod: "stir-fried", Rating: 4.6, Featured: false,
			ImageURL: "http://localhost:8080/images/h7.jpg",
		},

		// ── ỐC LEN ────────────────────────────────────────────────
		Dish{
			ID: "25", Name: "Ốc Len Xào Dừa",
			Description: "Ốc len tươi xào nước dừa thơm ngọt",
			Price: 50000, CookingMethod: "stir-fried", Rating: 4.4, Featured: false,
			ImageURL: "http://localhost:8080/images/h5.jpg",
		},

		// ── SÒ LÔNG ───────────────────────────────────────────────
		Dish{
			ID: "26", Name: "Sò Lông Nướng Mỡ Hành",
			Description: "Sò lông tươi nướng với mỡ hành thơm phức",
			Price: 40000, CookingMethod: "grilled", Rating: 4.8, Featured: true,
			ImageURL: "http://localhost:8080/images/h2.jpg",
		},
		Dish{
			ID: "27", Name: "Sò Lông Nướng Mắm Tỏi",
			Description: "Sò lông tươi nướng với mắm tỏi đậm đà",
			Price: 40000, CookingMethod: "grilled", Rating: 4.7, Featured: false,
			ImageURL: "http://localhost:8080/images/h2.jpg",
		},

		// ── HÀU ───────────────────────────────────────────────────
		Dish{
			ID: "28", Name: "Hàu Nướng Phô Mai",
			Description: "Hàu tươi nướng với phô mai béo ngậy",
			Price: 40000, CookingMethod: "grilled", Rating: 4.6, Featured: true,
			ImageURL: "http://localhost:8080/images/h11.jpg",
		},
		Dish{
			ID: "29", Name: "Hàu Nướng Mỡ Hành",
			Description: "Hàu tươi nướng với mỡ hành thơm phức",
			Price: 40000, CookingMethod: "grilled", Rating: 4.5, Featured: false,
			ImageURL: "http://localhost:8080/images/h11.jpg",
		},

		// ── SÒ ĐIỆP ───────────────────────────────────────────────
		Dish{
			ID: "30", Name: "Sò Điệp Nướng Phô Mai",
			Description: "Sò điệp tươi nướng với phô mai béo ngậy thơm lừng",
			Price: 50000, CookingMethod: "grilled", Rating: 4.8, Featured: true,
			ImageURL: "http://localhost:8080/images/h6.jpg",
		},
		Dish{
			ID: "31", Name: "Sò Điệp Nướng Mỡ Hành",
			Description: "Sò điệp tươi nướng với mỡ hành thơm phức",
			Price: 50000, CookingMethod: "grilled", Rating: 4.7, Featured: false,
			ImageURL: "http://localhost:8080/images/h6.jpg",
		},

		// ── ỐC BÚA ────────────────────────────────────────────────
		Dish{
			ID: "32", Name: "Ốc Búa Nướng Mắm Tỏi",
			Description: "Ốc búa tươi nướng với mắm tỏi đậm đà",
			Price: 40000, CookingMethod: "grilled", Rating: 4.4, Featured: false,
			ImageURL: "http://localhost:8080/images/h4.jpg",
		},

		// ── ỐC TỎI ────────────────────────────────────────────────
		Dish{
			ID: "33", Name: "Ốc Tỏi Nướng Mắm Tỏi",
			Description: "Ốc tỏi tươi nướng với mắm tỏi thơm nồng",
			Price: 40000, CookingMethod: "grilled", Rating: 4.4, Featured: false,
			ImageURL: "http://localhost:8080/images/h4.jpg",
		},
		Dish{
			ID: "34", Name: "Ốc Tỏi Nướng Mỡ Hành",
			Description: "Ốc tỏi tươi nướng với mỡ hành thơm phức",
			Price: 40000, CookingMethod: "grilled", Rating: 4.4, Featured: false,
			ImageURL: "http://localhost:8080/images/h4.jpg",
		},

		// ── CÀ NA ─────────────────────────────────────────────────
		Dish{
			ID: "35", Name: "Cà Na Bơ Bắp",
			Description: "Cà na xào bơ bắp béo ngọt lạ miệng",
			Price: 40000, CookingMethod: "stir-fried", Rating: 4.2, Featured: false,
			ImageURL: "http://localhost:8080/images/h10.jpg",
		},
		Dish{
			ID: "36", Name: "Cà Na Cháy Tỏi",
			Description: "Cà na rang cháy tỏi vàng giòn thơm",
			Price: 40000, CookingMethod: "stir-fried", Rating: 4.3, Featured: false,
			ImageURL: "http://localhost:8080/images/h10.jpg",
		},
		Dish{
			ID: "37", Name: "Cà Na Xào Me",
			Description: "Cà na xào me chua ngọt đậm đà",
			Price: 40000, CookingMethod: "stir-fried", Rating: 4.3, Featured: false,
			ImageURL: "http://localhost:8080/images/h10.jpg",
		},
		Dish{
			ID: "38", Name: "Cà Na Rang Muối",
			Description: "Cà na rang muối thơm giòn hấp dẫn",
			Price: 40000, CookingMethod: "stir-fried", Rating: 4.2, Featured: false,
			ImageURL: "http://localhost:8080/images/h10.jpg",
		},

		// ── ỐC BƯƠU ───────────────────────────────────────────────
		Dish{
			ID: "39", Name: "Ốc Bươu Nướng Tiêu Xanh",
			Description: "Ốc bươu tươi nướng với tiêu xanh thơm cay",
			Price: 40000, CookingMethod: "grilled", Rating: 4.5, Featured: false,
			ImageURL: "http://localhost:8080/images/h4.jpg",
		},
	}

	result, err := collection.InsertMany(context.Background(), dishes)
	if err != nil {
		log.Fatalf("Failed to insert dishes: %v", err)
	}

	fmt.Printf("✅ Đã insert thành công %d món vào MongoDB\n", len(result.InsertedIDs))
	fmt.Println("🎉 Seed data hoàn tất!")
}

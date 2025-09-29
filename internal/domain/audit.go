// package domain

// import (
// 	"time"
// )

// type Dish struct {
// 	ID            string  `bson:"_id" json:"id"`
// 	Name          string  `bson:"name" json:"name"`
// 	Price         float64 `bson:"price" json:"price"`
// 	Description   string  `bson:"description" json:"description"`
// 	ImageURL      string  `bson:"image_url" json:"image_url"`
// 	CookingMethod string  `bson:"cooking_method" json:"cooking_method"`
// 	Rating        float64 `bson:"rating" json:"rating"`
// 	Featured      bool    `bson:"featured" json:"featured"`
// }

// type DishFilter struct {
// 	CookingMethod string `json:"cooking_method"`
// 	Page          int    `json:"page"`
// 	Limit         int    `json:"limit"`
// }

// type DishListResponse struct {
// 	Success bool `json:"success"`
// 	Data    struct {
// 		Dishes     []Dish             `json:"dishes"`
// 		Pagination PaginationResponse `json:"pagination"`
// 	} `json:"data"`
// }

// type DishCardResponse struct {
// 	Img           string   `json:"img"`
// 	Title         string   `json:"title"`
// 	Rating        *float64 `json:"rating"`
// 	Price         *float64 `json:"price"`
// 	Description   *string  `json:"description"`
// 	CookingMethod *string  `json:"cooking_method"`
// 	ShowDetails   *bool    `json:"show_details"`
// }

// type DishCardListResponse struct {
// 	Success bool `json:"success"`
// 	Data    struct {
// 		Dishes     []DishCardResponse `json:"dishes"`
// 		Pagination PaginationResponse `json:"pagination"`
// 	} `json:"data"`
// }

// type PaginationResponse struct {
// 	Page       int   `json:"page"`
// 	Limit      int   `json:"limit"`
// 	Total      int64 `json:"total"`
// 	TotalPages int   `json:"total_pages"`
// }

// type AuditLog struct {
// 	ID        string
// 	Action    string
// 	UserID    string
// 	CreatedAt time.Time
// 	Metadata  map[string]string
// }

// type Order struct {
// 	ID          string      `bson:"_id" json:"id"`
// 	UserID      string      `bson:"user_id" json:"user_id"`
// 	TotalPrice  float64     `bson:"total_price" json:"total_price"`
// 	Status      string      `bson:"status" json:"status"`
// 	CreatedAt   time.Time   `bson:"created_at" json:"created_at"`
// 	Items       []OrderItem `bson:"items" json:"items"`
// 	TableNumber int         `bson:"table_number" json:"table_number"` // Mới: Số bàn (1-20)
// }

// type OrderItem struct {
// 	DishID   string  `bson:"dish_id" json:"dish_id"`
// 	Quantity int     `bson:"quantity" json:"quantity"`
// 	Price    float64 `bson:"price" json:"price"`
// 	Title    string  `bson:"title" json:"title"`
// }

// // Request model từ frontend (cập nhật để nhận table_number)
// type CreateOrderRequest struct {
// 	UserID      string                   `json:"user_id"`
// 	Items       []CreateOrderRequestItem `json:"items"`
// 	Total       float64                  `json:"total"`
// 	CreatedAt   time.Time                `json:"created_at"`
// 	TableNumber int                      `json:"table_number"` // Mới: Số bàn từ frontend
// }

// type CreateOrderRequestItem struct {
// 	ID            string  `json:"id"`
// 	Title         string  `json:"title"`
// 	Price         float64 `json:"price"`
// 	Quantity      int     `json:"quantity"`
// 	Description   string  `json:"description"`
// 	ImageURL      string  `json:"image_url"`
// 	CookingMethod string  `json:"cooking_method"`
// 	Rating        float64 `json:"rating"`
// }

// // Response cho filter theo bàn (nếu cần)
// type OrdersByTableResponse struct {
// 	Success bool    `json:"success"`
// 	Data    []Order `json:"data"`
// }

package domain

import (
	"time"
)

type OrderStatus string

const (
	StatusPending   OrderStatus = "pending"
	StatusPreparing OrderStatus = "preparing"
	StatusReady     OrderStatus = "ready"
	StatusCompleted OrderStatus = "completed"
	StatusCancelled OrderStatus = "cancelled"
)

type Dish struct {
	ID            string  `bson:"_id" json:"id"`
	Name          string  `bson:"name" json:"name"`
	Price         float64 `bson:"price" json:"price"`
	Description   string  `bson:"description" json:"description"`
	ImageURL      string  `bson:"image_url" json:"image_url"`
	CookingMethod string  `bson:"cooking_method" json:"cooking_method"`
	Rating        float64 `bson:"rating" json:"rating"`
	Featured      bool    `bson:"featured" json:"featured"`
}

type DishFilter struct {
	CookingMethod string `json:"cooking_method"`
	Page          int    `json:"page"`
	Limit         int    `json:"limit"`
}

type DishListResponse struct {
	Success bool `json:"success"`
	Data    struct {
		Dishes     []Dish             `json:"dishes"`
		Pagination PaginationResponse `json:"pagination"`
	} `json:"data"`
}

type DishCardResponse struct {
	Img           string   `json:"img"`
	Title         string   `json:"title"`
	Rating        *float64 `json:"rating"`
	Price         *float64 `json:"price"`
	Description   *string  `json:"description"`
	CookingMethod *string  `json:"cooking_method"`
	ShowDetails   *bool    `json:"show_details"`
}

type DishCardListResponse struct {
	Success bool `json:"success"`
	Data    struct {
		Dishes     []DishCardResponse `json:"dishes"`
		Pagination PaginationResponse `json:"pagination"`
	} `json:"data"`
}

type PaginationResponse struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

type AuditLog struct {
	ID        string
	Action    string
	UserID    string
	CreatedAt time.Time
	Metadata  map[string]string
}

type Order struct {
	ID          string      `bson:"_id" json:"id"`
	UserID      string      `bson:"user_id" json:"user_id"`
	TableNumber int         `bson:"table_number" json:"table_number"`
	TotalPrice  float64     `bson:"total_price" json:"total_price"`
	Status      OrderStatus `bson:"status" json:"status"`
	CreatedAt   time.Time   `bson:"created_at" json:"created_at"`
	Items       []OrderItem `bson:"items" json:"items"`
	UpdatedAt   time.Time   `bson:"updated_at" json:"updated_at"`
}

type OrderItem struct {
	DishID   string  `bson:"dish_id" json:"dish_id"`
	Quantity int     `bson:"quantity" json:"quantity"`
	Price    float64 `bson:"price" json:"price"`
	Title    string  `bson:"title" json:"title"`
}

// Request model từ frontend
type CreateOrderRequest struct {
	UserID      string                   `json:"user_id"`
	TableNumber int                      `json:"table_number"` // Mới: Nhận từ frontend
	Items       []CreateOrderRequestItem `json:"items"`
	Total       float64                  `json:"total"`
	CreatedAt   time.Time                `json:"created_at"`
}

type CreateOrderRequestItem struct {
	ID            string  `json:"id"`
	Title         string  `json:"title"`
	Price         float64 `json:"price"`
	Quantity      int     `json:"quantity"`
	Description   string  `json:"description"`
	ImageURL      string  `json:"image_url"`
	CookingMethod string  `json:"cooking_method"`
	Rating        float64 `json:"rating"`
}

type UpdateOrderRequest struct {
	Items  []UpdateOrderRequestItem `json:"items"`
	Status OrderStatus              `json:"status"`
}

type UpdateOrderRequestItem struct {
	ID       string  `json:"id"`
	Title    string  `json:"title"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
}

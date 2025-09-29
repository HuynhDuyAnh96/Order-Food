// package repository

// import (
// 	"fms_audit/internal/domain"
// 	"math"
// 	"sort"
// 	"strings"
// 	"time"
// )

// type DishRepository struct {
// 	dishes []domain.Dish
// }

// func NewDishRepository() *DishRepository {
// 	repo := &DishRepository{}
// 	repo.seedData()
// 	return repo
// }

// func (r *DishRepository) seedData() {
// 	r.dishes = []domain.Dish{
// 		{
// 			ID:            "1",
// 			Name:          "Sò điệp nướng phô mai",
// 			Description:   "Sò điệp tươi nướng với phô mai béo ngậy, thơm lừng",
// 			Price:         85000,
// 			Category:      "seafood",
// 			CookingMethod: "grilled",
// 			IsPopular:     true,
// 			Rating:        4.8,
// 			ImageURL:      "http://localhost:8080/images/h6.jpg",
// 			CreatedAt:     time.Now().AddDate(0, -1, 0),
// 			UpdatedAt:     time.Now(),
// 		},
// 		{
// 			ID:            "2",
// 			Name:          "Càng ghẹ rang me",
// 			Description:   "Càng ghẹ tươi rang với mắm ruốc đậm đà, cay nồng",
// 			Price:         90000,
// 			Category:      "seafood",
// 			CookingMethod: "stir-fried",
// 			IsPopular:     true,
// 			Rating:        4.7,
// 			ImageURL:      "http://localhost:8080/images/h7.jpg",
// 			CreatedAt:     time.Now().AddDate(0, -1, -10),
// 			UpdatedAt:     time.Now(),
// 		},
// 		{
// 			ID:            "3",
// 			Name:          "Sò lông nướng mở hành",
// 			Description:   "Tôm hùm tươi nướng với bơ tỏi thơm phức",
// 			Price:         150000,
// 			Category:      "seafood",
// 			CookingMethod: "grilled",
// 			IsPopular:     true,
// 			Rating:        4.9,
// 			ImageURL:      "http://localhost:8080/images/h2.jpg",
// 			CreatedAt:     time.Now().AddDate(0, -1, -15),
// 			UpdatedAt:     time.Now(),
// 		},
// 		{
// 			ID:            "4",
// 			Name:          "Ốc hương xào bắp",
// 			Description:   "Cua biển rang me chua ngọt đậm đà",
// 			Price:         120000,
// 			Category:      "seafood",
// 			CookingMethod: "stir-fried",
// 			IsPopular:     false,
// 			Rating:        4.5,
// 			ImageURL:      "http://localhost:8080/images/h1.jpg",
// 			CreatedAt:     time.Now().AddDate(0, -3, 0),
// 			UpdatedAt:     time.Now(),
// 		},
// 		{
// 			ID:            "5",
// 			Name:          "Ốc bươu nướng tiêu",
// 			Description:   "Ốc hương tươi nướng với mỡ hành thơm lừng",
// 			Price:         75000,
// 			Category:      "seafood",
// 			CookingMethod: "grilled",
// 			IsPopular:     true,
// 			Rating:        4.6,
// 			ImageURL:      "http://localhost:8080/images/h4.jpg",
// 			CreatedAt:     time.Now().AddDate(0, -2, -10),
// 			UpdatedAt:     time.Now(),
// 		},
// 		{
// 			ID:            "6",
// 			Name:          "Ốc len xào dừa",
// 			Description:   "Mực tươi nướng với sa tế cay nồng",
// 			Price:         95000,
// 			Category:      "seafood",
// 			CookingMethod: "steamed",
// 			IsPopular:     false,
// 			Rating:        4.4,
// 			ImageURL:      "http://localhost:8080/images/h5.jpg",
// 			CreatedAt:     time.Now().AddDate(0, -1, -5),
// 			UpdatedAt:     time.Now(),
// 		},
// 		{
// 			ID:            "7",
// 			Name:          "Mực xào mì",
// 			Description:   "Cá lăng tươi nướng trong lá chuối giữ nguyên hương vị",
// 			Price:         110000,
// 			Category:      "seafood",
// 			CookingMethod: "stir-fried",
// 			IsPopular:     false,
// 			Rating:        4.3,
// 			ImageURL:      "http://localhost:8080/images/h8.jpg",
// 			CreatedAt:     time.Now().AddDate(0, -2, -20),
// 			UpdatedAt:     time.Now(),
// 		},
// 		{
// 			ID:            "8",
// 			Name:          "Hào nướng phô mai",
// 			Description:   "Nghêu tươi hấp với sả thơm, thanh mát",
// 			Price:         65000,
// 			Category:      "seafood",
// 			CookingMethod: "grilled",
// 			IsPopular:     true,
// 			Rating:        4.5,
// 			ImageURL:      "http://localhost:8080/images/h11.jpg",
// 			CreatedAt:     time.Now().AddDate(0, -1, -8),
// 			UpdatedAt:     time.Now(),
// 		},
// 	}
// }

// func (r *DishRepository) GetAll(filter domain.DishFilter) ([]domain.Dish, int, error) {
// 	filteredDishes := r.filterDishes(filter)
// 	total := len(filteredDishes)

// 	// Sort dishes
// 	r.sortDishes(filteredDishes, filter.Sort)

// 	// Apply pagination
// 	start := (filter.Page - 1) * filter.Limit
// 	end := start + filter.Limit

// 	if start >= len(filteredDishes) {
// 		return []domain.Dish{}, total, nil
// 	}

// 	if end > len(filteredDishes) {
// 		end = len(filteredDishes)
// 	}

// 	return filteredDishes[start:end], total, nil
// }

// func (r *DishRepository) GetFeatured(limit int) ([]domain.Dish, error) {
// 	var popularDishes []domain.Dish

// 	for _, dish := range r.dishes {
// 		if dish.IsPopular {
// 			popularDishes = append(popularDishes, dish)
// 		}
// 	}

// 	// Sort by rating descending
// 	sort.Slice(popularDishes, func(i, j int) bool {
// 		return popularDishes[i].Rating > popularDishes[j].Rating
// 	})

// 	if len(popularDishes) > limit {
// 		popularDishes = popularDishes[:limit]
// 	}

// 	return popularDishes, nil
// }

// func (r *DishRepository) filterDishes(filter domain.DishFilter) []domain.Dish {
// 	var filtered []domain.Dish

// 	for _, dish := range r.dishes {
// 		// Filter by category
// 		if filter.Category != "" && !strings.EqualFold(dish.Category, filter.Category) {
// 			continue
// 		}

// 		// Filter by cooking method
// 		if filter.CookingMethod != "" && !strings.EqualFold(dish.CookingMethod, filter.CookingMethod) {
// 			continue
// 		}

// 		// Filter by popularity
// 		if filter.IsPopular != nil && dish.IsPopular != *filter.IsPopular {
// 			continue
// 		}

// 		filtered = append(filtered, dish)
// 	}

// 	return filtered
// }

// func (r *DishRepository) sortDishes(dishes []domain.Dish, sortBy string) {
// 	switch sortBy {
// 	case "price_asc":
// 		sort.Slice(dishes, func(i, j int) bool {
// 			return dishes[i].Price < dishes[j].Price
// 		})
// 	case "price_desc":
// 		sort.Slice(dishes, func(i, j int) bool {
// 			return dishes[i].Price > dishes[j].Price
// 		})
// 	case "rating_desc":
// 		sort.Slice(dishes, func(i, j int) bool {
// 			return dishes[i].Rating > dishes[j].Rating
// 		})
// 	default:
// 		// Default sort by name
// 		sort.Slice(dishes, func(i, j int) bool {
// 			return dishes[i].Name < dishes[j].Name
// 		})
// 	}
// }

//	func CalculateTotalPages(total, limit int) int {
//		return int(math.Ceil(float64(total) / float64(limit)))
//	}
package repository

import (
	"context"
	"fms_audit/internal/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DishRepository struct {
	collection *mongo.Collection
}

func NewDishRepository(db *mongo.Database) *DishRepository {
	return &DishRepository{
		collection: db.Collection("dishes"),
	}
}

func (r *DishRepository) GetAll(ctx context.Context, filter domain.DishFilter) ([]domain.Dish, int64, error) {
	findOptions := options.Find()
	findOptions.SetSkip(int64((filter.Page - 1) * filter.Limit))
	findOptions.SetLimit(int64(filter.Limit))

	query := bson.M{}
	if filter.CookingMethod != "" {
		query["cooking_method"] = filter.CookingMethod
	}

	// Đếm tổng số bản ghi
	total, err := r.collection.CountDocuments(ctx, query)
	if err != nil {
		return nil, 0, err
	}

	// Lấy danh sách món ăn
	cursor, err := r.collection.Find(ctx, query, findOptions)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var dishes []domain.Dish
	if err := cursor.All(ctx, &dishes); err != nil {
		return nil, 0, err
	}

	return dishes, total, nil
}

func (r *DishRepository) GetFeatured(ctx context.Context, limit int) ([]domain.Dish, error) {
	findOptions := options.Find()
	findOptions.SetLimit(int64(limit))

	cursor, err := r.collection.Find(ctx, bson.M{"featured": true}, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var dishes []domain.Dish
	if err := cursor.All(ctx, &dishes); err != nil {
		return nil, err
	}
	return dishes, nil
}

func (r *DishRepository) GetStirFriedDishes(ctx context.Context) ([]domain.Dish, error) {
	return r.getDishesByCookingMethod(ctx, "stir-fried")
}

func (r *DishRepository) GetSteamedDishes(ctx context.Context) ([]domain.Dish, error) {
	return r.getDishesByCookingMethod(ctx, "steamed")
}

func (r *DishRepository) GetGrilledDishes(ctx context.Context) ([]domain.Dish, error) {
	return r.getDishesByCookingMethod(ctx, "grilled")
}

func (r *DishRepository) getDishesByCookingMethod(ctx context.Context, method string) ([]domain.Dish, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"cooking_method": method})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var dishes []domain.Dish
	if err := cursor.All(ctx, &dishes); err != nil {
		return nil, err
	}
	return dishes, nil
}

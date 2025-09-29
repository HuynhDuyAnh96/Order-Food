// package repository

// import (
// 	"context"
// 	"fms_audit/internal/domain"

// 	"go.mongodb.org/mongo-driver/bson"
// 	"go.mongodb.org/mongo-driver/mongo"
// )

// type OrderRepository struct {
// 	collection *mongo.Collection
// }

// func NewOrderRepository(db *mongo.Database) *OrderRepository {
// 	return &OrderRepository{collection: db.Collection("orders")}
// }

// func (r *OrderRepository) CreateOrder(ctx context.Context, order *domain.Order) error {
// 	_, err := r.collection.InsertOne(ctx, order)
// 	return err
// }

// func (r *OrderRepository) GetAllOrders(ctx context.Context) ([]domain.Order, error) {
// 	cursor, err := r.collection.Find(ctx, bson.M{})
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer cursor.Close(ctx)

// 	var orders []domain.Order
// 	if err = cursor.All(ctx, &orders); err != nil {
// 		return nil, err
// 	}
// 	return orders, nil
// }

// // Mới: Lấy orders theo table_number
// func (r *OrderRepository) GetOrdersByTable(ctx context.Context, tableNumber int) ([]domain.Order, error) {
// 	filter := bson.M{"table_number": tableNumber}
// 	cursor, err := r.collection.Find(ctx, filter)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer cursor.Close(ctx)

// 	var orders []domain.Order
// 	if err = cursor.All(ctx, &orders); err != nil {
// 		return nil, err
// 	}
// 	return orders, nil
// }

package repository

import (
	"context"
	"fms_audit/internal/domain"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type OrderRepository struct {
	collection *mongo.Collection
}

func NewOrderRepository(db *mongo.Database) *OrderRepository {
	return &OrderRepository{collection: db.Collection("orders")}
}

func (r *OrderRepository) CreateOrder(ctx context.Context, order *domain.Order) error {
	fmt.Printf("Inserting order: %+v\n", order)
	_, err := r.collection.InsertOne(ctx, order)
	if err != nil {
		return fmt.Errorf("insert failed: %w", err)
	}
	return nil
}

func (r *OrderRepository) GetAllOrders(ctx context.Context) ([]domain.Order, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var orders []domain.Order
	err = cursor.All(ctx, &orders)
	return orders, err
}

// Mới: Lấy orders theo bàn
func (r *OrderRepository) GetOrdersByTable(ctx context.Context, tableNumber int) ([]domain.Order, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"table_number": tableNumber})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var orders []domain.Order
	err = cursor.All(ctx, &orders)
	return orders, err
}

func (r *OrderRepository) GetOrderByID(ctx context.Context, orderID string) (*domain.Order, error) {
	filter := bson.M{"_id": orderID}
	order := &domain.Order{}
	err := r.collection.FindOne(ctx, filter).Decode(order)
	if err != nil {
		return nil, err
	}
	return order, nil
}

// func (r *OrderRepository) UpdateOrder(ctx context.Context, order *domain.Order) error {
// 	filter := bson.M{"_id": order.ID}
// 	update := bson.M{"$set": bson.M{
// 		"items":       order.Items,
// 		"total_price": order.TotalPrice,
// 		"status":      order.Status,
// 	}}
// 	_, err := r.collection.UpdateOne(ctx, filter, update)
// 	return err
// }

func (r *OrderRepository) UpdateOrder(ctx context.Context, order *domain.Order) error {
	filter := bson.M{"_id": order.ID}
	update := bson.M{"$set": bson.M{
		"items":       order.Items,
		"total_price": order.TotalPrice,
		"status":      order.Status,
		"updated_at":  order.UpdatedAt, // Thêm updated_at
	}}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

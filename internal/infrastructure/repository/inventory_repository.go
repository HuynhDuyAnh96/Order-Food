package repository

import (
	"context"
	"fms_audit/internal/domain"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type InventoryRepository struct {
	ingredientTypes *mongo.Collection
	receipts        *mongo.Collection
	batches         *mongo.Collection
	recipeCosts     *mongo.Collection
	sessions        *mongo.Collection
}

func NewInventoryRepository(db *mongo.Database) *InventoryRepository {
	return &InventoryRepository{
		ingredientTypes: db.Collection("ingredient_types"),
		receipts:        db.Collection("inventory_receipts"),
		batches:         db.Collection("processing_batches"),
		recipeCosts:     db.Collection("recipe_costs"),
		sessions:        db.Collection("evening_sessions"),
	}
}

// ── IngredientType ────────────────────────────────────────────────────────────

func (r *InventoryRepository) CreateIngredientType(ctx context.Context, ing *domain.IngredientType) error {
	_, err := r.ingredientTypes.InsertOne(ctx, ing)
	return err
}

func (r *InventoryRepository) GetAllIngredientTypes(ctx context.Context) ([]domain.IngredientType, error) {
	opts := options.Find().SetSort(bson.D{{Key: "name", Value: 1}})
	cursor, err := r.ingredientTypes.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var result []domain.IngredientType
	return result, cursor.All(ctx, &result)
}

func (r *InventoryRepository) GetIngredientTypeByID(ctx context.Context, id string) (*domain.IngredientType, error) {
	var ing domain.IngredientType
	err := r.ingredientTypes.FindOne(ctx, bson.M{"_id": id}).Decode(&ing)
	if err != nil {
		return nil, err
	}
	return &ing, nil
}

// ── InventoryReceipt ──────────────────────────────────────────────────────────

func (r *InventoryRepository) CreateReceipt(ctx context.Context, receipt *domain.InventoryReceipt) error {
	_, err := r.receipts.InsertOne(ctx, receipt)
	return err
}

func (r *InventoryRepository) GetReceipts(ctx context.Context, ingredientTypeID string) ([]domain.InventoryReceipt, error) {
	filter := bson.M{}
	if ingredientTypeID != "" {
		filter["ingredient_type_id"] = ingredientTypeID
	}
	opts := options.Find().SetSort(bson.D{{Key: "received_at", Value: -1}})
	cursor, err := r.receipts.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var result []domain.InventoryReceipt
	return result, cursor.All(ctx, &result)
}

func (r *InventoryRepository) GetReceiptByID(ctx context.Context, id string) (*domain.InventoryReceipt, error) {
	var receipt domain.InventoryReceipt
	err := r.receipts.FindOne(ctx, bson.M{"_id": id}).Decode(&receipt)
	if err != nil {
		return nil, err
	}
	return &receipt, nil
}

// ── ProcessingBatch ───────────────────────────────────────────────────────────

func (r *InventoryRepository) CreateProcessingBatch(ctx context.Context, batch *domain.ProcessingBatch) error {
	_, err := r.batches.InsertOne(ctx, batch)
	return err
}

func (r *InventoryRepository) GetBatches(ctx context.Context, ingredientTypeID string) ([]domain.ProcessingBatch, error) {
	filter := bson.M{}
	if ingredientTypeID != "" {
		filter["ingredient_type_id"] = ingredientTypeID
	}
	opts := options.Find().SetSort(bson.D{{Key: "processed_at", Value: -1}})
	cursor, err := r.batches.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var result []domain.ProcessingBatch
	return result, cursor.All(ctx, &result)
}

// GetTotalOutputBaskets - tổng rổ đã xử lý cho 1 loại sò
func (r *InventoryRepository) GetTotalOutputBaskets(ctx context.Context, ingredientTypeID string) (int, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"ingredient_type_id": ingredientTypeID}}},
		{{Key: "$group", Value: bson.M{
			"_id":           nil,
			"total_baskets": bson.M{"$sum": "$output_baskets"},
		}}},
	}
	cursor, err := r.batches.Aggregate(ctx, pipeline)
	if err != nil {
		return 0, err
	}
	defer cursor.Close(ctx)
	var results []struct {
		TotalBaskets int `bson:"total_baskets"`
	}
	if err = cursor.All(ctx, &results); err != nil || len(results) == 0 {
		return 0, err
	}
	return results[0].TotalBaskets, nil
}

// GetLatestCostPerBasket - lấy chi phí/rổ từ mẻ xử lý gần nhất
func (r *InventoryRepository) GetLatestCostPerBasket(ctx context.Context, ingredientTypeID string) (float64, error) {
	opts := options.FindOne().SetSort(bson.D{{Key: "processed_at", Value: -1}})
	var batch domain.ProcessingBatch
	err := r.batches.FindOne(ctx, bson.M{"ingredient_type_id": ingredientTypeID}, opts).Decode(&batch)
	if err != nil {
		return 0, nil // không có mẻ nào thì cost = 0
	}
	return batch.CostPerBasket, nil
}

// ── RecipeCost ────────────────────────────────────────────────────────────────

func (r *InventoryRepository) UpsertRecipeCost(ctx context.Context, rc *domain.RecipeCost) error {
	filter := bson.M{"ingredient_type_id": rc.IngredientTypeID}
	update := bson.M{"$set": rc}
	opts := options.Update().SetUpsert(true)
	_, err := r.recipeCosts.UpdateOne(ctx, filter, update, opts)
	return err
}

func (r *InventoryRepository) GetRecipeCost(ctx context.Context, ingredientTypeID string) (*domain.RecipeCost, error) {
	var rc domain.RecipeCost
	err := r.recipeCosts.FindOne(ctx, bson.M{"ingredient_type_id": ingredientTypeID}).Decode(&rc)
	if err != nil {
		return nil, err // mongo.ErrNoDocuments nếu chưa có
	}
	return &rc, nil
}

// ── EveningSession ────────────────────────────────────────────────────────────

func (r *InventoryRepository) CreateSession(ctx context.Context, session *domain.EveningSession) error {
	_, err := r.sessions.InsertOne(ctx, session)
	return err
}

func (r *InventoryRepository) GetSessionByID(ctx context.Context, id string) (*domain.EveningSession, error) {
	var session domain.EveningSession
	err := r.sessions.FindOne(ctx, bson.M{"_id": id}).Decode(&session)
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *InventoryRepository) GetOpenSession(ctx context.Context) (*domain.EveningSession, error) {
	var session domain.EveningSession
	err := r.sessions.FindOne(ctx, bson.M{"status": domain.SessionOpen}).Decode(&session)
	if err != nil {
		return nil, err // mongo.ErrNoDocuments = không có ca đang mở
	}
	return &session, nil
}

func (r *InventoryRepository) GetSessions(ctx context.Context, limit int) ([]domain.EveningSession, error) {
	opts := options.Find().SetSort(bson.D{{Key: "opened_at", Value: -1}})
	if limit > 0 {
		opts.SetLimit(int64(limit))
	}
	cursor, err := r.sessions.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var result []domain.EveningSession
	return result, cursor.All(ctx, &result)
}

func (r *InventoryRepository) CloseSession(ctx context.Context, sessionID string, usage []domain.SessionUsage, note string) error {
	now := time.Now()
	filter := bson.M{"_id": sessionID}
	update := bson.M{"$set": bson.M{
		"status":    domain.SessionClosed,
		"usage":     usage,
		"closed_at": now,
		"note":      note,
	}}
	result, err := r.sessions.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("close session failed: %w", err)
	}
	if result.MatchedCount == 0 {
		return fmt.Errorf("session %s not found", sessionID)
	}
	return nil
}

// GetTotalUsedBaskets - tổng rổ đã dùng cho 1 loại sò (từ các ca đã đóng)
func (r *InventoryRepository) GetTotalUsedBaskets(ctx context.Context, ingredientTypeID string) (int, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"status": domain.SessionClosed}}},
		{{Key: "$unwind", Value: "$usage"}},
		{{Key: "$match", Value: bson.M{"usage.ingredient_type_id": ingredientTypeID}}},
		{{Key: "$group", Value: bson.M{
			"_id":          nil,
			"total_used":   bson.M{"$sum": "$usage.used_baskets"},
			"total_wasted": bson.M{"$sum": "$usage.wasted_baskets"},
		}}},
	}
	cursor, err := r.sessions.Aggregate(ctx, pipeline)
	if err != nil {
		return 0, err
	}
	defer cursor.Close(ctx)
	var results []struct {
		TotalUsed   int `bson:"total_used"`
		TotalWasted int `bson:"total_wasted"`
	}
	if err = cursor.All(ctx, &results); err != nil || len(results) == 0 {
		return 0, err
	}
	return results[0].TotalUsed + results[0].TotalWasted, nil
}

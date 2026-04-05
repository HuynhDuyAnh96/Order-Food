package domain

import "time"

// IngredientType - loại nguyên liệu (sò huyết, sò lông, ốc len, ...)
type IngredientType struct {
	ID             string    `bson:"_id" json:"id"`
	Name           string    `bson:"name" json:"name"`                         // "Sò huyết"
	AvgKgPerBasket float64   `bson:"avg_kg_per_basket" json:"avg_kg_per_basket"` // ~3kg/rổ, dùng để ước tính chi phí
	CreatedAt      time.Time `bson:"created_at" json:"created_at"`
}

// InventoryReceipt - phiếu nhập hàng (mua X kg với giá Y)
type InventoryReceipt struct {
	ID               string    `bson:"_id" json:"id"`
	IngredientTypeID string    `bson:"ingredient_type_id" json:"ingredient_type_id"`
	IngredientName   string    `bson:"ingredient_name" json:"ingredient_name"` // denormalized
	RawWeightKg      float64   `bson:"raw_weight_kg" json:"raw_weight_kg"`     // kg mua vào
	PricePerKg       float64   `bson:"price_per_kg" json:"price_per_kg"`       // giá/kg
	TotalCost        float64   `bson:"total_cost" json:"total_cost"`           // = RawWeightKg * PricePerKg
	ReceivedAt       time.Time `bson:"received_at" json:"received_at"`
	Note             string    `bson:"note,omitempty" json:"note,omitempty"`
}

// ProcessingBatch - mẻ xử lý: từ X kg → ra Y rổ
type ProcessingBatch struct {
	ID               string    `bson:"_id" json:"id"`
	ReceiptID        string    `bson:"receipt_id" json:"receipt_id"`
	IngredientTypeID string    `bson:"ingredient_type_id" json:"ingredient_type_id"`
	IngredientName   string    `bson:"ingredient_name" json:"ingredient_name"` // denormalized
	InputWeightKg    float64   `bson:"input_weight_kg" json:"input_weight_kg"` // kg đưa vào xử lý
	OutputBaskets    int       `bson:"output_baskets" json:"output_baskets"`   // ra được bao nhiêu rổ
	WastePercent     float64   `bson:"waste_percent" json:"waste_percent"`     // % hao hụt, tính tự động
	CostPerBasket    float64   `bson:"cost_per_basket" json:"cost_per_basket"` // chi phí sò/rổ, tính từ receipt
	ProcessedAt      time.Time `bson:"processed_at" json:"processed_at"`
	Note             string    `bson:"note,omitempty" json:"note,omitempty"`
}

// RecipeCostItem - 1 loại gia vị/nguyên liệu phụ
type RecipeCostItem struct {
	Name       string  `bson:"name" json:"name"`               // "Tỏi", "Ớt", "Dầu ăn", "Than"
	CostAmount float64 `bson:"cost_amount" json:"cost_amount"` // chi phí cho 1 rổ (đồng)
	Note       string  `bson:"note,omitempty" json:"note,omitempty"` // "~50g tỏi ≈ 3,000đ"
}

// RecipeCost - chi phí nguyên liệu phụ cho 1 loại sò (không track tồn kho, chỉ track tiền)
type RecipeCost struct {
	ID                 string           `bson:"_id" json:"id"`
	IngredientTypeID   string           `bson:"ingredient_type_id" json:"ingredient_type_id"`
	IngredientName     string           `bson:"ingredient_name" json:"ingredient_name"` // denormalized
	Items              []RecipeCostItem `bson:"items" json:"items"`
	TotalCostPerBasket float64          `bson:"total_cost_per_basket" json:"total_cost_per_basket"` // tổng chi phí phụ/rổ
	UpdatedAt          time.Time        `bson:"updated_at" json:"updated_at"`
}

// SessionUsage - lượng rổ sử dụng cho mỗi loại sò trong 1 ca
type SessionUsage struct {
	IngredientTypeID string  `bson:"ingredient_type_id" json:"ingredient_type_id"`
	IngredientName   string  `bson:"ingredient_name" json:"ingredient_name"` // denormalized
	PlannedBaskets   int     `bson:"planned_baskets" json:"planned_baskets"` // chuẩn bị mấy rổ
	UsedBaskets      int     `bson:"used_baskets" json:"used_baskets"`       // thực tế dùng/bán
	WastedBaskets    int     `bson:"wasted_baskets" json:"wasted_baskets"`   // dư/hỏng bỏ
}

type SessionStatus string

const (
	SessionOpen   SessionStatus = "open"
	SessionClosed SessionStatus = "closed"
)

// EveningSession - ca bán tối
type EveningSession struct {
	ID        string        `bson:"_id" json:"id"`
	Date      string        `bson:"date" json:"date"` // "2024-01-15" để dễ query
	Status    SessionStatus `bson:"status" json:"status"`
	Usage     []SessionUsage `bson:"usage" json:"usage"`
	Note      string        `bson:"note,omitempty" json:"note,omitempty"`
	OpenedAt  time.Time     `bson:"opened_at" json:"opened_at"`
	ClosedAt  *time.Time    `bson:"closed_at,omitempty" json:"closed_at,omitempty"`
}

// StockSummaryItem - tồn kho của 1 loại sò
type StockSummaryItem struct {
	IngredientTypeID   string  `json:"ingredient_type_id"`
	IngredientName     string  `json:"ingredient_name"`
	AvailableBaskets   int     `json:"available_baskets"`    // rổ còn sẵn
	TotalCostPerBasket float64 `json:"total_cost_per_basket"` // giá thành thực tế/rổ (sò + gia vị)
}

// StockSummary - tổng quan tồn kho tất cả loại sò
type StockSummary struct {
	Items     []StockSummaryItem `json:"items"`
	UpdatedAt time.Time          `json:"updated_at"`
}

// ── Request types ────────────────────────────────────────────────────────────

type CreateIngredientTypeRequest struct {
	Name           string  `json:"name" binding:"required"`
	AvgKgPerBasket float64 `json:"avg_kg_per_basket"`
}

type CreateReceiptRequest struct {
	IngredientTypeID string    `json:"ingredient_type_id" binding:"required"`
	RawWeightKg      float64   `json:"raw_weight_kg" binding:"required"`
	PricePerKg       float64   `json:"price_per_kg" binding:"required"`
	ReceivedAt       time.Time `json:"received_at"`
	Note             string    `json:"note"`
}

type CreateProcessingBatchRequest struct {
	ReceiptID     string    `json:"receipt_id" binding:"required"`
	InputWeightKg float64   `json:"input_weight_kg" binding:"required"`
	OutputBaskets int       `json:"output_baskets" binding:"required"`
	ProcessedAt   time.Time `json:"processed_at"`
	Note          string    `json:"note"`
}

type UpsertRecipeCostRequest struct {
	IngredientTypeID string           `json:"ingredient_type_id" binding:"required"`
	Items            []RecipeCostItem `json:"items" binding:"required"`
}

type OpenSessionRequest struct {
	Date  string         `json:"date" binding:"required"` // "2024-01-15"
	Usage []SessionUsage `json:"usage" binding:"required"`
	Note  string         `json:"note"`
}

type CloseSessionRequest struct {
	Usage []SessionUsage `json:"usage" binding:"required"` // actual used/wasted
	Note  string         `json:"note"`
}

// ── Report types ─────────────────────────────────────────────────────────────

// YieldReportItem - báo cáo hao hụt cho 1 mẻ xử lý
type YieldReportItem struct {
	BatchID        string    `json:"batch_id"`
	IngredientName string    `json:"ingredient_name"`
	InputWeightKg  float64   `json:"input_weight_kg"`
	OutputBaskets  int       `json:"output_baskets"`
	WastePercent   float64   `json:"waste_percent"`
	CostPerBasket  float64   `json:"cost_per_basket"`
	ProcessedAt    time.Time `json:"processed_at"`
}

// CostReportItem - báo cáo chi phí tổng cho 1 ca
type CostReportItem struct {
	SessionID      string    `json:"session_id"`
	Date           string    `json:"date"`
	IngredientName string    `json:"ingredient_name"`
	UsedBaskets    int       `json:"used_baskets"`
	CostPerBasket  float64   `json:"cost_per_basket"`  // sò + gia vị
	TotalCost      float64   `json:"total_cost"`        // UsedBaskets * CostPerBasket
}

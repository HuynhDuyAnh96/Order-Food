package service

import (
	"context"
	"fms_audit/internal/domain"
	"fms_audit/internal/infrastructure/repository"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type InventoryService struct {
	repo *repository.InventoryRepository
}

func NewInventoryService(repo *repository.InventoryRepository) *InventoryService {
	return &InventoryService{repo: repo}
}

// ── IngredientType ────────────────────────────────────────────────────────────

func (s *InventoryService) CreateIngredientType(ctx context.Context, req domain.CreateIngredientTypeRequest) (*domain.IngredientType, error) {
	ing := &domain.IngredientType{
		ID:             fmt.Sprintf("ing_%d", time.Now().UnixNano()),
		Name:           req.Name,
		AvgKgPerBasket: req.AvgKgPerBasket,
		CreatedAt:      time.Now(),
	}
	if err := s.repo.CreateIngredientType(ctx, ing); err != nil {
		return nil, fmt.Errorf("tạo loại nguyên liệu thất bại: %w", err)
	}
	return ing, nil
}

func (s *InventoryService) GetIngredientTypes(ctx context.Context) ([]domain.IngredientType, error) {
	return s.repo.GetAllIngredientTypes(ctx)
}

// ── InventoryReceipt ──────────────────────────────────────────────────────────

func (s *InventoryService) CreateReceipt(ctx context.Context, req domain.CreateReceiptRequest) (*domain.InventoryReceipt, error) {
	if req.RawWeightKg <= 0 {
		return nil, fmt.Errorf("khối lượng phải lớn hơn 0")
	}
	if req.PricePerKg <= 0 {
		return nil, fmt.Errorf("giá/kg phải lớn hơn 0")
	}

	ing, err := s.repo.GetIngredientTypeByID(ctx, req.IngredientTypeID)
	if err != nil {
		return nil, fmt.Errorf("không tìm thấy loại nguyên liệu: %w", err)
	}

	receivedAt := req.ReceivedAt
	if receivedAt.IsZero() {
		receivedAt = time.Now()
	}

	receipt := &domain.InventoryReceipt{
		ID:               fmt.Sprintf("receipt_%d", time.Now().UnixNano()),
		IngredientTypeID: req.IngredientTypeID,
		IngredientName:   ing.Name,
		RawWeightKg:      req.RawWeightKg,
		PricePerKg:       req.PricePerKg,
		TotalCost:        req.RawWeightKg * req.PricePerKg,
		ReceivedAt:       receivedAt,
		Note:             req.Note,
	}

	if err := s.repo.CreateReceipt(ctx, receipt); err != nil {
		return nil, fmt.Errorf("lưu phiếu nhập thất bại: %w", err)
	}
	return receipt, nil
}

func (s *InventoryService) GetReceipts(ctx context.Context, ingredientTypeID string) ([]domain.InventoryReceipt, error) {
	return s.repo.GetReceipts(ctx, ingredientTypeID)
}

// ── ProcessingBatch ───────────────────────────────────────────────────────────

func (s *InventoryService) CreateProcessingBatch(ctx context.Context, req domain.CreateProcessingBatchRequest) (*domain.ProcessingBatch, error) {
	if req.InputWeightKg <= 0 {
		return nil, fmt.Errorf("khối lượng đầu vào phải lớn hơn 0")
	}
	if req.OutputBaskets <= 0 {
		return nil, fmt.Errorf("số rổ đầu ra phải lớn hơn 0")
	}

	receipt, err := s.repo.GetReceiptByID(ctx, req.ReceiptID)
	if err != nil {
		return nil, fmt.Errorf("không tìm thấy phiếu nhập: %w", err)
	}
	if req.InputWeightKg > receipt.RawWeightKg {
		return nil, fmt.Errorf("kg xử lý (%.2f) không thể vượt quá kg nhập (%.2f)", req.InputWeightKg, receipt.RawWeightKg)
	}

	// Tính % hao hụt dựa trên avg_kg_per_basket
	ing, err := s.repo.GetIngredientTypeByID(ctx, receipt.IngredientTypeID)
	if err != nil {
		return nil, fmt.Errorf("không tìm thấy loại nguyên liệu: %w", err)
	}

	// wastePercent = (input - ước tính output kg) / input * 100
	var wastePercent float64
	if ing.AvgKgPerBasket > 0 {
		estimatedOutputKg := float64(req.OutputBaskets) * ing.AvgKgPerBasket
		wastePercent = (req.InputWeightKg - estimatedOutputKg) / req.InputWeightKg * 100
		if wastePercent < 0 {
			wastePercent = 0
		}
	}

	// Chi phí sò/rổ = TotalCost của receipt / số rổ ra được
	costPerBasket := receipt.TotalCost / float64(req.OutputBaskets)

	processedAt := req.ProcessedAt
	if processedAt.IsZero() {
		processedAt = time.Now()
	}

	batch := &domain.ProcessingBatch{
		ID:               fmt.Sprintf("batch_%d", time.Now().UnixNano()),
		ReceiptID:        req.ReceiptID,
		IngredientTypeID: receipt.IngredientTypeID,
		IngredientName:   receipt.IngredientName,
		InputWeightKg:    req.InputWeightKg,
		OutputBaskets:    req.OutputBaskets,
		WastePercent:     wastePercent,
		CostPerBasket:    costPerBasket,
		ProcessedAt:      processedAt,
		Note:             req.Note,
	}

	if err := s.repo.CreateProcessingBatch(ctx, batch); err != nil {
		return nil, fmt.Errorf("lưu mẻ xử lý thất bại: %w", err)
	}
	return batch, nil
}

func (s *InventoryService) GetBatches(ctx context.Context, ingredientTypeID string) ([]domain.ProcessingBatch, error) {
	return s.repo.GetBatches(ctx, ingredientTypeID)
}

// ── RecipeCost ────────────────────────────────────────────────────────────────

func (s *InventoryService) UpsertRecipeCost(ctx context.Context, req domain.UpsertRecipeCostRequest) (*domain.RecipeCost, error) {
	if len(req.Items) == 0 {
		return nil, fmt.Errorf("phải có ít nhất 1 nguyên liệu phụ")
	}

	ing, err := s.repo.GetIngredientTypeByID(ctx, req.IngredientTypeID)
	if err != nil {
		return nil, fmt.Errorf("không tìm thấy loại nguyên liệu: %w", err)
	}

	var total float64
	for _, item := range req.Items {
		total += item.CostAmount
	}

	rc := &domain.RecipeCost{
		ID:                 fmt.Sprintf("recipe_%s", req.IngredientTypeID),
		IngredientTypeID:   req.IngredientTypeID,
		IngredientName:     ing.Name,
		Items:              req.Items,
		TotalCostPerBasket: total,
		UpdatedAt:          time.Now(),
	}

	if err := s.repo.UpsertRecipeCost(ctx, rc); err != nil {
		return nil, fmt.Errorf("lưu chi phí gia vị thất bại: %w", err)
	}
	return rc, nil
}

func (s *InventoryService) GetRecipeCost(ctx context.Context, ingredientTypeID string) (*domain.RecipeCost, error) {
	rc, err := s.repo.GetRecipeCost(ctx, ingredientTypeID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // chưa có thì trả về nil, không phải lỗi
		}
		return nil, err
	}
	return rc, nil
}

// ── Stock Summary ─────────────────────────────────────────────────────────────

func (s *InventoryService) GetStockSummary(ctx context.Context) (*domain.StockSummary, error) {
	types, err := s.repo.GetAllIngredientTypes(ctx)
	if err != nil {
		return nil, err
	}

	var items []domain.StockSummaryItem
	for _, ing := range types {
		totalOut, err := s.repo.GetTotalOutputBaskets(ctx, ing.ID)
		if err != nil {
			return nil, err
		}
		totalUsed, err := s.repo.GetTotalUsedBaskets(ctx, ing.ID)
		if err != nil {
			return nil, err
		}

		available := totalOut - totalUsed
		if available < 0 {
			available = 0
		}

		// Chi phí thực tế/rổ = chi phí sò (mẻ gần nhất) + chi phí gia vị
		costPerBasket, _ := s.repo.GetLatestCostPerBasket(ctx, ing.ID)
		rc, _ := s.repo.GetRecipeCost(ctx, ing.ID)
		if rc != nil {
			costPerBasket += rc.TotalCostPerBasket
		}

		items = append(items, domain.StockSummaryItem{
			IngredientTypeID:   ing.ID,
			IngredientName:     ing.Name,
			AvailableBaskets:   available,
			TotalCostPerBasket: costPerBasket,
		})
	}

	return &domain.StockSummary{
		Items:     items,
		UpdatedAt: time.Now(),
	}, nil
}

// ── EveningSession ────────────────────────────────────────────────────────────

func (s *InventoryService) OpenSession(ctx context.Context, req domain.OpenSessionRequest) (*domain.EveningSession, error) {
	if req.Date == "" {
		req.Date = time.Now().Format("2006-01-02")
	}

	// Kiểm tra không có ca nào đang mở
	existing, err := s.repo.GetOpenSession(ctx)
	if err != nil && err != mongo.ErrNoDocuments {
		return nil, fmt.Errorf("kiểm tra ca hiện tại thất bại: %w", err)
	}
	if existing != nil {
		return nil, fmt.Errorf("đang có ca chưa đóng (ngày %s), vui lòng đóng ca trước", existing.Date)
	}

	// Enrich ingredient names nếu thiếu
	for i, u := range req.Usage {
		if u.IngredientName == "" {
			ing, err := s.repo.GetIngredientTypeByID(ctx, u.IngredientTypeID)
			if err == nil {
				req.Usage[i].IngredientName = ing.Name
			}
		}
	}

	session := &domain.EveningSession{
		ID:       fmt.Sprintf("session_%d", time.Now().UnixNano()),
		Date:     req.Date,
		Status:   domain.SessionOpen,
		Usage:    req.Usage,
		Note:     req.Note,
		OpenedAt: time.Now(),
	}

	if err := s.repo.CreateSession(ctx, session); err != nil {
		return nil, fmt.Errorf("mở ca thất bại: %w", err)
	}
	return session, nil
}

func (s *InventoryService) CloseSession(ctx context.Context, sessionID string, req domain.CloseSessionRequest) (*domain.EveningSession, error) {
	session, err := s.repo.GetSessionByID(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("không tìm thấy ca: %w", err)
	}
	if session.Status == domain.SessionClosed {
		return nil, fmt.Errorf("ca này đã được đóng rồi")
	}

	// Validate: used + wasted không vượt planned
	for _, u := range req.Usage {
		for _, planned := range session.Usage {
			if planned.IngredientTypeID == u.IngredientTypeID {
				if u.UsedBaskets+u.WastedBaskets > planned.PlannedBaskets {
					return nil, fmt.Errorf("%s: tổng dùng (%d) + hỏng (%d) vượt quá chuẩn bị (%d rổ)",
						planned.IngredientName, u.UsedBaskets, u.WastedBaskets, planned.PlannedBaskets)
				}
				break
			}
		}
	}

	// Enrich ingredient names
	for i, u := range req.Usage {
		if u.IngredientName == "" {
			for _, planned := range session.Usage {
				if planned.IngredientTypeID == u.IngredientTypeID {
					req.Usage[i].IngredientName = planned.IngredientName
					req.Usage[i].PlannedBaskets = planned.PlannedBaskets
					break
				}
			}
		}
	}

	note := req.Note
	if note == "" {
		note = session.Note
	}

	if err := s.repo.CloseSession(ctx, sessionID, req.Usage, note); err != nil {
		return nil, err
	}

	return s.repo.GetSessionByID(ctx, sessionID)
}

func (s *InventoryService) GetSessions(ctx context.Context, limit int) ([]domain.EveningSession, error) {
	return s.repo.GetSessions(ctx, limit)
}

func (s *InventoryService) GetCurrentOpenSession(ctx context.Context) (*domain.EveningSession, error) {
	session, err := s.repo.GetOpenSession(ctx)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return session, nil
}

// ── Reports ───────────────────────────────────────────────────────────────────

func (s *InventoryService) GetYieldReport(ctx context.Context, ingredientTypeID string) ([]domain.YieldReportItem, error) {
	batches, err := s.repo.GetBatches(ctx, ingredientTypeID)
	if err != nil {
		return nil, err
	}

	var report []domain.YieldReportItem
	for _, b := range batches {
		report = append(report, domain.YieldReportItem{
			BatchID:        b.ID,
			IngredientName: b.IngredientName,
			InputWeightKg:  b.InputWeightKg,
			OutputBaskets:  b.OutputBaskets,
			WastePercent:   b.WastePercent,
			CostPerBasket:  b.CostPerBasket,
			ProcessedAt:    b.ProcessedAt,
		})
	}
	return report, nil
}

func (s *InventoryService) GetCostReport(ctx context.Context, limit int) ([]domain.CostReportItem, error) {
	sessions, err := s.repo.GetSessions(ctx, limit)
	if err != nil {
		return nil, err
	}

	var report []domain.CostReportItem
	for _, sess := range sessions {
		if sess.Status != domain.SessionClosed {
			continue
		}
		for _, u := range sess.Usage {
			if u.UsedBaskets == 0 {
				continue
			}
			// Lấy chi phí/rổ từ mẻ gần nhất + gia vị
			costPerBasket, _ := s.repo.GetLatestCostPerBasket(ctx, u.IngredientTypeID)
			rc, _ := s.repo.GetRecipeCost(ctx, u.IngredientTypeID)
			if rc != nil {
				costPerBasket += rc.TotalCostPerBasket
			}
			report = append(report, domain.CostReportItem{
				SessionID:      sess.ID,
				Date:           sess.Date,
				IngredientName: u.IngredientName,
				UsedBaskets:    u.UsedBaskets,
				CostPerBasket:  costPerBasket,
				TotalCost:      float64(u.UsedBaskets) * costPerBasket,
			})
		}
	}
	return report, nil
}

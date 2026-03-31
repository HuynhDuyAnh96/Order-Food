// package service

// import (
// 	"context"
// 	"fms_audit/internal/domain"
// 	"fms_audit/internal/infrastructure/repository"
// 	"fmt"
// 	"time"
// )

// type OrderService struct {
// 	orderRepo *repository.OrderRepository
// }

// func NewOrderService(orderRepo *repository.OrderRepository) *OrderService {
// 	return &OrderService{orderRepo: orderRepo}
// }

// func (s *OrderService) CreateOrder(ctx context.Context, request domain.CreateOrderRequest) (*domain.Order, error) {
// 	// Validate table_number (1-20)
// 	if request.TableNumber < 1 || request.TableNumber > 20 {
// 		return nil, fmt.Errorf("invalid table number: must be between 1 and 20")
// 	}

// 	// Validate items
// 	if len(request.Items) == 0 {
// 		return nil, fmt.Errorf("order must have at least one item")
// 	}

// 	// Set default user_id if empty
// 	userID := request.UserID
// 	if userID == "" {
// 		userID = "guest_user"
// 	}

// 	// Generate order ID
// 	orderID := generateOrderID()

// 	// Create order với table number từ frontend
// 	order := &domain.Order{
// 		ID:          orderID,
// 		UserID:      userID,
// 		TableNumber: request.TableNumber, // Nhận table_number từ frontend request
// 		Status:      "pending",           // Có thể dùng request status nếu cần
// 		CreatedAt:   time.Now(),          // Hoặc dùng request.CreatedAt
// 		Items:       make([]domain.OrderItem, 0, len(request.Items)),
// 	}

// 	// Convert items và tính total
// 	var calculatedTotal float64
// 	for _, reqItem := range request.Items {
// 		orderItem := domain.OrderItem{
// 			DishID:   reqItem.ID,
// 			Title:    reqItem.Title,
// 			Quantity: reqItem.Quantity,
// 			Price:    reqItem.Price,
// 		}
// 		order.Items = append(order.Items, orderItem)
// 		calculatedTotal += reqItem.Price * float64(reqItem.Quantity)
// 	}

// 	// Verify total từ frontend để tránh manipulation
// 	if calculatedTotal != request.Total {
// 		return nil, fmt.Errorf("total price mismatch: calculated %.2f, received %.2f", calculatedTotal, request.Total)
// 	}
// 	order.TotalPrice = calculatedTotal

// 	// Lưu order vào database
// 	err := s.orderRepo.CreateOrder(ctx, order)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to save order: %w", err)
// 	}

// 	return order, nil
// }

// // Lấy tất cả orders (có thể filter theo table)
// func (s *OrderService) GetAllOrders(ctx context.Context, tableNumber *int) ([]domain.Order, error) {
// 	if tableNumber != nil {
// 		return s.orderRepo.GetOrdersByTable(ctx, *tableNumber)
// 	}
// 	return s.orderRepo.GetAllOrders(ctx)
// }

// // Lấy orders theo số bàn cụ thể
// func (s *OrderService) GetOrdersByTable(ctx context.Context, tableNumber int) ([]domain.Order, error) {
// 	if tableNumber < 1 || tableNumber > 20 {
// 		return nil, fmt.Errorf("invalid table number: must be between 1 and 20")
// 	}
// 	return s.orderRepo.GetOrdersByTable(ctx, tableNumber)
// }

// func generateOrderID() string {
// 	return fmt.Sprintf("order_%d", time.Now().Unix())
// }

package service

import (
	"context"
	"fms_audit/internal/domain"
	"fms_audit/internal/infrastructure/repository"
	"fmt"
	"time"
)

type OrderService struct {
	orderRepo *repository.OrderRepository
}

func NewOrderService(orderRepo *repository.OrderRepository) *OrderService {
	return &OrderService{orderRepo: orderRepo}
}

func (s *OrderService) CreateOrder(ctx context.Context, request domain.CreateOrderRequest) (*domain.Order, error) {
	// Set default order_type
	if request.OrderType == "" {
		request.OrderType = domain.OrderTypeDineIn
	}

	// Validate table_number chỉ khi dine_in
	if request.OrderType == domain.OrderTypeDineIn {
		if request.TableNumber == 0 {
			request.TableNumber = 1
		}
		if request.TableNumber < 1 || request.TableNumber > 20 {
			return nil, fmt.Errorf("invalid table number: must be between 1 and 20")
		}
	} else {
		request.TableNumber = 0 // takeaway không cần bàn
	}

	// Validate items
	if len(request.Items) == 0 {
		return nil, fmt.Errorf("order must have at least one item")
	}

	// Set default user_id if empty
	userID := request.UserID
	if userID == "" {
		userID = "guest_user"
	}

	// Generate order ID
	orderID := generateOrderID()

	// Create order
	order := &domain.Order{
		ID:          orderID,
		UserID:      userID,
		OrderType:   request.OrderType,
		TableNumber: request.TableNumber,
		Status:      domain.StatusPending,
		CreatedAt:   time.Now(),
		Items:       make([]domain.OrderItem, 0, len(request.Items)),
	}

	// Convert items va tinh total
	var calculatedTotal float64
	for i, reqItem := range request.Items {
		dishID := reqItem.ID
		if reqItem.IsCustom {
			dishID = "" // món lậu không có DishID trong menu
		}
		orderItem := domain.OrderItem{
			ItemID:   fmt.Sprintf("%s_item_%d", orderID, i),
			DishID:   dishID,
			Title:    reqItem.Title,
			Quantity: reqItem.Quantity,
			Price:    reqItem.Price,
			Status:   domain.ItemStatusPending,
			IsCustom: reqItem.IsCustom,
			Note:     reqItem.Note,
		}
		order.Items = append(order.Items, orderItem)
		calculatedTotal += reqItem.Price * float64(reqItem.Quantity)
	}

	// Verify total
	if calculatedTotal != request.Total {
		return nil, fmt.Errorf("total price mismatch: calculated %.2f, received %.2f", calculatedTotal, request.Total)
	}
	order.TotalPrice = calculatedTotal

	// Debug trước khi lưu
	fmt.Printf("Saving order %s with table_number: %d\n", orderID, order.TableNumber)

	// Lưu order
	err := s.orderRepo.CreateOrder(ctx, order)
	if err != nil {
		return nil, fmt.Errorf("failed to save order: %w", err)
	}

	return order, nil
}

func (s *OrderService) GetAllOrders(ctx context.Context, tableNumber *int) ([]domain.Order, error) {
	if tableNumber != nil {
		return s.orderRepo.GetOrdersByTable(ctx, *tableNumber)
	}
	return s.orderRepo.GetAllOrders(ctx)
}

func (s *OrderService) GetOrdersByTable(ctx context.Context, tableNumber int) ([]domain.Order, error) {
	if tableNumber < 1 || tableNumber > 20 {
		return nil, fmt.Errorf("invalid table number: must be between 1 and 20")
	}
	return s.orderRepo.GetOrdersByTable(ctx, tableNumber)
}

func generateOrderID() string {
	return fmt.Sprintf("order_%d", time.Now().Unix())
}

func (s *OrderService) UpdateOrder(ctx context.Context, orderID string, request domain.UpdateOrderRequest) (*domain.Order, error) {
	order, err := s.orderRepo.GetOrderByID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	// Cập nhật items nếu có
	if request.Items != nil {
		order.Items = []domain.OrderItem{} // Reset và thêm mới
		for _, reqItem := range request.Items {
			orderItem := domain.OrderItem{
				DishID:   reqItem.ID,
				Title:    reqItem.Title,
				Quantity: reqItem.Quantity,
				Price:    reqItem.Price,
			}
			order.Items = append(order.Items, orderItem)
		}
		order.TotalPrice = calculateTotalPrice(order.Items) // Tính lại total
	}

	// Cập nhật status nếu có, kiểm tra quy trình workflow
	if request.Status != "" {
		if isValidStatusTransition(order.Status, request.Status) {
			order.Status = request.Status
		} else {
			return nil, fmt.Errorf("invalid status transition from %s to %s", order.Status, request.Status)
		}
	}

	err = s.orderRepo.UpdateOrder(ctx, order)
	if err != nil {
		return nil, err
	}

	return order, nil
}

// PayOrder - nhân viên xác nhận đã thu tiền
// dine_in: phải qua completed trước → paid → bàn trống
// takeaway: có thể pay từ bất kỳ status nào (thu tiền xong đưa túi cho khách là xong)
func (s *OrderService) PayOrder(ctx context.Context, orderID string) (*domain.Order, error) {
	order, err := s.orderRepo.GetOrderByID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("order not found: %w", err)
	}

	if order.OrderType == domain.OrderTypeTakeaway {
		// Takeaway: bếp phải xác nhận xong (ready) mới được thu tiền
		if order.Status == domain.StatusPaid {
			return nil, fmt.Errorf("đơn này đã được thanh toán rồi")
		}
		if order.Status == domain.StatusCancelled {
			return nil, fmt.Errorf("không thể thanh toán đơn đã huỷ")
		}
		if order.Status == domain.StatusPending || order.Status == domain.StatusPreparing {
			return nil, fmt.Errorf("bếp chưa hoàn thành đơn, trạng thái hiện tại: %s", order.Status)
		}
	} else {
		// Dine-in: phải qua completed trước
		if order.Status != domain.StatusCompleted {
			return nil, fmt.Errorf("chỉ có thể thanh toán đơn đã hoàn thành, trạng thái hiện tại: %s", order.Status)
		}
	}

	order.Status = domain.StatusPaid
	order.UpdatedAt = time.Now()
	if err := s.orderRepo.UpdateOrder(ctx, order); err != nil {
		return nil, fmt.Errorf("failed to pay order: %w", err)
	}
	return order, nil
}

func calculateTotalPrice(items []domain.OrderItem) float64 {
	var total float64
	for _, item := range items {
		total += item.Price * float64(item.Quantity)
	}
	return total
}

// Kiểm tra chuyển đổi status có hợp lệ không
// ConfirmOrder - Admin xác nhận đơn hàng (chuyển từ pending sang preparing)
// func (s *OrderService) ConfirmOrder(ctx context.Context, orderID string) (*domain.Order, error) {
// 	// Lấy order hiện tại
// 	order, err := s.orderRepo.GetOrderByID(ctx, orderID)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get order: %w", err)
// 	}

// 	// Kiểm tra trạng thái hiện tại
// 	if order.Status != domain.StatusPending {
// 		return nil, fmt.Errorf("cannot confirm order with status %s, order must be pending", order.Status)
// 	}

// 	// Cập nhật status thành preparing
// 	order.Status = domain.StatusPreparing
// 	order.UpdatedAt = time.Now()

// 	// Lưu vào database
// 	err = s.orderRepo.UpdateOrder(ctx, order)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to update order: %w", err)
// 	}

// 	return order, nil
// }

func (s *OrderService) ConfirmOrder(ctx context.Context, orderID string) (*domain.Order, error) {
	order, err := s.orderRepo.GetOrderByID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	// Chỉ cho phép xác nhận đơn hàng đang pending
	if order.Status != domain.StatusPending {
		return nil, fmt.Errorf("can only confirm orders with status 'pending', current status: %s", order.Status)
	}

	// Kiểm tra chuyển đổi trạng thái
	if !isValidStatusTransition(order.Status, domain.StatusPreparing) {
		return nil, fmt.Errorf("invalid status transition from %s to %s", order.Status, domain.StatusPreparing)
	}

	// Chuyển trạng thái sang preparing
	order.Status = domain.StatusPreparing
	order.UpdatedAt = time.Now() // Thêm thời gian cập nhật

	err = s.orderRepo.UpdateOrder(ctx, order)
	if err != nil {
		return nil, fmt.Errorf("failed to confirm order: %w", err)
	}

	fmt.Printf("Order %s confirmed and moved to preparing status\n", orderID)
	return order, nil
}

func isValidStatusTransition(currentStatus, newStatus domain.OrderStatus) bool {
	validTransitions := map[domain.OrderStatus][]domain.OrderStatus{
		domain.StatusPending: {
			domain.StatusPreparing,
			domain.StatusCancelled,
		},
		domain.StatusPreparing: {
			domain.StatusReady,
			domain.StatusCompleted,
			domain.StatusCancelled,
		},
		domain.StatusReady: {
			domain.StatusCompleted,
			domain.StatusCancelled,
		},
		domain.StatusCompleted: {
			domain.StatusPaid, // nhân viên thu tiền xong
		},
		domain.StatusPaid:      {},
		domain.StatusCancelled: {},
	}

	allowedStatuses, exists := validTransitions[currentStatus]
	if !exists {
		return false
	}

	for _, allowedStatus := range allowedStatuses {
		if allowedStatus == newStatus {
			return true
		}
	}
	return false
}

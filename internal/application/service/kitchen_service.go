package service

import (
	"context"
	"fms_audit/internal/domain"
	"fms_audit/internal/infrastructure/repository"
	"fmt"
	"time"
)

type KitchenService struct {
	orderRepo *repository.OrderRepository
}

func NewKitchenService(orderRepo *repository.OrderRepository) *KitchenService {
	return &KitchenService{orderRepo: orderRepo}
}

// GetKDSBoard tra ve board hien thi cho bep:
//   - Orders sap theo thu tu vao (uu tien bang vao truoc lam truoc)
//   - Moi mon co flag is_duplicate neu ban khac cung goi cung mon do
//   - DishSummary o footer tong hop cac mon trung de bep tranh thu nau chung
func (s *KitchenService) GetKDSBoard(ctx context.Context) (*domain.KDSBoard, error) {
	orders, err := s.orderRepo.GetActiveOrders(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get active orders: %w", err)
	}

	// Xay dung ban do mon trung: dish_id -> danh sach (table, orderID, quantity)
	// Chi tinh mon chua xong (pending hoac cooking)
	type occurrence struct {
		tableNumber int
		orderID     string
		quantity    int
	}
	dishOccurrences := make(map[string][]occurrence)

	for _, order := range orders {
		for _, item := range order.Items {
			if item.Status == domain.ItemStatusReady || item.Status == domain.ItemStatusServed {
				continue
			}
			dishOccurrences[item.DishID] = append(dishOccurrences[item.DishID], occurrence{
				tableNumber: order.TableNumber,
				orderID:     order.ID,
				quantity:    item.Quantity,
			})
		}
	}

	// Build KDS orders theo thu tu uu tien (index = priority)
	kdsOrders := make([]domain.KDSOrder, 0, len(orders))
	for i, order := range orders {
		// Bao ve truong hop order cu trong DB khong co _id hop le
		if order.ID == "" {
			continue
		}
		kdsItems := make([]domain.KDSOrderItem, 0, len(order.Items))

		for _, item := range order.Items {
			occs := dishOccurrences[item.DishID]
			isDuplicate := len(occs) > 1

			var dupInfo []domain.DuplicateInfo
			if isDuplicate {
				for _, occ := range occs {
					if occ.orderID == order.ID {
						continue
					}
					dupInfo = append(dupInfo, domain.DuplicateInfo{
						TableNumber: occ.tableNumber,
						OrderID:     occ.orderID,
						Quantity:    occ.quantity,
					})
				}
			}

			kdsItems = append(kdsItems, domain.KDSOrderItem{
				ItemID:        item.ItemID,
				DishID:        item.DishID,
				Title:         item.Title,
				Quantity:      item.Quantity,
				Price:         item.Price,
				Status:        item.Status,
				IsDuplicate:   isDuplicate && len(dupInfo) > 0,
				DuplicateInfo: dupInfo,
			})
		}

		kdsOrders = append(kdsOrders, domain.KDSOrder{
			Priority:    i + 1,
			OrderID:     order.ID,
			OrderType:   order.OrderType,
			TableNumber: order.TableNumber,
			CreatedAt:   order.CreatedAt,
			WaitMinutes: int(time.Since(order.CreatedAt).Minutes()),
			Status:      order.Status,
			Items:       kdsItems,
		})
	}

	// Build dish summary: chi nhung mon co >= 2 ban cung goi
	dishSummary := make([]domain.DishSummary, 0)
	for dishID, occs := range dishOccurrences {
		if len(occs) < 2 {
			continue
		}
		totalQty := 0
		tables := make([]domain.DuplicateInfo, 0, len(occs))
		var title string
		for _, occ := range occs {
			totalQty += occ.quantity
			tables = append(tables, domain.DuplicateInfo{
				TableNumber: occ.tableNumber,
				OrderID:     occ.orderID,
				Quantity:    occ.quantity,
			})
		}
		// Tim ten mon tu orders
		for _, order := range orders {
			for _, item := range order.Items {
				if item.DishID == dishID {
					title = item.Title
					break
				}
			}
			if title != "" {
				break
			}
		}
		dishSummary = append(dishSummary, domain.DishSummary{
			DishID:   dishID,
			Title:    title,
			TotalQty: totalQty,
			Tables:   tables,
		})
	}

	return &domain.KDSBoard{
		Orders:      kdsOrders,
		DishSummary: dishSummary,
	}, nil
}

// CompleteOrder - bep hoan thanh 1 ban, chuyen thang sang completed
func (s *KitchenService) CompleteOrder(ctx context.Context, orderID string) (*domain.Order, error) {
	order, err := s.orderRepo.GetOrderByID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("order not found: %w", err)
	}
	if order.Status == domain.StatusCompleted || order.Status == domain.StatusCancelled {
		return nil, fmt.Errorf("order is already %s", order.Status)
	}
	order.Status = domain.StatusCompleted
	order.UpdatedAt = time.Now()
	if err := s.orderRepo.UpdateOrder(ctx, order); err != nil {
		return nil, fmt.Errorf("failed to complete order: %w", err)
	}
	return order, nil
}

// UpdateItemStatus - bep cap nhat trang thai 1 mon (cooking / ready / served)
// Neu tat ca mon trong don da ready/served -> tu dong chuyen order sang ready
func (s *KitchenService) UpdateItemStatus(ctx context.Context, orderID string, itemID string, status domain.ItemStatus) (*domain.Order, error) {
	validStatuses := map[domain.ItemStatus]bool{
		domain.ItemStatusCooking: true,
		domain.ItemStatusReady:   true,
		domain.ItemStatusServed:  true,
	}
	if !validStatuses[status] {
		return nil, fmt.Errorf("invalid item status: %s", status)
	}

	if err := s.orderRepo.UpdateOrderItemStatus(ctx, orderID, itemID, status); err != nil {
		return nil, err
	}

	order, err := s.orderRepo.GetOrderByID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	// Neu tat ca mon da ready hoac served -> chuyen order sang ready
	allDone := true
	for _, item := range order.Items {
		if item.Status != domain.ItemStatusReady && item.Status != domain.ItemStatusServed {
			allDone = false
			break
		}
	}
	if allDone && (order.Status == domain.StatusPreparing || order.Status == domain.StatusPending) {
		order.Status = domain.StatusReady
		order.UpdatedAt = time.Now()
		if err := s.orderRepo.UpdateOrder(ctx, order); err != nil {
			return nil, fmt.Errorf("failed to update order status: %w", err)
		}
	}

	return order, nil
}

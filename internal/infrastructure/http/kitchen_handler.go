package http

import (
	"encoding/json"
	"fms_audit/internal/application/service"
	"fms_audit/internal/domain"
	"net/http"

	"github.com/gin-gonic/gin"
)

type KitchenHandler struct {
	kitchenService *service.KitchenService
	hub            *Hub
}

func NewKitchenHandler(kitchenService *service.KitchenService, hub *Hub) *KitchenHandler {
	return &KitchenHandler{
		kitchenService: kitchenService,
		hub:            hub,
	}
}

// GetKDSBoard GET /api/kitchen/board
// Tra ve danh sach don sap theo thu tu vao, kem flag mon trung giua cac ban
func (h *KitchenHandler) GetKDSBoard(c *gin.Context) {
	board, err := h.kitchenService.GetKDSBoard(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": board})
}

// StartCooking POST /api/kitchen/orders/:orderId/items/:itemId/start
// Bep bat dau nau 1 mon -> item status: cooking
func (h *KitchenHandler) StartCooking(c *gin.Context) {
	orderID := c.Param("orderId")
	itemID := c.Param("itemId")

	order, err := h.kitchenService.UpdateItemStatus(c.Request.Context(), orderID, itemID, domain.ItemStatusCooking)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": order})

	// Broadcast cho staff app biet mon dang duoc nau
	msg, _ := json.Marshal(map[string]interface{}{
		"event":       "dish_cooking",
		"orderID":     orderID,
		"itemID":      itemID,
		"orderType":   order.OrderType,
		"tableNumber": order.TableNumber,
		"orderStatus": order.Status,
	})
	h.hub.Send(msg)
}

// MarkDishReady POST /api/kitchen/orders/:orderId/items/:itemId/ready
// Bep hoan thanh mon -> item status: ready, alert staff lay mon
func (h *KitchenHandler) MarkDishReady(c *gin.Context) {
	orderID := c.Param("orderId")
	itemID := c.Param("itemId")

	order, err := h.kitchenService.UpdateItemStatus(c.Request.Context(), orderID, itemID, domain.ItemStatusReady)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": order})

	// Tim ten mon vua xong de broadcast
	itemTitle := ""
	for _, item := range order.Items {
		if item.ItemID == itemID {
			itemTitle = item.Title
			break
		}
	}

	// Broadcast: staff app nhan alert "mon nay xong roi, mang ra ban / goi khach"
	msg, _ := json.Marshal(map[string]interface{}{
		"event":       "dish_ready",
		"orderID":     orderID,
		"itemID":      itemID,
		"itemTitle":   itemTitle,
		"orderType":   order.OrderType,
		"tableNumber": order.TableNumber,
		"orderStatus": order.Status,
	})
	h.hub.Send(msg)
}

// CompleteOrder POST /api/kitchen/orders/:orderId/complete
// Bep hoan thanh ca ban, order chuyen sang completed va bien mat khoi board
func (h *KitchenHandler) CompleteOrder(c *gin.Context) {
	orderID := c.Param("orderId")

	order, err := h.kitchenService.CompleteOrder(c.Request.Context(), orderID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": order})

	msg, _ := json.Marshal(map[string]interface{}{
		"event":       "order_completed",
		"orderID":     orderID,
		"orderType":   order.OrderType,
		"tableNumber": order.TableNumber,
		"status":      order.Status,
	})
	h.hub.Send(msg)
}

// MarkDishServed POST /api/kitchen/orders/:orderId/items/:itemId/served
// Nhan vien xac nhan da mang mon ra ban -> item status: served
func (h *KitchenHandler) MarkDishServed(c *gin.Context) {
	orderID := c.Param("orderId")
	itemID := c.Param("itemId")

	order, err := h.kitchenService.UpdateItemStatus(c.Request.Context(), orderID, itemID, domain.ItemStatusServed)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": order})

	msg, _ := json.Marshal(map[string]interface{}{
		"event":       "dish_served",
		"orderID":     orderID,
		"itemID":      itemID,
		"tableNumber": order.TableNumber,
		"orderStatus": order.Status,
	})
	h.hub.Send(msg)
}

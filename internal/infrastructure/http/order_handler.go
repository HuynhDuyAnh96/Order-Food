package http

import (
	"encoding/json"
	"fms_audit/internal/application/service"
	"fms_audit/internal/domain"
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type OrderHandler struct {
	orderService *service.OrderService
	broadcast    chan []byte
	clients      map[*websocket.Conn]bool
	mutex        sync.RWMutex
}

var upgrade = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func NewOrderHandler(orderService *service.OrderService) *OrderHandler {
	h := &OrderHandler{
		orderService: orderService,
		broadcast:    make(chan []byte),
		clients:      make(map[*websocket.Conn]bool),
	}

	// Start message handler goroutine
	go h.handleMessages()

	return h
}

// func (h *OrderHandler) CreateOrder(c *gin.Context) {
// 	var request domain.CreateOrderRequest
// 	if err := c.ShouldBindJSON(&request); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	order, err := h.orderService.CreateOrder(c.Request.Context(), request)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"success": true, "data": order})

// 	orderJson, _ := json.Marshal(order)
// 	h.broadcast <- orderJson
// }

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var request domain.CreateOrderRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order, err := h.orderService.CreateOrder(c.Request.Context(), request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": order})

	// Broadcast new order
	message := map[string]interface{}{
		"event":       "new_order",
		"orderID":     order.ID,
		"status":      order.Status,
		"tableNumber": order.TableNumber,
		"items":       order.Items,
	}
	orderJSON, _ := json.Marshal(message)
	h.broadcast <- orderJSON
}

func (h *OrderHandler) GetAllOrders(c *gin.Context) {
	tableStr := c.Query("table_number") // Filter optional qua query param
	var tableNumber *int
	if tableStr != "" {
		t, err := strconv.Atoi(tableStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid table_number"})
			return
		}
		tableNumber = &t
	}

	orders, err := h.orderService.GetAllOrders(c.Request.Context(), tableNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": orders})
}

// WebSocket endpoint cho bếp
func (h *OrderHandler) WsEndpoint(c *gin.Context) {
	conn, err := upgrade.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upgrade connection"})
		return
	}
	defer conn.Close()

	// Thêm client vào danh sách với mutex
	h.mutex.Lock()
	h.clients[conn] = true
	h.mutex.Unlock()

	// Lắng nghe tin nhắn từ client (nếu cần)
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			// Client ngắt kết nối
			h.mutex.Lock()
			delete(h.clients, conn)
			h.mutex.Unlock()
			break
		}
	}
}

// Xử lý broadcast messages đến tất cả clients
func (h *OrderHandler) handleMessages() {
	for {
		message := <-h.broadcast

		// Gửi message đến tất cả clients đã kết nối
		h.mutex.RLock()
		for client := range h.clients {
			err := client.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				// Client ngắt kết nối, xóa khỏi danh sách
				client.Close()
				h.mutex.RUnlock()
				h.mutex.Lock()
				delete(h.clients, client)
				h.mutex.Unlock()
				h.mutex.RLock()
			}
		}
		h.mutex.RUnlock()
	}
}

// ConfirmOrder - Admin xác nhận đơn hàng

// func (h *OrderHandler) ConfirmOrder(c *gin.Context) {
// 	orderID := c.Param("orderId")

// 	confirmedOrder, err := h.orderService.ConfirmOrder(c.Request.Context(), orderID)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"success": true, "data": confirmedOrder})

// 	// Broadcast đến kitchen via WebSocket
// 	confirmedOrderJSON, _ := json.Marshal(confirmedOrder)
// 	h.broadcast <- confirmedOrderJSON
// }

func (h *OrderHandler) ConfirmOrder(c *gin.Context) {
	orderID := c.Param("orderId")

	confirmedOrder, err := h.orderService.ConfirmOrder(c.Request.Context(), orderID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": confirmedOrder})

	// Broadcast đến kitchen via WebSocket
	message := map[string]interface{}{
		"event":       "order_confirmed",
		"orderID":     confirmedOrder.ID,
		"status":      confirmedOrder.Status,
		"tableNumber": confirmedOrder.TableNumber,
	}
	confirmedOrderJSON, _ := json.Marshal(message)
	h.broadcast <- confirmedOrderJSON
}

// func (h *OrderHandler) UpdateOrder(c *gin.Context) {
// 	orderID := c.Param("orderId")
// 	var request domain.UpdateOrderRequest
// 	if err := c.ShouldBindJSON(&request); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	updatedOrder, err := h.orderService.UpdateOrder(c.Request.Context(), orderID, request)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"success": true, "data": updatedOrder})
// 	updatedOrderJSON, _ := json.Marshal(updatedOrder)
// 	h.broadcast <- updatedOrderJSON
// }

func (h *OrderHandler) UpdateOrder(c *gin.Context) {
	orderID := c.Param("orderId")
	var request domain.UpdateOrderRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedOrder, err := h.orderService.UpdateOrder(c.Request.Context(), orderID, request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": updatedOrder})

	// Broadcast updated order
	message := map[string]interface{}{
		"event":       "order_updated",
		"orderID":     updatedOrder.ID,
		"status":      updatedOrder.Status,
		"tableNumber": updatedOrder.TableNumber,
		"items":       updatedOrder.Items,
	}
	updatedOrderJSON, _ := json.Marshal(message)
	h.broadcast <- updatedOrderJSON
}

// GetOrdersByTable - Endpoint riêng cho get by table
func (h *OrderHandler) GetOrdersByTable(c *gin.Context) {
	tableStr := c.Param("table_number")
	tableNumber, err := strconv.Atoi(tableStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid table_number"})
		return
	}

	orders, err := h.orderService.GetOrdersByTable(c.Request.Context(), tableNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": orders})
}

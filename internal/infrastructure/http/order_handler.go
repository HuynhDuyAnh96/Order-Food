package http

import (
	"encoding/json"
	"fms_audit/internal/application/service"
	"fms_audit/internal/domain"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type OrderHandler struct {
	orderService *service.OrderService
	hub          *Hub
}

var upgrade = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func NewOrderHandler(orderService *service.OrderService, hub *Hub) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
		hub:          hub,
	}
}

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

	msg, _ := json.Marshal(map[string]interface{}{
		"event":       "new_order",
		"orderID":     order.ID,
		"status":      order.Status,
		"tableNumber": order.TableNumber,
		"items":       order.Items,
	})
	h.hub.Send(msg)
}

func (h *OrderHandler) GetAllOrders(c *gin.Context) {
	tableStr := c.Query("table_number")
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

func (h *OrderHandler) ConfirmOrder(c *gin.Context) {
	orderID := c.Param("orderId")

	confirmedOrder, err := h.orderService.ConfirmOrder(c.Request.Context(), orderID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": confirmedOrder})

	msg, _ := json.Marshal(map[string]interface{}{
		"event":       "order_confirmed",
		"orderID":     confirmedOrder.ID,
		"status":      confirmedOrder.Status,
		"tableNumber": confirmedOrder.TableNumber,
	})
	h.hub.Send(msg)
}

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

	msg, _ := json.Marshal(map[string]interface{}{
		"event":       "order_updated",
		"orderID":     updatedOrder.ID,
		"status":      updatedOrder.Status,
		"tableNumber": updatedOrder.TableNumber,
		"items":       updatedOrder.Items,
	})
	h.hub.Send(msg)
}

// PayOrder POST /api/orders/:orderId/pay
// Nhan vien xac nhan da thu tien → ban trong
func (h *OrderHandler) PayOrder(c *gin.Context) {
	orderID := c.Param("orderId")

	paidOrder, err := h.orderService.PayOrder(c.Request.Context(), orderID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": paidOrder})

	msg, _ := json.Marshal(map[string]interface{}{
		"event":       "order_paid",
		"orderID":     paidOrder.ID,
		"tableNumber": paidOrder.TableNumber,
		"status":      paidOrder.Status,
	})
	h.hub.Send(msg)
}

// WsEndpoint - WebSocket endpoint dung chung cho bep va staff
func (h *OrderHandler) WsEndpoint(c *gin.Context) {
	conn, err := upgrade.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upgrade connection"})
		return
	}
	defer func() {
		h.hub.Unregister(conn)
		conn.Close()
	}()

	h.hub.Register(conn)

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}
}

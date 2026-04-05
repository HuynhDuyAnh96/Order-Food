package http

import (
	"fms_audit/internal/application/service"
	"fms_audit/internal/domain"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type InventoryHandler struct {
	inventoryService *service.InventoryService
}

func NewInventoryHandler(inventoryService *service.InventoryService) *InventoryHandler {
	return &InventoryHandler{inventoryService: inventoryService}
}

// ── IngredientType ────────────────────────────────────────────────────────────

// POST /api/inventory/ingredient-types
func (h *InventoryHandler) CreateIngredientType(c *gin.Context) {
	var req domain.CreateIngredientTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	ing, err := h.inventoryService.CreateIngredientType(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": ing})
}

// GET /api/inventory/ingredient-types
func (h *InventoryHandler) GetIngredientTypes(c *gin.Context) {
	types, err := h.inventoryService.GetIngredientTypes(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": types})
}

// ── InventoryReceipt ──────────────────────────────────────────────────────────

// POST /api/inventory/receipts
func (h *InventoryHandler) CreateReceipt(c *gin.Context) {
	var req domain.CreateReceiptRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	receipt, err := h.inventoryService.CreateReceipt(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": receipt})
}

// GET /api/inventory/receipts?ingredient_type_id=xxx
func (h *InventoryHandler) GetReceipts(c *gin.Context) {
	ingredientTypeID := c.Query("ingredient_type_id")
	receipts, err := h.inventoryService.GetReceipts(c.Request.Context(), ingredientTypeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": receipts})
}

// ── ProcessingBatch ───────────────────────────────────────────────────────────

// POST /api/inventory/batches
func (h *InventoryHandler) CreateProcessingBatch(c *gin.Context) {
	var req domain.CreateProcessingBatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	batch, err := h.inventoryService.CreateProcessingBatch(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": batch})
}

// GET /api/inventory/batches?ingredient_type_id=xxx
func (h *InventoryHandler) GetBatches(c *gin.Context) {
	ingredientTypeID := c.Query("ingredient_type_id")
	batches, err := h.inventoryService.GetBatches(c.Request.Context(), ingredientTypeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": batches})
}

// ── RecipeCost ────────────────────────────────────────────────────────────────

// PUT /api/inventory/recipe-costs
func (h *InventoryHandler) UpsertRecipeCost(c *gin.Context) {
	var req domain.UpsertRecipeCostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	rc, err := h.inventoryService.UpsertRecipeCost(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": rc})
}

// GET /api/inventory/recipe-costs/:ingredientTypeId
func (h *InventoryHandler) GetRecipeCost(c *gin.Context) {
	ingredientTypeID := c.Param("ingredientTypeId")
	rc, err := h.inventoryService.GetRecipeCost(c.Request.Context(), ingredientTypeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}
	if rc == nil {
		c.JSON(http.StatusOK, gin.H{"success": true, "data": nil, "message": "chưa có chi phí gia vị cho loại này"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": rc})
}

// ── Stock Summary ─────────────────────────────────────────────────────────────

// GET /api/inventory/stock
func (h *InventoryHandler) GetStockSummary(c *gin.Context) {
	summary, err := h.inventoryService.GetStockSummary(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": summary})
}

// ── EveningSession ────────────────────────────────────────────────────────────

// POST /api/inventory/sessions/open
func (h *InventoryHandler) OpenSession(c *gin.Context) {
	var req domain.OpenSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	session, err := h.inventoryService.OpenSession(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": session})
}

// POST /api/inventory/sessions/:sessionId/close
func (h *InventoryHandler) CloseSession(c *gin.Context) {
	sessionID := c.Param("sessionId")
	var req domain.CloseSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	session, err := h.inventoryService.CloseSession(c.Request.Context(), sessionID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": session})
}

// GET /api/inventory/sessions?limit=30
func (h *InventoryHandler) GetSessions(c *gin.Context) {
	limit := 30
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	sessions, err := h.inventoryService.GetSessions(c.Request.Context(), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": sessions})
}

// GET /api/inventory/sessions/current
func (h *InventoryHandler) GetCurrentSession(c *gin.Context) {
	session, err := h.inventoryService.GetCurrentOpenSession(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}
	if session == nil {
		c.JSON(http.StatusOK, gin.H{"success": true, "data": nil, "message": "không có ca nào đang mở"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": session})
}

// ── Reports ───────────────────────────────────────────────────────────────────

// GET /api/inventory/reports/yield?ingredient_type_id=xxx
func (h *InventoryHandler) GetYieldReport(c *gin.Context) {
	ingredientTypeID := c.Query("ingredient_type_id")
	report, err := h.inventoryService.GetYieldReport(c.Request.Context(), ingredientTypeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": report})
}

// GET /api/inventory/reports/cost?limit=30
func (h *InventoryHandler) GetCostReport(c *gin.Context) {
	limit := 30
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	report, err := h.inventoryService.GetCostReport(c.Request.Context(), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": report})
}

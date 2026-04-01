package http

import (
	"fms_audit/internal/application/service"
	"fms_audit/internal/domain"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DishHandler struct {
	dishService *service.DishService
}

func NewDishHandler(dishService *service.DishService) *DishHandler {
	return &DishHandler{dishService: dishService}
}

func (h *DishHandler) GetDishes(c *gin.Context) {
	var filter domain.DishFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid query parameters"})
		return
	}

	response, err := h.dishService.GetDishes(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get dishes"})
		return
	}
	c.JSON(http.StatusOK, response)
}

func (h *DishHandler) GetFeaturedDishes(c *gin.Context) {
	dishes, err := h.dishService.GetFeaturedDishes(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get featured dishes"})
		return
	}
	c.JSON(http.StatusOK, dishes)
}

func (h *DishHandler) GetStirFriedDishes(c *gin.Context) {
	dishes, err := h.dishService.GetDishesByStirFried(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get stir-fried dishes"})
		return
	}
	c.JSON(http.StatusOK, dishes)
}

func (h *DishHandler) GetSteamedDishes(c *gin.Context) {
	dishes, err := h.dishService.GetDishesBySteamed(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get steamed dishes"})
		return
	}
	c.JSON(http.StatusOK, dishes)
}

func (h *DishHandler) GetGrilledDishes(c *gin.Context) {
	dishes, err := h.dishService.GetDishesByGrilled(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get grilled dishes"})
		return
	}
	c.JSON(http.StatusOK, dishes)
}

func (h *DishHandler) GetDrinks(c *gin.Context) {
	drinks, err := h.dishService.GetDrinks(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get drinks"})
		return
	}
	c.JSON(http.StatusOK, drinks)
}

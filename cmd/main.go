package main

import (
	"context"
	"fms_audit/internal/application/service"
	httpinfra "fms_audit/internal/infrastructure/http"
	"fms_audit/internal/infrastructure/repository"
	"log"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Kết nối MongoDB
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(context.Background())

	// Kiểm tra kết nối
	if err := client.Ping(context.Background(), nil); err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}

	// Khởi tạo repository với MongoDB database
	db := client.Database("fms_audit")
	dishRepo := repository.NewDishRepository(db)
	orderRepo := repository.NewOrderRepository(db)

	// Khởi tạo service
	dishService := service.NewDishService(dishRepo)
	orderService := service.NewOrderService(orderRepo)

	// Khởi tạo handler
	orderHandler := httpinfra.NewOrderHandler(orderService)
	dishHandler := httpinfra.NewDishHandler(dishService)

	// Initialize Gin router
	r := gin.Default()

	// Add CORS middleware
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Serve static files
	r.Static("/images", "./internal/img")

	// API routes
	api := r.Group("/api")
	{
		api.GET("/dishes", dishHandler.GetDishes)
		api.GET("/dishes/featured", dishHandler.GetFeaturedDishes)
		api.GET("/dishes/stir-fried", dishHandler.GetStirFriedDishes)
		api.GET("/dishes/steamed", dishHandler.GetSteamedDishes)
		api.GET("/dishes/grilled", dishHandler.GetGrilledDishes)
		api.POST("/orders", orderHandler.CreateOrder)
		api.GET("/orders", orderHandler.GetAllOrders)
		api.GET("/orders/table/:table_number", orderHandler.GetOrdersByTable)
		api.PUT("/orders/:orderId", orderHandler.UpdateOrder)
		api.POST("/orders/:orderId/confirm", orderHandler.ConfirmOrder) // Admin xác nhận đơn hàng
	}
	api.GET("/ws", orderHandler.WsEndpoint)

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Dish API is running",
		})
	})

	log.Println("HTTP server listening on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

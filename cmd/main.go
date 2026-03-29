package main

import (
	"context"
	"fms_audit/internal/application/service"
	httpinfra "fms_audit/internal/infrastructure/http"
	"fms_audit/internal/infrastructure/repository"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}
	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(context.Background())

	if err := client.Ping(context.Background(), nil); err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}

	db := client.Database("fms_audit")

	// Repositories
	dishRepo := repository.NewDishRepository(db)
	orderRepo := repository.NewOrderRepository(db)

	// Services
	dishService := service.NewDishService(dishRepo)
	orderService := service.NewOrderService(orderRepo)
	kitchenService := service.NewKitchenService(orderRepo)

	// WebSocket hub dung chung
	hub := httpinfra.NewHub()

	// Handlers
	orderHandler := httpinfra.NewOrderHandler(orderService, hub)
	dishHandler := httpinfra.NewDishHandler(dishService)
	kitchenHandler := httpinfra.NewKitchenHandler(kitchenService, hub)

	r := gin.Default()

	// CORS middleware
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

	r.Static("/images", "./internal/img")

	api := r.Group("/api")
	{
		// Dish routes
		api.GET("/dishes", dishHandler.GetDishes)
		api.GET("/dishes/featured", dishHandler.GetFeaturedDishes)
		api.GET("/dishes/stir-fried", dishHandler.GetStirFriedDishes)
		api.GET("/dishes/steamed", dishHandler.GetSteamedDishes)
		api.GET("/dishes/grilled", dishHandler.GetGrilledDishes)

		// Order routes
		api.POST("/orders", orderHandler.CreateOrder)
		api.GET("/orders", orderHandler.GetAllOrders)
		api.GET("/orders/table/:table_number", orderHandler.GetOrdersByTable)
		api.PUT("/orders/:orderId", orderHandler.UpdateOrder)
		api.POST("/orders/:orderId/confirm", orderHandler.ConfirmOrder)
		api.POST("/orders/:orderId/pay", orderHandler.PayOrder)

		// Kitchen routes
		api.GET("/kitchen/board", kitchenHandler.GetKDSBoard)
		api.POST("/kitchen/orders/:orderId/complete", kitchenHandler.CompleteOrder)
		api.POST("/kitchen/orders/:orderId/items/:itemId/start", kitchenHandler.StartCooking)
		api.POST("/kitchen/orders/:orderId/items/:itemId/ready", kitchenHandler.MarkDishReady)
		api.POST("/kitchen/orders/:orderId/items/:itemId/served", kitchenHandler.MarkDishServed)
	}

	// WebSocket endpoint
	api.GET("/ws", orderHandler.WsEndpoint)

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "message": "Dish API is running"})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Println("HTTP server listening on 0.0.0.0:" + port)
	if err := r.Run("0.0.0.0:" + port); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

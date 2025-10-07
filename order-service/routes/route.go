package routes

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"

	"order-service/common/kafka"
	"order-service/controllers"
	"order-service/repositories"
	"order-service/services"
)

func SetRoutes(router *gin.Engine, db *gorm.DB, redisClient *redis.Client, kafkaProducer *kafka.KafkaProducer) {
	// Buat instance repository
	orderRepo := repositories.NewOrderRepository(db)

	productURL := os.Getenv("PRODUCT_SERVICE_URL")

	orderService := services.NewOrderService(
		orderRepo,
		redisClient,
		kafkaProducer,
		productURL,
		http.DefaultClient,
	)

	orderController := controllers.NewOrderController(orderService)

	api := router.Group("/api")
	{
		orders := api.Group("/orders")
		{
			orders.POST("", orderController.CreateOrder)
			orders.GET("/:id", orderController.GetOrderByID)
			orders.GET("/product/:productId", orderController.GetOrdersByProductID)
		}
	}
}

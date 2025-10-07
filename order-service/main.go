package main

import (
	"log"
	"order-service/common/kafka"
	"order-service/common/redis"
	"order-service/config"
	"order-service/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Setup dependencies
	config.ConnectDatabase()
	redis.ConnectRedis()
	kafkaProducer := kafka.NewKafkaProducer()

	// Start Kafka consumer in background
	go kafka.StartOrderCreatedConsumer()

	// Setup Gin router
	router := gin.Default()

	// Set routes - kirim redis.Client bukan config
	routes.SetRoutes(router, config.DB, redis.Client, kafkaProducer)

	log.Println("🚀 Order service running on :3001")
	router.Run(":3001")
}

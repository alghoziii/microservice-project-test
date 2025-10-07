package kafka

import (
	"context"
	"encoding/json"
	"log"

	"github.com/segmentio/kafka-go"
)

func StartOrderCreatedConsumer() {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{"localhost:9092"},
		Topic:    "order.created",
		GroupID:  "order-service-group",
		MinBytes: 10e3,
		MaxBytes: 10e6,
	})

	log.Println("🔄 Kafka consumer started for order.created events")

	for {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Printf("Kafka consumer error: %v", err)
			continue
		}

		var event map[string]interface{}
		json.Unmarshal(msg.Value, &event)

		log.Printf("📨 KAFKA EVENT RECEIVED: order.created - OrderID: %v, ProductID: %v, Quantity: %v",
			event["orderId"], event["productId"], event["quantity"])
	}
}

package kafka

import (
	"context"
	"encoding/json"
	"log"

	"github.com/segmentio/kafka-go"
)

type KafkaProducer struct {
	writer *kafka.Writer
}

func NewKafkaProducer() *KafkaProducer {
	writer := &kafka.Writer{
		Addr:     kafka.TCP("localhost:9092"),
		Topic:    "order.created",
		Balancer: &kafka.LeastBytes{},
	}

	return &KafkaProducer{writer: writer}
}

func (k *KafkaProducer) SendOrderCreated(orderID uint, productID int, quantity int) {
	event := map[string]interface{}{
		"orderId":   orderID,
		"productId": productID,
		"quantity":  quantity,
	}

	message, _ := json.Marshal(event)

	err := k.writer.WriteMessages(context.Background(),
		kafka.Message{
			Value: message,
		},
	)

	if err != nil {
		log.Printf("Failed to send Kafka event: %v", err)
	} else {
		log.Printf("✅ KAFKA EVENT SENT: order.created - OrderID: %d", orderID)
	}
}

func (k *KafkaProducer) Close() {
	k.writer.Close()
}

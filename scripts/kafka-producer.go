package main

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	log "github.com/sirupsen/logrus"
	"os"
)

func main() {
	producer, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost:9092"})
	if err != nil {
		log.Fatalf("Failed to create producer: %s", err)
	}
	defer producer.Close()

	topic := "orders"
	messageByte, err := os.ReadFile("./scripts/input.json")
	if err != nil {
		log.Fatalf("Failed to read input file: %s", err)
	}
	message := string(messageByte)

	err = producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          []byte(message),
	}, nil)
	if err != nil {
		log.Fatalf("Failed to send message: %s", err)
	}

	producer.Flush(15 * 100)
	log.Println("Message sent successfully!")
}

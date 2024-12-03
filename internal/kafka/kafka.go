package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	log "github.com/sirupsen/logrus"
	"time"
	"wb-l-zero/internal/entity"
	"wb-l-zero/internal/service"
	"wb-l-zero/pkg/validator"
)

type Consumer struct {
	consumer  *kafka.Consumer
	topic     string
	service   *service.Service
	validator *validator.CustomValidator
	logger    *log.Logger
}

func NewConsumer(brokers, groupID, topic string, service *service.Service, valid *validator.CustomValidator, logger *log.Logger) (*Consumer, error) {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": brokers,
		"group.id":          groupID,
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		return nil, fmt.Errorf("error while creating new consumer: %w", err)
	}

	return &Consumer{
		consumer:  consumer,
		topic:     topic,
		service:   service,
		validator: valid,
		logger:    logger,
	}, nil
}

func (kc *Consumer) Start(ctx context.Context) error {
	err := kc.consumer.Subscribe(kc.topic, nil)
	if err != nil {
		return fmt.Errorf("failed to subscribe to topic: %w", err)
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			msg, err := kc.consumer.ReadMessage(100 * time.Millisecond)
			if err != nil {
				if kafkaError, ok := err.(kafka.Error); ok && kafkaError.Code() == kafka.ErrTimedOut {
					continue
				}
				return fmt.Errorf("consumer error: %w", err)
			}

			var order entity.Order
			if err := json.Unmarshal(msg.Value, &order); err != nil {
				return fmt.Errorf("failed to unmarshal message: %w", err)
			}

			if err := kc.validator.Validate(order); err != nil {
				return fmt.Errorf("order validation error: %w", err)
			}

			err = kc.service.Create(ctx, &order)
			if err != nil {
				return err
			}

			kc.logger.Infof("order with order_uid=%s added", order.OrderUID)
		}
	}
}

func (kc *Consumer) Close() {
	err := kc.consumer.Close()
	if err != nil {
		log.Printf("Error closing consumer: %v", err)
	}
}

package kafka

import (
	"Payment_Service/internal/models"
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/IBM/sarama"
)

type Consumer struct {
	brokers []string
	topics  []string
	handler PaymentHandler
}

type PaymentHandler interface {
	HandleOrderEvent(event models.OrderEvent) error
}

func NewConsumer(brokers []string, topics []string, handler PaymentHandler) *Consumer {
	return &Consumer{
		brokers: brokers,
		topics:  topics,
		handler: handler,
	}
}

func (c *Consumer) Start(ctx context.Context) error {
	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	config.Consumer.Return.Errors = true

	consumerGroup, err := sarama.NewConsumerGroup(c.brokers, "payment-service", config)
	if err != nil {
		return err
	}

	defer func() {
		if err := consumerGroup.Close(); err != nil {
			log.Printf("Error closing consumer group: %v", err)
		}
	}()

	go func() {
		for err := range consumerGroup.Errors() {
			log.Printf("Consumer error: %v", err)
		}
	}()

	consumer := &consumerGroupHandler{handler: c.handler}

	for {
		select {
		case <-ctx.Done():
			log.Println("Terminating: context cancelled")
			return nil
		default:
			if err := consumerGroup.Consume(ctx, c.topics, consumer); err != nil {
				log.Printf("Error from consumer: %v", err)
				time.Sleep(5 * time.Second)
			}
		}
	}
}

type consumerGroupHandler struct {
	handler PaymentHandler
}

func (h *consumerGroupHandler) Setup(sarama.ConsumerGroupSession) error   { return nil }
func (h *consumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error { return nil }

func (h *consumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message := <-claim.Messages():
			if message == nil {
				return nil
			}

			log.Printf("Received message from topic %s: %s", message.Topic, string(message.Value))

			// Сначала пытаемся определить тип события по содержимому
			var messageData map[string]interface{}
			if err := json.Unmarshal(message.Value, &messageData); err != nil {
				log.Printf("Error unmarshaling message: %v", err)
				session.MarkMessage(message, "")
				continue
			}

			eventType, ok := messageData["event_type"].(string)
			if !ok {
				log.Printf("No event_type found in message")
				session.MarkMessage(message, "")
				continue
			}

			// Определяем тип события и обрабатываем соответственно
			switch eventType {
			case "order_created", "order_status_updated":
				// Это OrderEvent
				var orderEvent models.OrderEvent
				if err := json.Unmarshal(message.Value, &orderEvent); err == nil {
					if err := h.handler.HandleOrderEvent(orderEvent); err != nil {
						log.Printf("Error handling order event: %v", err)
					}
				} else {
					log.Printf("Error unmarshaling OrderEvent: %v", err)
				}
			case "payment_required", "payment_completed":
				// Это PaymentEvent - Payment Service не обрабатывает свои собственные события
				log.Printf("Ignoring own payment event: %s", eventType)
			default:
				log.Printf("Unknown event type: %s", eventType)
			}

			session.MarkMessage(message, "")

		case <-session.Context().Done():
			return nil
		}
	}
}

package kafka

import (
	loggingtype "alle-task-manager-gunish/internal/common/logging"
	"context"
	"errors"
	"github.com/IBM/sarama"
)

type MessageHandler func(message *sarama.ConsumerMessage) error

type Consumer struct {
	consumer sarama.ConsumerGroup
	topics   []string
	handler  MessageHandler
}

func NewConsumer(brokers []string, groupID string, topics []string, handler MessageHandler) (*Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	consumer, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		consumer: consumer,
		topics:   topics,
		handler:  handler,
	}, nil
}

type ConsumerGroupHandler struct {
	handler MessageHandler
}

func (h *ConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (h *ConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (h *ConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	logger := loggingtype.GetLogger()
	for message := range claim.Messages() {
		if err := h.handler(message); err != nil {
			logger.Error("Error processing message", "error", err)
			continue
		}
		session.MarkMessage(message, "")
	}
	return nil
}

func (c *Consumer) Start(ctx context.Context) error {
	handler := &ConsumerGroupHandler{handler: c.handler}
	logger := loggingtype.GetLogger()
	for {
		err := c.consumer.Consume(ctx, c.topics, handler)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return nil
			}
			logger.Error("Error from consumer:", "error", err)
		}
	}
}

func (c *Consumer) Close() error {
	return c.consumer.Close()
}

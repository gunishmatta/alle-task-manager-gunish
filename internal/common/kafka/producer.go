package kafka

import (
	loggingtype "alle-task-manager-gunish/internal/common/logging"
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
)

type SyncProducer interface {
	SendMessage(msg *sarama.ProducerMessage) (partition int32, offset int64, err error)
	Close() error
}

type Producer struct {
	Producer SyncProducer
}

func NewProducer(brokers []string) (*Producer, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}

	return &Producer{Producer: producer}, nil
}

func (p *Producer) PublishMessage(topic string, key string, value interface{}) error {
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.StringEncoder(jsonValue),
	}

	partition, offset, err := p.Producer.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}
	loggingtype.GetLogger().Info("Message published", "topic", topic, "partition", partition, "offset", offset)
	return nil
}

func (p *Producer) Close() error {
	return p.Producer.Close()
}

package internalmessagebroker //nolint

import (
	"encoding/json"
	"fmt"

	"github.com/IBM/sarama"
	"github.com/dostrovskiy/otus-golang-home-work/hw12_13_14_15_16_calendar/internal/app"
	"github.com/dostrovskiy/otus-golang-home-work/hw12_13_14_15_16_calendar/internal/storage"
)

type Producer struct {
	logger  app.Logger
	brokers []string // [localhost:9092]
	topic   string
}

func NewProducer(logger app.Logger, brokers []string, topic string) *Producer {
	return &Producer{logger: logger, brokers: brokers, topic: topic}
}

func (s *Producer) Send(event *storage.Event) error {
	s.logger.Debug("Sending event [%+v] to Kafka", event)
	producer, err := sarama.NewSyncProducer(s.brokers, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err := producer.Close(); err != nil {
			s.logger.Error("failed to close producer: %s", err.Error())
		}
	}()

	eventJSON, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("sender failed to send event: %w", err)
	}
	msg := &sarama.ProducerMessage{Topic: s.topic, Value: sarama.StringEncoder(eventJSON)}
	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("sender failed to send event: %w", err)
	}
	s.logger.Debug("Event [%+v] sent to partition %d at offset %d", event, partition, offset)
	return nil
}
//nolint
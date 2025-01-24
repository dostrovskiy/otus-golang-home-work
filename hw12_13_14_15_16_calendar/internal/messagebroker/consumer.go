package internalmessagebroker //nolint

import (
	"context"
	"encoding/json"
	"time"

	"github.com/IBM/sarama"
	"github.com/dostrovskiy/otus-golang-home-work/hw12_13_14_15_16_calendar/internal/app"
	"github.com/dostrovskiy/otus-golang-home-work/hw12_13_14_15_16_calendar/internal/storage"
)

type Consumer struct {
	app        app.Application
	logger     app.Logger
	brokers    []string
	topic      string
	admin      *ClusterAdmin
	maxRetries int
}

func NewConsumer(app app.Application, logger app.Logger, brokers []string, topic string, maxRetries int) *Consumer {
	return &Consumer{app: app, logger: logger, brokers: brokers, topic: topic,
		admin: NewClusterAdmin(logger, brokers), maxRetries: maxRetries}
}

func (c *Consumer) Start(ctx context.Context) error {
	config := sarama.NewConfig()
	var kafkaConsumer sarama.Consumer
	var err error

	// Connect to Kafka
	retryDelay := time.Second
	for i := 0; i < c.maxRetries; i++ {
		kafkaConsumer, err = sarama.NewConsumer(c.brokers, config)
		if err == nil {
			break
		}
		c.logger.Warn("failed to create kafka consumer: %s, retrying in %s", err.Error(), retryDelay)
		time.Sleep(retryDelay)
		retryDelay *= 2 // Exponential backoff
	}
	defer func() {
		if err := kafkaConsumer.Close(); err != nil {
			c.logger.Error("failed to close kafka consumer: %s", err.Error())
		}
	}()

	// Connect to partition
	partitionConsumer, err := kafkaConsumer.ConsumePartition(c.topic, 0, sarama.OffsetNewest)
	if err != nil {
		return err
	}
	defer func() {
		if err := partitionConsumer.Close(); err != nil {
			c.logger.Error("failed to close kafka partition consumer: %s", err.Error())
		}
	}()

	// Consume messages
	for {
		select {
		case <-ctx.Done():
			return nil
		case msg := <-partitionConsumer.Messages():
			event := storage.Event{}
			if err := json.Unmarshal(msg.Value, &event); err != nil {
				c.logger.Error("consumer failed to unmarshal message: %s", err.Error())
				continue
			}
			if _, err := c.app.AddEventNotification(ctx, &event); err != nil {
				c.logger.Error("consumer failed to add notification: %s", err.Error())
			}
		case err := <-partitionConsumer.Errors():
			c.logger.Error("consumer error: %s", err.Error())
		}
	}
}
//nolint
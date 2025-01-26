package internalmessagebroker //nolint

import (
	"github.com/IBM/sarama"
	"github.com/dostrovskiy/otus-golang-home-work/hw12_13_14_15_16_calendar/internal/app"
)

type ClusterAdmin struct {
	logger  app.Logger
	brokers []string
}

func NewClusterAdmin(logger app.Logger, brokers []string) *ClusterAdmin {
	return &ClusterAdmin{logger: logger, brokers: brokers}
}

func (a *ClusterAdmin) CreateTopic(topic string) error {
	config := sarama.NewConfig()
	admin, err := sarama.NewClusterAdmin(a.brokers, config)
	if err != nil {
		return err
	}
	defer func() {
		if err := admin.Close(); err != nil {
			a.logger.Error("failed to close admin: %s", err.Error())
		}
	}()

	return admin.CreateTopic(topic, &sarama.TopicDetail{
		NumPartitions:     1,
		ReplicationFactor: 1,
	}, false)
}

func (a *ClusterAdmin) TopicExists(topic string) (bool, error) {
	config := sarama.NewConfig()
	admin, err := sarama.NewClusterAdmin(a.brokers, config)
	if err != nil {
		return false, err
	}
	defer func() {
		if err := admin.Close(); err != nil {
			a.logger.Error("failed to close admin: %s", err.Error())
		}
	}()

	topics, err := admin.ListTopics()
	if err != nil {
		return false, err
	}
	_, exists := topics[topic]
	return exists, nil
}
//nolint
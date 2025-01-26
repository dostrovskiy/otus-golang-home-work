package main

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Logger     LoggerConf
	DataSource DataSourceConf
	Kafka      KafkaConf
	Scheduler  SchedulerConf
}

type LoggerConf struct {
	Level string
}

type DataSourceConf struct {
	StorageType string
	Dsn         string
}

type KafkaConf struct {
	Brokers []string
	Topic   string
}

type SchedulerConf struct {
	Interval time.Duration
}

func NewConfig() Config {
	return Config{}
}

func LoadConfig(filename string) (*Config, error) {
	viper.SetConfigFile(filename)
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}
	return &config, nil
}

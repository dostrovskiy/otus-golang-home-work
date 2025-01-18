package main

import (
	"github.com/spf13/viper"
)

// При желании конфигурацию можно вынести в internal/config.
// Организация конфига в main принуждает нас сужать API компонентов, использовать
// при их конструировании только необходимые параметры, а также уменьшает вероятность циклической зависимости.
type Config struct {
	Logger     LoggerConf
	DataSource DataSourceConf
	Server     ServerConf
}

type LoggerConf struct {
	Level string
}

type DataSourceConf struct {
	StorageType string
	Dsn         string
}

type ServerConf struct {
	Address string
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

package integrationtests

import (
	"github.com/spf13/viper"
)

type Config struct {
	DataSource DataSourceConf
	Server     ServerConf
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

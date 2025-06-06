package configs

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	Kafka    KafkaConfig    `yaml:"kafka"`
}

type ServerConfig struct {
	Address         string        `yaml:"address"`
	GRPCAddress     string        `yaml:"grpc_address"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout"`
}

type DatabaseConfig struct {
	Type     string        `yaml:"type"`
	Database MongoDBConfig `yaml:"mongodb"`
}

type MongoDBConfig struct {
	URI      string `yaml:"uri"`
	Database string `yaml:"database"`
}

type KafkaConfig struct {
	Brokers      []string `yaml:"brokers"`
	ProduceTopic []string `yaml:"ptopics"`
	ConsumeTopic []string `yaml:"ctopics"`
	GroupID      string   `yaml:"groupID"`
}

func Load() *Config {
	data, err := os.ReadFile("./auction/configs/config.yaml")
	if err != nil {
		panic(fmt.Sprintf("Error reading config file: %v", err))
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		panic(fmt.Sprintf("Error parsing config: %v", err))
	}

	return &cfg
}

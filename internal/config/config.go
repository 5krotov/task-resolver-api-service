package config

import (
	"fmt"
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	ApiService   ApiServiceConfig   `yaml:"api-service" validate:"required"`
	Agent        AgentConfig        `yaml:"agent" validate:"required"`
	DataProvider DataProviderConfig `yaml:"data-provider" validate:"required"`
	GRPCClient   GRPCClientConfig   `yaml:"grpc-client" validate:"required"`
}

type ApiServiceConfig struct {
	HTTP HTTPConfig `yaml:"http" validate:"required"`
}

type AgentConfig struct {
	Addr           string `yaml:"addr" validate:"required"`
	UseTLS         bool   `yaml:"useTLS" validate:"required"`
	GrpcServerName string `yaml:"grpcServerName" validate:"required"`
}

type DataProviderConfig struct {
	Addr           string `yaml:"addr" validate:"required"`
	UseTLS         bool   `yaml:"useTLS" validate:"required"`
	GrpcServerName string `yaml:"grpcServerName" validate:"required"`
}

type GRPCClientConfig struct {
	CaCert     string `yaml:"caCert" validate:"required"`
	ClientCert string `yaml:"clientCert" validate:"required"`
	ClientKey  string `yaml:"clientKey" validate:"required"`
}

type HTTPConfig struct {
	Addr string `yaml:"addr" validate:"required"`
}

func NewConfig() *Config {
	return &Config{}
}

func (c *Config) Load(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	if err := yaml.Unmarshal(data, c); err != nil {
		return fmt.Errorf("failed to parse config file: %w", err)
	}
	return nil
}

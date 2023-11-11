package config

import (
	"fmt"
	"os"
)

const (
	envServicePort = "SERVICE_PORT"
)

type Config struct {
	
	Port string `yaml:"port"`
	
}

func NewServerConfig(servicePort string) *Config {
	return &Config{
		Port: servicePort,
	}
}

// EnvConfig returns a Config struct populated from environment variables.
func EnvServerConfig() (*Config, error) {
	var servicePort string
	var found bool

	if servicePort, found = os.LookupEnv(envServicePort); !found {
		return nil, fmt.Errorf("env var %s not found", envServicePort)
	}

	config := NewServerConfig( servicePort )

	return config, nil
}
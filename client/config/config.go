package config

import (
	"fmt"
	"os"
)

const (
	envUploadUrl = "UPLOAD_URL"
)

type Config struct {
	
	UploadUrl string `yaml:"uploadURL"`
	
}

func NewClientConfig(uploadURL string) *Config {
	return &Config{
		UploadUrl: uploadURL,
	}
}

// EnvConfig returns a Config struct populated from environment variables.
func EnvClientConfig() (*Config, error) {
	var uploadURL string
	var found bool

	if uploadURL, found = os.LookupEnv(envUploadUrl); !found {
		return nil, fmt.Errorf("env var %s not found", envUploadUrl)
	}

	config := NewClientConfig( uploadURL )

	return config, nil
}
package stores

import (
	"fmt"
	"os"
)

type StoreType string

const(
	envStoreType = "STORE_TYPE"
)

const (
	StoreTypeFile StoreType = "file"
	StoreTypeS3 StoreType = "s3"
)

type IStore interface {
	Connect() error
	StoreFile(name string) error
	GetFile() error
	StoreMerkleTree() error
	GetMerkleTree() error
	Disconnect() error
}

type Config struct {

	Type StoreType `yaml:"type"`
}

func NewConfig(storeType StoreType) *Config {
	return &Config{
		Type: storeType,
	}
}

// EnvConfig returns a Config struct populated from environment variables.
func EnvConfig() (*Config, error) {
	var storeType string
	var found bool

	if storeType, found = os.LookupEnv(envStoreType); !found {
		return nil, fmt.Errorf("env var %s not found", envStoreType)
	}

	config := NewConfig( StoreType(storeType) )

	return config, nil
}

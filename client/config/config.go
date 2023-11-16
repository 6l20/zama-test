package config

import (
	"fmt"
	"os"
)

const (
	envUploadUrl   = "UPLOAD_URL"
	envDownloadUrl = "DOWNLOAD_URL"
	envProofUrl    = "PROOF_URL"
	envRootFile    = "ROOT_FILE"
)

type Config struct {
	UploadURL   string `yaml:"uploadURL"`
	DownloadURL string `yaml:"downloadURL"`
	ProofURL    string `yaml:"proofURL"`
	RootFile    string `yaml:"rootFile"`
}

func NewClientConfig(uploadURL, downloadURL, proofURL, rootFile string) *Config {
	return &Config{
		UploadURL:   uploadURL,
		DownloadURL: downloadURL,
		ProofURL:   proofURL,
		RootFile:    rootFile,
	}
}

// EnvConfig returns a Config struct populated from environment variables.
func EnvClientConfig() (*Config, error) {
	var uploadURL, downloadURL, proofURL, rootFile string
	var found bool

	if uploadURL, found = os.LookupEnv(envUploadUrl); !found {
		return nil, fmt.Errorf("env var %s not found", envUploadUrl)
	}

	if downloadURL, found = os.LookupEnv(envDownloadUrl); !found {
		return nil, fmt.Errorf("env var %s not found", envDownloadUrl)
	}

	if proofURL, found = os.LookupEnv(envProofUrl); !found {
		return nil, fmt.Errorf("env var %s not found", envProofUrl)
	}

	if rootFile, found = os.LookupEnv(envRootFile); !found {
		return nil, fmt.Errorf("env var %s not found", envRootFile)
	}

	config := NewClientConfig(uploadURL, downloadURL, proofURL, rootFile)

	return config, nil
}

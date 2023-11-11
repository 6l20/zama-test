package zap

import (
	"fmt"
	"os"
)

type LoggerLevel string
type LoggerFormat string

const(
	envLogLevel = "LOG_LEVEL"
	envLogFormat = "LOG_FORMAT"
)

const (
	InfoLevel  LoggerLevel = "info"
	ErrorLevel LoggerLevel = "error"
	DebugLevel LoggerLevel = "debug"
	WarnLevel  LoggerLevel = "warn"
	PanicLevel LoggerLevel = "panic"
)

const (
	TextFormat LoggerFormat = "text"
	JSONFormat LoggerFormat = "json"
)

type Config struct {
	Level  LoggerLevel
	Format LoggerFormat
}

func NewConfig(level LoggerLevel, format LoggerFormat) *Config {
	return &Config{
		Level:  level,
		Format: format,
	}
}

func EnvConfig() (*Config, error) {
	var level string
	var format string
	var found bool

	if level, found = os.LookupEnv(envLogLevel); !found {
		return nil, fmt.Errorf("env var %s not found", envLogLevel)
	}

	if format, found = os.LookupEnv(envLogFormat); !found {
		return nil, fmt.Errorf("env var %s not found", envLogFormat)
	}

	config := NewConfig( LoggerLevel(level), LoggerFormat(format) )

	return config, nil
}
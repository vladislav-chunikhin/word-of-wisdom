package config

import (
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

const (
	configFileEnv = "CONFIG_FILE"
)

// Config - config structure
type Config struct {
	Server Server `yaml:"server" env-prefix:"SERVER_"`
	POW    POW    `yaml:"pow" env-prefix:"POW_"`
}

// Server - server config structure
type Server struct {
	LogLevel  string        `yaml:"log_level" env:"LOG_LEVEL" env-default:"DEBUG"`
	Address   string        `yaml:"address" env:"ADDRESS" env-default:":8080"`
	KeepAlive time.Duration `yaml:"keep_alive" env:"KEEP_ALIVE" env-default:"15s"`
	Deadline  time.Duration `yaml:"deadline" env:"DEADLINE" env-default:"30s"`
}

// POW - proof of work config structure
type POW struct {
	Complexity byte `yaml:"complexity" env:"COMPLEXITY" env-default:"2"`
}

// Parse - parse config from environment variables or file
func Parse() (config *Config, err error) {
	if path, ok := os.LookupEnv(configFileEnv); ok {
		return parseFromFile(path)
	}

	return parseFromEnv()
}

// parseFromEnv - parse config from environment variables
func parseFromEnv() (*Config, error) {
	var config Config

	if err := cleanenv.ReadEnv(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

// parseFromFile - parse config from file
func parseFromFile(path string) (*Config, error) {
	var config Config

	if err := cleanenv.ReadConfig(path, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

package config

import (
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

const (
	serverConfigFileEnv = "SERVER_CONFIG_FILE"
	clientConfigFileEnv = "CLIENT_CONFIG_FILE"
)

type ServerConfig struct {
	Server Server `yaml:"server" env-prefix:"SERVER_"`
	POW    POW    `yaml:"pow" env-prefix:"POW_"`
}

type ClientConfig struct {
	ServerAddr    string `yaml:"server_address" env:"SERVER_ADDRESS" env-default:"127.0.0.1:8080"`
	RPS           int    `yaml:"rps" env:"RPS" env-default:"5"`
	TotalRequests int    `yaml:"total_requests" env:"TOTAL_REQUESTS" env-default:"100"`
	LogLevel      string `yaml:"log_level" env:"LOG_LEVEL" env-default:"DEBUG"`
}

type Server struct {
	LogLevel        string        `yaml:"log_level" env:"LOG_LEVEL" env-default:"DEBUG"`
	Address         string        `yaml:"address" env:"ADDRESS" env-default:":8080"`
	WorkerCount     int           `yaml:"worker_count" env:"WORKER_COUNT" env-default:"10"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout" env:"SHUTDOWN_TIMEOUT" env-default:"10s"`
}

type POW struct {
	Complexity byte          `yaml:"complexity" env:"COMPLEXITY" env-default:"2"`
	Timeout    time.Duration `yaml:"timeout" env:"TIMEOUT" env-default:"5s"`
}

// ServerParse - parse server config from environment variables or file
func ServerParse() (config *ServerConfig, err error) {
	if path, ok := os.LookupEnv(serverConfigFileEnv); ok {
		return parseFromFile(path)
	}

	return parseFromEnv()
}

// ClientParse - parse client config from environment variables or file
func ClientParse() (config *ClientConfig, err error) {
	if path, ok := os.LookupEnv(clientConfigFileEnv); ok {
		return parseClientFromFile(path)
	}

	return parseClientFromEnv()
}

// parseFromEnv - parse config from environment variables
func parseFromEnv() (*ServerConfig, error) {
	var config ServerConfig

	if err := cleanenv.ReadEnv(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

// parseClientFromEnv - parse client config from environment variables
func parseClientFromEnv() (*ClientConfig, error) {
	var config ClientConfig

	if err := cleanenv.ReadEnv(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

// parseFromFile - parse config from file
func parseFromFile(path string) (*ServerConfig, error) {
	var config ServerConfig

	if err := cleanenv.ReadConfig(path, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

// parseClientFromFile - parse client config from file
func parseClientFromFile(path string) (*ClientConfig, error) {
	var config ClientConfig

	if err := cleanenv.ReadConfig(path, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

package config

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type LoggerKey string

const LoggerContextKey LoggerKey = "logger"

type Config struct {
	Main MainConfig `yaml:"main"`
	App  AppConfig  `yaml:"app"`
}

type MainConfig struct {
	Port              string        `yaml:"port"`
	ReadTimeout       time.Duration `yaml:"read_timeout"`
	WriteTimeout      time.Duration `yaml:"write_timeout"`
	ReadHeaderTimeout time.Duration `yaml:"read_header_timeout"`
	IdleTimeout       time.Duration `yaml:"idle_timeout"`
	ShutdownTimeout   time.Duration `yaml:"shutdown_timeout"`
}

type AppConfig struct {
	MazeCount      int `yaml:"maze_count"`
	AlgorithmCount int `yaml:"algorithm_count"`
}

func MustLoadConfig(path string, logger *slog.Logger) *Config {
	cfg := &Config{}

	file, err := os.Open(path)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to open config file: %v", err))
		return &Config{}
	}
	defer file.Close()

	if err = yaml.NewDecoder(file).Decode(cfg); err != nil {
		logger.Error(fmt.Sprintf("failed to decode config file: %v", err))
		return &Config{}
	}

	return cfg
}

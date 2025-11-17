package config

import (
	"encoding/json"
	"fmt"
	"os"
)

var Env string = "development"

type DBConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type Config struct {
	Database DBConfig `json:"database"`
	Secret   string   `json:"secret"`
}

var AppConfig *Config

func SetConfig(cfg *Config) {
	AppConfig = cfg
}

func GetConfig() *Config {
	return AppConfig
}

func LoadConfig(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	var cfg Config
	if err := decoder.Decode(&cfg); err != nil {
		return fmt.Errorf("failed to decode config file: %w", err)
	}

	AppConfig = &cfg
	return nil
}

func InitTestConfig(secret string) {
	AppConfig = &Config{Secret: secret}
}

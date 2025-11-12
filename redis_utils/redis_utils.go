package redis_utils

import (
	"encoding/json"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

var rdb *redis.Client
var AppConfig EnvConfig

type RedisConfig struct {
	Addr     string `json:"addr"`
	Password string `json:"password"`
}

type CookieConfig struct {
	Domain string `json:"domain"`
}

type EnvConfig struct {
	Redis  RedisConfig  `json:"redis"`
	Cookie CookieConfig `json:"cookie"`
}

func InitRedis(env string) {
	configData, err := os.ReadFile("./config/config.json")
	if err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}

	var fullConfig map[string]EnvConfig
	if err := json.Unmarshal(configData, &fullConfig); err != nil {
		log.Fatalf("Failed to parse config JSON: %v", err)
	}

	envConfig, ok := fullConfig[env]
	if !ok {
		log.Fatalf("Environment '%s' not found in config file", env)
	}

	AppConfig = envConfig

	rdb = redis.NewClient(&redis.Options{
		Addr:            envConfig.Redis.Addr,
		Password:        envConfig.Redis.Password,
		ClientName:      "myapp",
		DisableIdentity: true,
	})
}

func GetRedisClient() *redis.Client {
	return rdb
}

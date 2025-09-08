package config

import (
	"log"
	"os"
)

type Config struct {
	ApiKey    string
	RedisAddr string
	GRPCAddr  string
}

func Load() *Config {
	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		log.Fatal("API_KEY is required")
	}

	return &Config{
		ApiKey:    apiKey,
		RedisAddr: getEnv("REDIS_ADDR", "localhost:6379"),
		GRPCAddr:  getEnv("GRPC_ADDR", ":50051"),
	}
}

func getEnv(key, def string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return def
}

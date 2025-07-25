package config

import (
	"context"
	"os"
)

type Config struct {
	MongoURI    string
	MongoDBName string
	JWTSecret   string
	Ctx         context.Context
}

func Load() *Config {
	return &Config{
		MongoURI:    getEnv("MONGO_URI", "mongodb://localhost:27017"),
		MongoDBName: getEnv("MONGO_DB_NAME", "mydb"),
		JWTSecret:   getEnv("JWT_SECRET", "supersecret"),
		Ctx:         context.Background(),
	}
}

func getEnv(key, defaultVal string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	return val
}

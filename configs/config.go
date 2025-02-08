package configs

import (
	"errors"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port          string
	DatabaseURL   string
	RedisAddress  string
	RedisPassword string
	RedisDB       int
}

func InitConfig() (Config, error) {
	err := godotenv.Load()
	if err != nil {
		return Config{}, err
	}

	if os.Getenv("PORT") == "" {
		return Config{}, errors.New("PORT is not set")
	}

	if os.Getenv("DATABASE_URL") == "" {
		return Config{}, errors.New("DATABASE_URL is not set")
	}

	if os.Getenv("REDIS_ADDR") == "" {
		return Config{}, errors.New("REDIS_ADDR is not set")
	}

	redisDB, err := strconv.ParseInt(os.Getenv("REDIS_DB"), 10, 64)
	if err != nil {
		// DEFAULT TO 0
		redisDB = 0
	}

	return Config{
		Port:          os.Getenv("PORT"),
		DatabaseURL:   os.Getenv("DATABASE_URL"),
		RedisAddress:  os.Getenv("REDIS_ADDR"),
		RedisPassword: os.Getenv("REDIS_PASSWORD"),
		RedisDB:       int(redisDB),
	}, nil
}

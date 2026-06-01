package cache

import (
	"context"
	"log"
	"os"
	"strconv"

	"github.com/redis/go-redis/v9"
)

func NewRedisClient() *redis.Client {
	db := 0

	if os.Getenv("REDIS_DB") != "" {
		parsedDB, err := strconv.Atoi(os.Getenv("REDIS_DB"))
		if err == nil {
			db = parsedDB
		}
	}

	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Username: os.Getenv("REDIS_USERNAME"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       db,
	})

	ctx := context.Background()

	if err := client.Ping(ctx).Err(); err != nil {
		log.Println("failed to connect to Redis:", err)
	} else {
		log.Println("connected to Redis successfully")
	}

	return client
}
package cache

import (
	"context"
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
)

// Ctx is context variable for Redis
var Ctx = context.Background()

// RedisConnection returns redis client
func RedisConnection() *redis.Client {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	envs, err := godotenv.Read(".env")
	if err != nil {
		log.Fatal("Unable to read .env file")
	}
	redisAddr := envs["REDIS_ADDR"]
	redisPswrd := envs["REDIS_PSWRD"]

	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPswrd,
		DB:       0,
	})

	if _, err := rdb.Ping(Ctx).Result(); err != nil {
		log.Fatalf("Could not ping redis server due to err: %s \n", err)
	}

	fmt.Println("Successfully connected to Redis!")
	return rdb
}

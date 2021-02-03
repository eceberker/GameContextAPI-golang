package cache

import (
	"context"
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
)

var Ctx = context.Background()

// RedisConnection returns redis client
func RedisConnection() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	if _, err := rdb.Ping(Ctx).Result(); err != nil {
		log.Fatalf("Could not ping redis server due to err: %s \n", err)
	}

	fmt.Println("Successfully connected to Redis!")
	return rdb
}

package cache

import (
	"context"
	"fmt"
	"log"
	"time"

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
	redisAddr := envs["REDIS_URL"]

	var rdb *redis.Client
	var er error

	tries := 10

	for tries > 0 {
		rdb = redis.NewClient(&redis.Options{
			Addr:     redisAddr,
			Password: "",
			DB:       0,
		})
		if _, er = rdb.Ping(Ctx).Result(); er != nil {
			fmt.Printf("Could not ping redis server due to err: %s \n", er)
			time.Sleep(time.Second * 2)
			tries--
		} else {
			er = nil
			break
		}
	}
	if er != nil {
		rdb = nil
		return rdb
	}

	fmt.Println("Successfully connected to Redis!")
	return rdb

}

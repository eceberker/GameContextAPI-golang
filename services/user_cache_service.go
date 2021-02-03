package services

import (
	"fmt"
	"strconv"

	"github.com/eceberker/gamecontextdb/cache"
	"github.com/eceberker/gamecontextdb/helpers"
	"github.com/eceberker/gamecontextdb/models"
	"github.com/go-redis/redis/v8"
)

// GetUserCache retrieves user by Id
func GetUserCache(id int64) (models.User, error) {
	// create redis connection
	rdb := cache.RedisConnection()

	// Return model
	var user models.User

	// UserId to string
	userID := helpers.Int64ToString(id)

	// Check Redis to key exists or not
	_, err := rdb.Get(cache.Ctx, userID).Result()
	if err == redis.Nil {
		user = GetUserFromDb(id)

		rdb.HSet(cache.Ctx, userID, []string{"display_name", user.Name, "country", user.Country, "points", strconv.Itoa(int(user.Points))})
	}

	cachedUserName, err := rdb.HGet(cache.Ctx, userID, "display_name").Result()
	if err != nil {
		fmt.Printf("Unable to retrieve user struct into redis due to: %s \n", err)
	}
	cachedUserCountry, err := rdb.HGet(cache.Ctx, userID, "country").Result()
	if err != nil {
		fmt.Printf("Unable to retrieve user struct into redis due to: %s \n", err)
	}
	cachedUserPoints, err := rdb.HGet(cache.Ctx, userID, "points").Result()
	if err != nil {
		fmt.Printf("Unable to retrieve user struct into redis due to: %s \n", err)
	}
	points, err := helpers.StringToInt64(cachedUserPoints)
	if err != nil {
		fmt.Printf("Unable to convert string to int")
	}

	user.ID = id
	user.Name = cachedUserName
	user.Country = cachedUserCountry
	user.Points = points

	return user, nil
}

// InsertUserToCache inserts given user
func InsertUserToCache(user models.User) models.User {
	// create redis connection
	rdb := cache.RedisConnection()

	rdb.HSet(cache.Ctx, helpers.Int64ToString(user.ID), []string{"display_name", user.Name, "country", user.Country, "points", helpers.Int64ToString(user.Points)})

	return user
}

// UpdateUserScoreCache updates given user's score from hash set
func UpdateUserScoreCache(score models.ScoreSubmit) models.User {
	// create redis connection
	rdb := cache.RedisConnection()

	// retrieve user cache
	user, _ := GetUserCache(score.UserID)

	// UserId and points to string
	userID := helpers.Int64ToString(user.ID)
	points := helpers.Int64ToString(user.Points + score.ScoreWorth)

	// Update Cache
	rdb.HSet(cache.Ctx, userID, "points", points)

	// Retrieve updated cache
	user, _ = GetUserCache(score.UserID)
	return user
}

// GetAllUsersCache returns all users
func GetAllUsersCache() []models.User {

	//CACHE CONTROL TO BE IMPLEMENTED

	users, err := GetAllUsers()
	if err != nil {
		fmt.Printf("Unable to get users %v", err)
	}
	return users
}

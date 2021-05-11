package main

import (
	"fmt"

	"github.com/go-redis/redis"
)

var redisClient *redis.Client

func initRedis() (err error) {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})

	_, err = redisClient.Ping().Result()
	return
}

func main() {
	err := initRedis()
	if err != nil {
		fmt.Println("connect redis faild, err: ", err)
		return
	}

	redisClient.Set().Result()
}

package db

import "github.com/go-redis/redis/v8"

var Redis *redis.Client

func RedisInit() {
	Redis = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
}

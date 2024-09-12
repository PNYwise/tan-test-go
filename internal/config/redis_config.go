package config

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

func RedisConn(ctx context.Context, conf *viper.Viper) *redis.Client {

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", conf.GetString("redis.host"), conf.GetString("redis.port")),
		Password: conf.GetString("redis.password"),
		DB:       0,
	})

	// Check the connection
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	} else {
		log.Printf("Connected to Redis \n")
	}

	return rdb
}

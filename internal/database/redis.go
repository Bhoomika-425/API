package database

import (
	"fmt"
	"project/config"
	"strconv"

	"github.com/redis/go-redis/v9"
)

func ConnectionToRedis() *redis.Client {
	cfg := config.GetConfig()
	PASS, _ := strconv.Atoi(cfg.Redisconfig.Database)

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf(":%s", cfg.Redisconfig.Address),
		Password: cfg.Redisconfig.Password,
		DB:       PASS,
	})
	return rdb
}


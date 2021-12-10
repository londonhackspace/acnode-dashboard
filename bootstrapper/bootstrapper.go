package main

import (
	"github.com/go-redis/redis/v8"
	"github.com/londonhackspace/acnode-dashboard/auth"
	"github.com/londonhackspace/acnode-dashboard/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	conf := config.GetCombinedConfig("acserverdash.json")

	if !conf.RedisEnable {
		log.Fatal().Msg("Redis not enabled so nowhere to create a user")
		os.Exit(1)
	}

	redisConn := redis.NewClient(&redis.Options{
		Addr:     conf.RedisServer,
		Password: "",
		DB:       0,
	})

	userStore := auth.CreateRedisProvider(redisConn)

	auth.CreateInitialUser(userStore)
}

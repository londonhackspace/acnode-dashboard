package main

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"github.com/londonhackspace/acnode-dashboard/acserver_api"
	"github.com/londonhackspace/acnode-dashboard/config"
	"github.com/londonhackspace/acnode-dashboard/usagelogs"
	"strconv"
)

func main() {
	conf := config.GetCombinedConfig("acserverdash.json")

	acserver := acserver_api.CreateACServer(&conf)

	redisConn := redis.NewClient(&redis.Options{
		Addr: conf.RedisServer,
		Password: "",
		DB: 0,
	})

	ctx := context.Background()

	usageKeys := redisConn.Keys(ctx, "nodeusage:*").Val()

	for _,val := range usageKeys {
		// delete the scratch space if it exists
		nt := redisConn.Type(ctx, "new_"+val).Val()
		if nt != "none" {
			println("Deleting scratch space new_"+val)
			redisConn.Del(ctx, "new_"+val)
		}

		llen := redisConn.LLen(ctx, val).Val()
		for idx,data := range redisConn.LRange(ctx, val, 0, llen).Val() {
			if idx % 10 == 0 {
				println("Processing " + strconv.Itoa(idx) + "/" + strconv.Itoa(int(llen)) + " of " + val)
			}

			entry := usagelogs.LogEntry{}
			json.Unmarshal([]byte(data), &entry)

			if len(entry.Name) == 0 {
				apirec := acserver.GetUserRecordForCard(entry.Card)
				if apirec != nil {
					entry.Name = apirec.UserName
					entry.UserId = apirec.Id
				} else {
					println("Couldn't get name for card " + entry.Card)
				}
			}

			remarshalled,_ := json.Marshal(&entry)
			redisConn.LPush(ctx, "new_"+val, remarshalled).Result()
		}
		redisConn.Del(ctx, val).Result()
		redisConn.Rename(ctx, "new_"+val, val).Result()
	}
}

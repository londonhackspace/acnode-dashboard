package usagelogs

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"github.com/londonhackspace/acnode-dashboard/acnode"
	"github.com/londonhackspace/acnode-dashboard/acserver_api"
	"github.com/rs/zerolog/log"
	"strings"
	"time"
)

type RedisUsageLogger struct {
	redis *redis.Client
	ctx context.Context
	acserver *acserver_api.ACServer
}

func CreateRedisUsageLogger(redis *redis.Client, acserver *acserver_api.ACServer) UsageLogger {
	return &RedisUsageLogger{
		redis: redis,
		ctx: context.Background(),
		acserver: acserver,
	}
}

func (rul *RedisUsageLogger) AddUsageLog(node *acnode.ACNode, msg acnode.Announcement) {


	log := LogEntry{
		Timestamp: time.Now(),
		Card:      msg.Card,
		Node:      (*node).GetMqttName(),
		Success:   msg.Granted != 0,
		// Name: user_name,
	}

	data,_ := json.Marshal(log)

	rul.redis.LPush(rul.ctx, "nodeusage:"+(*node).GetMqttName(), string(data)).Result()
	// maybe trim the list?
	//rul.redis.LTrim(rul.ctx, "nodeusageraw:"+(*node).GetMqttName(), 0, 10000)
}

func (rul *RedisUsageLogger) GetUsageLogNodes() []string {
	result := []string{}

	for _,key := range rul.redis.Keys(rul.ctx, "nodeusage:*").Val() {
		parts := strings.Split(key, ":")
		result = append(result, parts[1])
	}

	return result
}

func (rul *RedisUsageLogger) GetUsageLogCountForNode(node string) int64 {
	res,err := rul.redis.LLen(rul.ctx, "nodeusage:"+node).Result()
	if err != nil {
		log.Err(err).Str("node", node).Msg("Error getting count of usage logger entries for node")
	}
	return res
}

func (rul *RedisUsageLogger) GetUsageLogsForNode(node string, from int64, to int64) []LogEntry {
	out := []LogEntry{}
	res, err := rul.redis.LRange(rul.ctx, "nodeusage:"+node, from, to).Result()

	if err != nil {
		log.Err(err).Str("node", node).Msg("Error getting logs")
	}

	for _,item := range res {
		le := LogEntry{}
		json.Unmarshal([]byte(item), &le)
		out = append(out, le)
	}

	return out
}
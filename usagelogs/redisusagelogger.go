package usagelogs

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"github.com/londonhackspace/acnode-dashboard/acnode"
	"time"
)

type RedisUsageLogger struct {
	redis *redis.Client
	ctx context.Context
}

func CreateRedisUsageLogger(redis *redis.Client) UsageLogger {
	return &RedisUsageLogger{
		redis: redis,
		ctx: context.Background(),
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
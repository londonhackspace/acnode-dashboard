package acnode

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog/log"
)

type RedisNodePersistence struct {
	redis *redis.Client
}

func GetRedisNodePersistence(redis *redis.Client) NodePersistence {
	return &RedisNodePersistence{
		redis: redis,
	}
}

func (np *RedisNodePersistence) GetNodeByMqttName(name string) (*ACNodeRec,error) {
	ctx := context.Background()
	return np.getNodeFromRedis(ctx, "node:" + name)
}

func (np *RedisNodePersistence) StoreNode(node *ACNodeRec) (*ACNodeRec, error) {
	data, _ := json.Marshal(node)
	res := np.redis.Set(context.Background(), "node:" + node.MqttName, string(data), 0)
	if err := res.Err(); err != nil {
		return nil, err
	}

	return node, nil
}

func (np *RedisNodePersistence) getNodeFromRedis(ctx context.Context,key string) (*ACNodeRec, error) {
	item,err := np.redis.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	rec := ACNodeRec{}
	err = json.Unmarshal([]byte(item), &rec)
	if err != nil {
		return nil, err
	}

	return &rec, nil
}

func (np *RedisNodePersistence) GetAllNodes() ([]ACNodeRec, error) {
	res := []ACNodeRec{}
	ctx := context.Background()

	var cursor uint64
	for {
		var keys []string
		var err error
		keys, cursor, err = np.redis.Scan(ctx, cursor, "node:*", 0).Result()
		if err == nil {
			for _,key := range keys {
				item,err := np.getNodeFromRedis(ctx, key)
				if err == nil {
					res = append(res, *item)
				} else {
					log.Err(err).Str("Key", key).Msg("Error getting node data")
				}
			}

		} else {
			log.Err(err).Msg("Error getting nodes from Redis")
		}

		if cursor == 0 {
			break
		}
	}
	return res, nil
}
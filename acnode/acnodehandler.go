package acnode

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog/log"
	"time"
)

type ACNodeHandler struct {
	nodes []ACNodeRec

	redis *redis.Client
	syncchannel chan bool
}

func CreateACNodeHandler() ACNodeHandler {
	return ACNodeHandler{}
}

// Every minute, this syncs to redis
func (h *ACNodeHandler) syncer(stopper chan bool) {
	for {
		select {
		case <- stopper: {
			return
		}
		default: {
			time.Sleep(time.Minute*1)
			for i,_ := range h.nodes {
				node := &h.nodes[i]
				data, _ := json.Marshal(node)
				h.redis.Set(context.Background(), "node:" + node.MqttName, string(data), 0)
			}
			log.Info().Msg("Dumped nodes to Redis")
		}
		}
	}
}

func (h *ACNodeHandler) SetRedis(r *redis.Client) {
	h.redis = r
	// if we have a sync running already, stop it
	if h.syncchannel != nil {
		h.syncchannel <- true
	}

	// read from redis
	ctx := context.Background()
	iter := h.redis.Scan(ctx, 0, "node*", 0).Iterator()

	for iter.Next(ctx) {
		data := h.redis.Get(ctx, iter.Val()).Val()
		var node ACNodeRec
		err := json.Unmarshal([]byte(data), &node)
		if err != nil {
			log.Err(err).Str("node", iter.Val()).Msg("Error unmarshalling node")
		} else {
			found := false
			for i,_ := range h.nodes {
				if h.nodes[i].MqttName == node.MqttName {
					found = true
					h.nodes[i] = node
					log.Info().Str("node", node.MqttName).Msg("Updated node from Redis")
					break
				}
			}
			if !found {
				log.Info().Str("node", node.MqttName).Msg("Added node from Redis")
				h.AddNode(node)
			}
		}
	}

	h.syncchannel = make(chan bool)
	go h.syncer(h.syncchannel)
}

func (h *ACNodeHandler) AddNode(node ACNodeRec) {
	h.nodes = append(h.nodes, node)
}

func (h *ACNodeHandler) GetNodeByMqttName(name string) ACNode {
	for i, _ := range h.nodes {
		if h.nodes[i].GetMqttName() == name {
			return &h.nodes[i]
		}
	}

	node := ACNodeRec{
		NodeType: NodeTypeTool,
		Name:     name,
		MqttName: name,
	}
	h.AddNode(node)

	// return a ref to the last entry we just added
	return &h.nodes[len(h.nodes)-1]
}

func (h *ACNodeHandler) GetNodes() []ACNode {
	var ret []ACNode
	for i,_ := range h.nodes {
		ret = append(ret, &h.nodes[i])
	}

	return ret
}
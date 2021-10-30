package acnode

import (
	"github.com/go-redis/redis/v8"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/rs/zerolog/log"
)

type acnodeUpdateTrigger interface {
	OnNodeUpdate(node ACNode)
}

var (
	nodeCounter = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "acnode_count",
		Help: "Number of ACNodes tracked",
	})
	updateCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "acnode_update_count",
		Help: "Number of updates to ACNode status",
	})
)

type ACNodeHandler struct {
	persistence NodePersistence

	redis *redis.Client

	listeners map[*HandlerListener]bool
}

func CreateACNodeHandler(persistence NodePersistence) ACNodeHandler {
	return ACNodeHandler{
		persistence: persistence,
		listeners:   make(map[*HandlerListener]bool),
	}
}

func (h *ACNodeHandler) AddListener(l *HandlerListener) {
	h.listeners[l] = true
}

func (h *ACNodeHandler) RemoveListener(l *HandlerListener) {
	delete(h.listeners, l)
	close(l.nodeAdded)
	close(l.nodeChanged)
}

func (h *ACNodeHandler) AddNode(node ACNodeRec) {
	node.updateTrigger = h
	nodeCounter.Inc()
	nrec, err := h.persistence.StoreNode(&node)

	if err != nil {
		log.Err(err).Msg("Error adding node")
	}

	for l := range h.listeners {
		go func() {
			l.nodeAdded <- nrec
		}()
	}
}

func (h *ACNodeHandler) OnNodeUpdate(node ACNode) {
	updateCounter.Inc()
	h.persistence.StoreNode(node.(*ACNodeRec))
	for l := range h.listeners {
		go func() {
			l.nodeChanged <- node
		}()
	}
}

func (h *ACNodeHandler) GetNodeByMqttName(name string) ACNode {
	noderec, err := h.persistence.GetNodeByMqttName(name)
	if err == nil {
		noderec.updateTrigger = h
		return noderec
	}

	node := ACNodeRec{
		NodeType: NodeTypeTool,
		Name:     name,
		MqttName: name,
	}
	h.AddNode(node)

	// return a ref to the last entry we just added
	noderec, _ = h.persistence.GetNodeByMqttName(name)
	noderec.updateTrigger = h
	return noderec
}

func (h *ACNodeHandler) GetNodes() []ACNode {
	var ret []ACNode
	nodes,_ := h.persistence.GetAllNodes()
	for i,_ := range nodes {
		nodes[i].updateTrigger = h
		ret = append(ret, &nodes[i])
	}

	return ret
}
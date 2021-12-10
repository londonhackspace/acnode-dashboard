package acserverwatcher

import (
	"github.com/londonhackspace/acnode-dashboard/acnode"
	"github.com/londonhackspace/acnode-dashboard/acserver_api"
	"time"
)

type Watcher struct {
	Api     acserver_api.ACServer
	Handler *acnode.ACNodeHandler
}

func (w *Watcher) Run() {
	for {
		tools := w.Api.GetTools()
		for _, t := range tools {
			if len(t.MqttName) == 0 {
				continue
			}

			var deducedType int = acnode.NodeTypeTool

			if t.Type == "Tool" {
				deducedType = acnode.NodeTypeTool
			} else if t.Type == "Unrestricted Door" || t.Type == "Restricted Door" {
				deducedType = acnode.NodeTypeDoor
			}

			node := w.Handler.GetNodeByMqttName(t.MqttName)
			node.SetName(t.Name)
			node.SetId(t.Id)
			currentType := node.GetType()

			// Our printer designation is more specific than
			// ACServer's tool designation, so don't overwrite it
			if currentType != acnode.NodeTypePrinter &&
				currentType != deducedType {
				node.SetType(deducedType)
			}

			node.SetInService(t.Status == "Operational")
		}
		time.Sleep(5 * time.Minute)
	}
}

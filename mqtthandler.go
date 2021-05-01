package main

import (
	"encoding/json"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/londonhackspace/acnode-dashboard/acnode"
	"github.com/londonhackspace/acnode-dashboard/config"
	"github.com/londonhackspace/acnode-dashboard/usagelogs"
	"github.com/rs/zerolog/log"
	"strings"
	"time"
)

type MqttHandler struct {
	config *config.Config
	acnodehandler *acnode.ACNodeHandler
	usageLogger usagelogs.UsageLogger
	running bool

	conn mqtt.Client
}

func CreateMQTTHandler(config *config.Config, acnodehandler *acnode.ACNodeHandler, usageLogger usagelogs.UsageLogger) MqttHandler {
	return MqttHandler{
		config: config,
		acnodehandler: acnodehandler,
		usageLogger: usageLogger,
		running: true,
	}
}

func (handler *MqttHandler) cbMessage(client mqtt.Client, msg mqtt.Message) {
	log.Debug().
		Str("Topic", msg.Topic()).
		Msg("Handing message start")

	topicParts := strings.Split(msg.Topic(), "/")

	if len(topicParts) < 4 {
		// Not enough to work with ere
		log.Debug().
			Str("Topic", msg.Topic()).
			Msg("Handing message end short")
		return
	}

	node := handler.acnodehandler.GetNodeByMqttName(topicParts[2])

	// first, fix the type
	if topicParts[1] == "tool"  && node.GetType() == acnode.NodeTypeDoor {
		node.SetType(acnode.NodeTypeTool)
	} else if topicParts[1] == "door" && node.GetType() != acnode.NodeTypeDoor {
		node.SetType(acnode.NodeTypeDoor)
	}

	if (msg.Topic() == "/tool/" + topicParts[2] + "/event/PrinterStateChanged") &&
			(node.GetType() != acnode.NodeTypePrinter) {
		node.SetType(acnode.NodeTypePrinter)
	}

	if topicParts[3] == "announcements" {
		announcement := acnode.Announcement{}
		json.Unmarshal(msg.Payload(), &announcement)

		if announcement.Type == "RFID" {
			if handler.usageLogger != nil {
				handler.usageLogger.AddUsageLog(&node, announcement)
			}
		}

		log.Info().
			Str("Node", node.GetName()).
			Msg("Got announcement from node")
	} else if topicParts[3] == "status" {
		status := acnode.Status{}
		json.Unmarshal(msg.Payload(), &status)

		if status.Type == "START" {
			node.SetLastStarted(time.Now())
			if status.SettingsVersion != 0 {
				node.SetSettingsVersion(status.SettingsVersion)
			}

			if status.EEPROMSettingsVersion != 0 {
				node.SetEepromSettingsVersion(status.EEPROMSettingsVersion)
			}

			if status.ResetCause != "" {
				node.SetResetCause(status.ResetCause)
			}
		} else if status.Type == "ALIVE" {
			node.SetStatusMessage(status.Message)
			// We know the nodes have more than zero memory total,
			// so use that to sanity check the results
			if (status.Mem.HeapUsed + status.Mem.HeapFree) > 0 {
				node.SetMemoryStats(status.Mem.HeapFree, status.Mem.HeapUsed)
			}
		}

		log.Info().
			Str("Node", node.GetName()).
			Msg("Got Node Status")
	}

	// this check prevents Octoprint from marking nodes as recently alive
	if topicParts[3] == "announcements" || topicParts[3] == "status" || topicParts[3] == "bell" {
		node.SetLastSeenMQTT(time.Now())
	}

	log.Debug().
		Str("Topic", msg.Topic()).
		Msg("Handing message end")
}

func (handler *MqttHandler) handleMqtt() {
	for handler.running {
		// delay first so retries work sensibly
		time.Sleep(1 * time.Second)

		// do we need to try to connect?

		if !handler.conn.IsConnected() {
			tok := handler.conn.Connect()
			if tok.Wait() && tok.Error() != nil {
				log.Err(tok.Error()).Msg("Error Connecting to MQTT Server")
				continue
			}
			log.Info().Msg("Connected to MQTT server")
		}
	}
}

func (handler *MqttHandler) Init() {
	if handler.conn != nil {
		handler.conn.Disconnect(250)
	}

	opts := mqtt.NewClientOptions().SetClientID(handler.config.MqttClientId)
	opts.AddBroker(handler.config.MqttServer)
	opts.OnConnect = func(client mqtt.Client) {
		// ok we connected. Now try to set our subscriptions
		tok := handler.conn.Subscribe("/tool/#", 0, handler.cbMessage)
		if tok.Wait() && tok.Error() != nil {
			log.Err(tok.Error()).
				Str("Topic", "/tool/#").
				Msg("Error adding subscription")
		}
		tok = handler.conn.Subscribe("/door/#", 0, handler.cbMessage)
		if tok.Wait() && tok.Error() != nil {
			log.Err(tok.Error()).
				Str("Topic", "/door/#").
				Msg("Error adding subscription")
		}
		log.Info().Msg("MQTT Subscriptions set up")
	}

	handler.conn = mqtt.NewClient(opts)

	go handler.handleMqtt()
}
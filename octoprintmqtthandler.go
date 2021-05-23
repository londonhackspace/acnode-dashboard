package main

import (
	"encoding/json"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/londonhackspace/acnode-dashboard/acnode"
	"github.com/rs/zerolog/log"
	"strings"
)

var (
	// list of suffixes of topics octoprint produces
	suffixes = []string{
		"event/Disconnected",
		"event/Connected",

		"event/PrinterStateChanged",
		"event/FirmwareData",
		"event/ZChange",
		"event/plugin_pi_support_throttle_state",
		"temperature/tool0",
		"temperature/bed",
	}
)

type octoprintMessage struct {
	Event string `json:"_event"`
	Timestamp uint64 `json:"_timestamp"`
}

type octoprintPrinterStateChange struct {
	octoprintMessage

	StateString string `json:"state_string"`
	StateId string `json:"state_id"`
}

type firmwareVersion struct {
	FirmwareVersion string `json:"FIRMWARE_VERSION"`
}

type octoprintFirmwareVersion struct {
	octoprintMessage

	Data firmwareVersion `json:"data"`
	Name string `json:"name"`
}

type octoprintChangeMessage struct {
	octoprintMessage

	Old float32 `json:"old"`
	New float32 `json:"new"`
}

type octoprintTemperatureMessage struct {
	Timestamp uint64 `json:"_timestamp"`

	Actual float32 `json:"actual"`
	Target float32 `json:"target"`
}

type octoprintThrottleStateMessage struct {
	octoprintMessage

	Overheat bool `json:"current_overheat"`
	Undervoltage bool `json:"current_undervoltage"`
}

// Is this an octoprint message we care about?
func isOctoprintTopic(topic string) bool {
	for _,suffix := range suffixes {
		if strings.HasSuffix(topic, suffix) {
			return true
		}
	}

	return false
}

func handleOctoprintMessage(node acnode.ACNode, message mqtt.Message) {
	rec := node.GetPrinterStatus()

	if strings.HasSuffix(message.Topic(), "event/Disconnected") {
		msg := octoprintMessage{}
		err := json.Unmarshal(message.Payload(), &msg)
		if err != nil {
			log.Err(err).Str("node", node.GetName()).Msg("Error unmarshalling event/Disconnect message")
		}
		if msg.Timestamp > rec.OctoprintConnectedTimestamp {
			log.Info().Str("node", node.GetName()).Msg("Received updated Octoprint connection status")
			rec.OctoprintConnected = false
			rec.OctoprintConnectedTimestamp = msg.Timestamp
		}
	} else if strings.HasSuffix(message.Topic(), "event/Connected") {
		msg := octoprintMessage{}
		err := json.Unmarshal(message.Payload(), &msg)
		if err != nil {
			log.Err(err).Str("node", node.GetName()).Msg("Error unmarshalling event/Connected message")
		}
		if msg.Timestamp > rec.OctoprintConnectedTimestamp {
			log.Info().Str("node", node.GetName()).Msg("Received updated Octoprint connection status")
			rec.OctoprintConnected = true
			rec.OctoprintConnectedTimestamp = msg.Timestamp
		}
	} else if strings.HasSuffix(message.Topic(), "event/PrinterStateChanged") {
		msg := octoprintPrinterStateChange{}
		err := json.Unmarshal(message.Payload(), &msg)
		if err != nil {
			log.Err(err).Str("node", node.GetName()).Msg("Error unmarshalling event/PrinterStateChanged message")
		}
		if msg.Timestamp > rec.PrinterStatusTimestamp {
			log.Info().Str("node", node.GetName()).Msg("Received updated printer status")
			rec.PrinterStatus = msg.StateString
			rec.PrinterStatusTimestamp = msg.Timestamp
		}
	} else if strings.HasSuffix(message.Topic(), "event/FirmwareData") {
		msg := octoprintFirmwareVersion{}
		err := json.Unmarshal(message.Payload(), &msg)
		if err != nil {
			log.Err(err).Str("node", node.GetName()).Msg("Error unmarshalling event/FirmwareData message")
		}
		if msg.Timestamp > rec.FirmwareVersionTimestamp {
			log.Info().Str("node", node.GetName()).Msg("Received updated firmware info")
			rec.FirmwareVersion = msg.Data.FirmwareVersion
			rec.FirmwareVersionTimestamp = msg.Timestamp
		}
	} else if strings.HasSuffix(message.Topic(), "event/ZChange") {
		msg := octoprintChangeMessage{}
		err := json.Unmarshal(message.Payload(), &msg)
		if err != nil {
			log.Err(err).Str("node", node.GetName()).Msg("Error unmarshalling event/ZChange message")
		}
		if msg.Timestamp > rec.ZHeightTimestamp {
			log.Info().Str("node", node.GetName()).Msg("Received updated Z height")
			rec.ZHeight = msg.New
			rec.ZHeightTimestamp = msg.Timestamp
		}
	} else if strings.HasSuffix(message.Topic(), "event/plugin_pi_support_throttle_state") {
		msg := octoprintThrottleStateMessage{}
		err := json.Unmarshal(message.Payload(), &msg)
		if err != nil {
			log.Err(err).Str("node", node.GetName()).Msg("Error unmarshalling event/plugin_pi_support_throttle_state message")
		}
		if msg.Timestamp > rec.PiThrottleTimestamp {
			log.Info().Str("node", node.GetName()).Msg("Received updated Pi throttle status")
			rec.PiUndervoltage = msg.Undervoltage
			rec.PiOverheat = msg.Overheat
			rec.PiThrottleTimestamp = msg.Timestamp
		}
	} else if strings.HasSuffix(message.Topic(), "temperature/tool0") {
		msg := octoprintTemperatureMessage{}
		err := json.Unmarshal(message.Payload(), &msg)
		if err != nil {
			log.Err(err).Str("node", node.GetName()).Msg("Error unmarshalling temperature/tool0 message")
		}
		if msg.Timestamp > rec.HotendTemperatureTimestamp {
			log.Info().Str("node", node.GetName()).Msg("Received updated hotend temperature")
			rec.HotendTemperature = msg.Actual
		}
	} else if strings.HasSuffix(message.Topic(), "temperature/bed") {
		msg := octoprintTemperatureMessage{}
		err := json.Unmarshal(message.Payload(), &msg)
		if err != nil {
			log.Err(err).Str("node", node.GetName()).Msg("Error unmarshalling temperature/bed message")
		}
		if msg.Timestamp > rec.BedTemperatureTimestamp {
			log.Info().Str("node", node.GetName()).Msg("Received updated bed temperature")
			rec.BedTemperature = msg.Actual
			rec.BedTemperatureTimestamp = msg.Timestamp
		}
	}
}
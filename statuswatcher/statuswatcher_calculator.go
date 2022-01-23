package main

import (
	"github.com/londonhackspace/acnode-dashboard/apitypes"
	"time"
)

func (sw *StatusWatcher) calculateNodeState(node *apitypes.ACNode) (int, []string) {
	var healthHints []string
	health := STATE_GOOD

	if !node.InService {
		healthHints = append(healthHints, "Tool out of service")
	}

	// Calculate relative values for decisions
	lastSeen := time.Now().Sub(time.Unix(node.LastSeen, 0)).Seconds()
	lastSeenApi := time.Now().Sub(time.Unix(node.LastSeenAPI, 0)).Seconds()
	lastSeenMQTT := time.Now().Sub(time.Unix(node.LastSeenMQTT, 0)).Seconds()
	lastStarted := time.Now().Sub(time.Unix(node.LastStarted, 0)).Seconds()

	// how often should a node check in with the API?
	var apiThreshold float64 = 610
	apiThresholdText := "over 10 minutes"

	// how often should a node check in with MQTT?
	var mqttThreshold float64 = 130
	mqttThresholdText := "over 2 minutes"

	// unrestricted doors don't check in nearly so often if they're not
	// running firmware new enough to periodically revalidate the cache
	// or aren't running a maintainer cache
	if node.Type == "Door" {
		apiThreshold = (3600 * 12) + 10
		apiThresholdText = "over 12 hours"
	}

	// Only consider time if the node isn't flagged as transient
	if !node.IsTransient {
		if node.LastSeenMQTT > -1 || node.LastSeenAPI > -1 {
			// if we're seeing neither MQTT or ACServer log entries,
			// it's probably dead
			if (lastSeenApi > apiThreshold || node.LastSeenAPI == -1) &&
				(lastSeenMQTT > mqttThreshold || node.LastSeenMQTT == -1) {
				text := mqttThresholdText
				if apiThreshold > mqttThreshold {
					text = apiThresholdText
				}
				healthHints = append(healthHints, "Has not been seen online in any form in "+text)

				return STATE_BAD, healthHints
			}

			if node.LastSeenMQTT == -1 || lastSeenMQTT > mqttThreshold {
				healthHints = append(healthHints, "Has not sent a message via MQTT in over "+mqttThresholdText)
				health = STATE_WARN
			}

			if (node.LastSeenAPI == -1 || lastSeenApi > apiThreshold) && !node.InUse {
				healthHints = append(healthHints, "Has not contacted ACServer in over "+apiThresholdText)
				health = STATE_WARN
			}

		} else {
			// I'm not actually convinced we need this now, since I think all nodes have API and MQTT
			// values separately, but since we're porting the TS logic...
			if node.LastSeen == -1 {
				healthHints = append(healthHints, "Has never been seen online")
				return STATE_UNKNOWN, healthHints
			}

			if lastSeen > 610 {
				healthHints = append(healthHints, "Has not been seen online in over 10 minutes")
				health = STATE_BAD
			} else if lastSeen > 70 {
				healthHints = append(healthHints, "Has not been seen online in over a minute")
				return STATE_WARN, healthHints
			}
		}
	}

	// lower the health if the node watchdog'd recently
	if node.LastStarted > 0 && lastStarted < 610 {
		if node.ResetCause == "Watchdog" {
			healthHints = append(healthHints, "Watchdog reset detected in last 10 minutes")
			health = STATE_WARN
		}
	}

	// low on memory?
	totalMem := node.MemUsed + node.MemFree
	if node.MemUsed > 0 && node.MemFree < (totalMem/10) {
		healthHints = append(healthHints, "Very low on memory (<10% left)")
		health = STATE_BAD
	}

	return health, healthHints
}

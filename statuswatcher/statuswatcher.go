package main

import (
	"context"
	"github.com/londonhackspace/acnode-dashboard/apiclient"
	"github.com/londonhackspace/acnode-dashboard/apitypes"
	"github.com/londonhackspace/acnode-dashboard/statuswatcher/slackapi"
	"github.com/rs/zerolog/log"
	"time"
)

type StatusWatcher struct {
	client apiclient.APIClient
	slack slackapi.SlackAPI
	slackChannel string

	previousStates map[string]int
}

const (
	STATE_BAD = iota
	STATE_UNKNOWN = iota
	STATE_WARN = iota
	STATE_GOOD = iota
)

func getStringFromNodeState(state int) string {
	switch(state) {
	case STATE_BAD:
		return "BAD"
	case STATE_UNKNOWN:
		return "UNKNOWN"
	case STATE_WARN:
		return "WARN"
	case STATE_GOOD:
		return "GOOD"
	}
	return "<Unknown>"
}

func CreateStatusWatcher(client apiclient.APIClient, slack slackapi.SlackAPI, slackChannel string) *StatusWatcher {
	return &StatusWatcher{
		client: client,
		previousStates: map[string]int{},
		slack: slack,
		slackChannel: slackChannel,
	}
}

func (sw *StatusWatcher) Run(ctx context.Context) {
	ticker := time.NewTicker(time.Second * 10)
	log.Info().Msg("Node Watcher Running")
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			sw.runWatch(ctx)
		}
	}
	log.Info().Msg("Node Watcher Exiting")
}

func (sw *StatusWatcher) postSlackMessage(node *apitypes.ACNode, hints []string, oldState int, newState int) {
	if sw.slack != nil{
		action := "healthier"

		if oldState > newState {
			action = "less healthy"
		}

		msg := "Node " + node.MqttName + " has become " + action + "\n" +
			"New health: " + getStringFromNodeState(newState) + "\n" +
			"Hints: "

		for i, h := range hints {
			if i > 0 {
				msg += ", "
			}
			msg += h
		}

		if len(hints) == 0 {
			msg += "(None)"
		}
		log.Debug().Str("channel", sw.slackChannel).
			Msg("Posting message to Slack")
		sw.slack.PostMessage(msg, sw.slackChannel)
	}
}

func (sw *StatusWatcher) checkNode(name string) {
	status, err := sw.client.GetNode(name)
	if err != nil {
		log.Err(err).Str("node", name).Msg("Error getting node status")
		return
	}

	newState,hints := sw.calculateNodeState(status)

	oldState, ok := sw.previousStates[status.MqttName]
	if ok {
		if newState > oldState {
			log.Info().Str("node", status.MqttName).
				Str("State", getStringFromNodeState(newState)).
				Msg("Node has improved")
			sw.postSlackMessage(status, hints, oldState, newState)
		} else if newState < oldState {
			log.Info().Str("node", status.MqttName).
				Str("State", getStringFromNodeState(newState)).
				Msg("Node has degraded")
			sw.postSlackMessage(status, hints, oldState, newState)
		}
	} else {
		log.Info().Str("node", status.MqttName).
			Str("State", getStringFromNodeState(newState)).
			Msg("New node has appeared")
		if newState < STATE_GOOD {
			combined := ""
			for i,s := range hints {
				if i > 0 {
					combined += ", "
				}
				combined += s
			}
			log.Warn().Str("node", status.MqttName).
				Str("status", combined).
				Msg("Dubious health status")
		}
	}

	sw.previousStates[status.MqttName] = newState
}

func (sw *StatusWatcher) runWatch(ctx context.Context) {
	log.Info().Msg("Checking Status of Nodes")
	nodes, err := sw.client.GetNodes()
	if err != nil {
		log.Err(err).Msg("Error getting nodes")
		return
	}

	for _,node := range nodes {
		sw.checkNode(node)
	}
}

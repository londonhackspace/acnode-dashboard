StatusWatcher
===
The StatusWatcher is a tool to monitor the state of nodes, and raise alerts when things change.
It uses the same API as the frontend, and can be seen as a port of the logic within the frontend into a background monitoring tool

Config
---
Configuration is via environment variables
* API_KEY - the ACNode-Dash API Key
* HTTP_LISTEN - the HTTP listener config. Defaults to ":8081"
* ACNODE_DASH - the address of acnode-dash - defaults to "https://acnode-dash.london.hackspace.org.uk/api/"
* SLACK_TOKEN - a Slack API key - defaults to an empty string, disabling Slack functionality
* SLACK_CHANNEL - a Slack channel to post to - defaults to "bot-test"

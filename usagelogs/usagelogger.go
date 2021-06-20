package usagelogs

import "github.com/londonhackspace/acnode-dashboard/acnode"

type UsageLogger interface {
	AddUsageLog(node *acnode.ACNode, msg acnode.Announcement)

	GetUsageLogNodes() []string

	GetUsageLogCountForNode(node string) int64
	GetUsageLogsForNode(node string, from int64, to int64) []LogEntry
}

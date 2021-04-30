package usagelogs

import "github.com/londonhackspace/acnode-dashboard/acnode"

type UsageLogger interface {
	AddUsageLog(node *acnode.ACNode, msg acnode.Announcement)
}

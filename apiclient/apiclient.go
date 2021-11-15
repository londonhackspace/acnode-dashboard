package apiclient

import "github.com/londonhackspace/acnode-dashboard/apitypes"

type APIClient interface {
	GetNodes() ([]string, error)
	GetNode(name string) (*apitypes.ACNode, error)
}

package software

import "github.com/mittwald/spacectl/client/lowlevel"

type SoftwareClient interface {
	List() ([]Software, error)
	ListWithVersions() ([]Software, error)
	Get(softwareID string) (*Software, error)
}

func NewSoftwareClient(client *lowlevel.SpacesLowlevelClient, group string) SoftwareClient {
	return &softwareClient{client, group}
}

type softwareClient struct {
	client *lowlevel.SpacesLowlevelClient
	group  string
}

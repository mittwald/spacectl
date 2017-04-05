package teams

import "github.com/mittwald/spacectl/client/lowlevel"

type TeamsClient interface {
	List() ([]Team, error)
	Create(string) (*Team, error)
}

func NewTeamsClient(client *lowlevel.SpacesLowlevelClient) (TeamsClient) {
	return &teamsClient{client}
}

type teamsClient struct {
	client *lowlevel.SpacesLowlevelClient
}
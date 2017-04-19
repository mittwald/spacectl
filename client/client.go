package client

import (
	"fmt"
	"github.com/mittwald/spacectl/client/lowlevel"
	"github.com/mittwald/spacectl/client/teams"
	"github.com/spf13/viper"
)

type SpacesClient interface {
	Teams() teams.TeamsClient
}

type spacesClient struct {
	client *lowlevel.SpacesLowlevelClient
}

func NewSpacesClientAutoConf() (SpacesClient, error) {
	token := viper.GetString("token")

	lowlevelClient, err := lowlevel.NewSpacesLowlevelClient(token, viper.GetString("apiServer"))
	if err != nil {
		return nil, fmt.Errorf("could not create SPACES client: %s", err)
	}

	return &spacesClient{
		lowlevelClient,
	}, nil
}

func (c *spacesClient) Teams() teams.TeamsClient {
	return teams.NewTeamsClient(c.client)
}
package spaces

import (
	"fmt"
	"github.com/mittwald/spacectl/client/teams"
	"errors"
)

func (c *spacesClient) List() ([]Space, error) {
	var spaces []Space

	err := c.client.Get("/spaces/", &spaces)
	if err != nil {
		return nil, fmt.Errorf("could not load spaces: %s", err)
	}

	return spaces, nil
}

func (c *spacesClient) ListByTeam(teamID string) ([]Space, error) {
	var spaces []Space
	var team teams.Team

	if teamID == "" {
		return nil, errors.New("team ID must not be empty")
	}

	err := c.client.Get("/teams/" + teamID, &team)
	if err != nil {
		return nil, err
	}

	link, err := team.Links.GetLinkByRel("spaces")
	if err != nil {
		return nil, err
	}

	err = link.Get(c.client, &spaces)
	if err != nil {
		return nil, fmt.Errorf("could not load spaces: %s", err)
	}

	return spaces, nil
}
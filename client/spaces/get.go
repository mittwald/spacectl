package spaces

import (
	"fmt"
	"github.com/mittwald/spacectl/client/teams"
	"github.com/mittwald/spacectl/client/errors"
)

func (c *spacesClient) GetByID(spaceID string) (*Space, error) {
	var space Space

	err := c.client.Get("/spaces/" + spaceID, &space)
	if err != nil {
		return nil, fmt.Errorf("could not load space %s: %s", spaceID, err)
	}

	return &space, nil
}

func (c *spacesClient) GetByTeamAndName(teamIDOrName string, spaceIDOrName string) (*Space, error) {
	var team teams.Team
	var space Space

	err := c.client.Get("/teams/" + teamIDOrName, &team)
	if err != nil {
		return nil, fmt.Errorf("could not load team %s: %s", teamIDOrName, err)
	}

	link, err := team.Links.GetLinkByRel("space")
	if err != nil {
		return nil, errors.ErrUnauthorized{Msg: "not authorized to access team's Spaces", Inner: err}
	}

	err = link.WithParam("id", spaceIDOrName).Get(c.client, &space)
	if err != nil {
		return nil, fmt.Errorf("could not load space %s: %s", spaceIDOrName, err)
	}

	return &space, nil
}
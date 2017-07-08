package spaces

import (
	"github.com/mittwald/spacectl/client/teams"
	"github.com/mittwald/spacectl/client/lowlevel"
	"github.com/mittwald/spacectl/client/errors"
)

func (c *spacesClient) Declare(teamID string, declaration *SpaceDeclaration) (*Space, error) {
	var team teams.Team
	var existingSpaces []Space

	err := c.client.Get("/teams/" + teamID, &team)
	if err != nil {
		return nil, err
	}

	c.logger.Printf("Space '%s' is declared in team %s", declaration.Name.DNSName, team.ID)

	l, err := team.Links.GetLinkByRel("spaces")
	if err != nil {
		return nil, err
	}

	err = l.Get(c.client, &existingSpaces)
	if err != nil {
		return nil, err
	}

	var existing *Space
	var created Space

	for i := range existingSpaces {
		if existingSpaces[i].Name.DNSName == declaration.Name.DNSName {
			existing = &existingSpaces[i]
		}
	}

	if existing == nil {
		c.logger.Printf("Space '%s' does not yet exist", declaration.Name.DNSName)

		l, err := team.Links.GetLinkByRel("spaces")
		if err != nil {
			switch err.(type) {
			case lowlevel.ErrLinkNotFound:
				return nil, errors.ErrUnauthorized{Msg: "Not authorized to create Spaces in this team", Inner: err}
			default:
				return nil, err
			}
		}

		err = l.Post(c.client, declaration, &created)
		if err != nil {
			return nil, errors.ErrNested{Msg: "Error occurred while creating a new Space", Inner: err}
		}

		return &created, nil
	}

	c.logger.Printf("Space '%s' already exists with ID %s", declaration.Name.DNSName, existing.ID)

	err = c.client.Put(existing.HREF, declaration, &created)
	if err != nil {
		return nil, errors.ErrNested{Msg: "Error occurred while updating an existing Space", Inner: err}
	}

	return nil, nil
}

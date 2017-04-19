package teams

import "fmt"

func (c *teamsClient) List() ([]Team, error) {
	var teams []Team

	err := c.client.Get("/teams", &teams)
	if err != nil {
		return nil, fmt.Errorf("could not load teams: %s", err)
	}

	return teams, nil
}
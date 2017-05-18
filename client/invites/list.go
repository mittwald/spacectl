package invites

import "fmt"

func (c *invitesClient) List() ([]Invite, error) {
	var invites []Invite

	err := c.client.Get("/invites", &invites)
	if err != nil {
		return nil, fmt.Errorf("could not load invites: %s", err)
	}

	return invites, nil
}
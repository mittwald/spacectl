package invites

import "fmt"

func (c *invitesClient) ListIncoming() ([]Invite, error) {
	var invites []Invite

	err := c.client.Get("/invites/incoming", &invites)
	if err != nil {
		return nil, fmt.Errorf("could not load invites: %s", err)
	}

	return invites, nil
}

func (c *invitesClient) ListOutgoing() ([]Invite, error) {
	var invites []Invite

	err := c.client.Get("/invites/outgoing", &invites)
	if err != nil {
		return nil, fmt.Errorf("could not load invites: %s", err)
	}

	return invites, nil
}

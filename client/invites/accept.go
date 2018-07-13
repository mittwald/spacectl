package invites

import "fmt"

func (c *invitesClient) Accept(inviteID string) error {
	var invite Invite

	err := c.client.Get("/invites/"+inviteID, &invite)
	if err != nil {
		return fmt.Errorf("could not load invite %s: %s", inviteID, err)
	}

	action, err := invite.Actions.GetLinkByRel("accept")
	if err != nil {
		return fmt.Errorf("not authorized to accept invite %s", inviteID)
	}

	err = action.Post(c.client, nil, nil)
	if err != nil {
		return fmt.Errorf("error while accepting the invite: %s", err)
	}

	return nil
}

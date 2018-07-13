package invites

import "fmt"

func (c *invitesClient) Revoke(inviteID string) error {
	var invite Invite

	err := c.client.Get("/invites/"+inviteID, &invite)
	if err != nil {
		return fmt.Errorf("could not load invite %s: %s", inviteID, err)
	}

	action, err := invite.Actions.GetLinkByRel("revoke")
	if err != nil {
		return fmt.Errorf("not authorized to revoke invite %s", inviteID)
	}

	err = action.Post(c.client, nil, nil)
	if err != nil {
		return fmt.Errorf("error while revoking the invite: %s", err)
	}

	return nil
}

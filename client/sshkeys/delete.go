package sshkeys

import (
	"errors"
	"fmt"

	"github.com/mittwald/spacectl/client/lowlevel"
)

func (c *sshKeyClient) Delete(keyID string) error {
	response := lowlevel.Message{}
	profile := lowlevel.Linkeable{}

	err := c.client.Get("/users/me", &profile)
	if err != nil {
		return fmt.Errorf("Error while loading user profile: %s", err)
	}

	link, err := profile.Links.GetLinkByRel("key")
	if err != nil {
		return errors.New("You are not authorized to manage SSH public keys")
	}

	err = link.WithParam("keyID", keyID).Delete(c.client, &response)
	if err != nil {
		return err
	}

	return nil
}

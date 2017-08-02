package sshkeys

import (
	"github.com/mittwald/spacectl/client/lowlevel"
	"fmt"
	"errors"
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

	err = link.WithParam("id", keyID).Delete(c.client, &response)
	if err != nil {
		return err
	}

	return nil
}
package sshkeys

import (
	"errors"
	"fmt"
	"github.com/mittwald/spacectl/client/lowlevel"
)

func (c *sshKeyClient) List() ([]SSHKey, error) {
	user := lowlevel.Linkeable{}

	err := c.client.Get("/users/me", &user)
	if err != nil {
		return nil, fmt.Errorf("Error while loading user profile: %s", err)
	}

	link, err := user.Links.GetLinkByRel("keys")
	if err != nil {
		return nil, errors.New("You are not authorized to manage SSH public keys")
	}

	keys := make([]SSHKey, 0)
	err = link.Get(c.client, &keys)

	if err != nil {
		return nil, err
	}

	return keys, nil
}

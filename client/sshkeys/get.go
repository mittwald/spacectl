package sshkeys

import (
	"errors"
	"fmt"

	"github.com/mittwald/spacectl/client/lowlevel"
)

func (c *sshKeyClient) Get(keyID string) (*SSHKey, error) {
	user := lowlevel.Linkeable{}

	err := c.client.Get("/users/me", &user)
	if err != nil {
		return nil, fmt.Errorf("Error while loading user profile: %s", err)
	}

	link, err := user.Links.GetLinkByRel("key")
	if err != nil {
		return nil, errors.New("You are not authorized to manage SSH public keys")
	}

	key := SSHKey{}
	err = link.WithParam("keyID", keyID).Get(c.client, &key)

	if err != nil {
		return nil, err
	}

	return &key, nil
}

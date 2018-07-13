package sshkeys

import (
	"errors"
	"fmt"
	"github.com/mittwald/spacectl/client/lowlevel"
)

func (c *sshKeyClient) Add(keyBytes []byte, algorithm, comment string) (*SSHKey, error) {
	resultkey := SSHKey{}
	key := sshKeyInput{
		CipherAlgorithm: algorithm,
		Comment:         comment,
		Key:             keyBytes,
	}

	profile := lowlevel.Linkeable{}
	err := c.client.Get("/users/me", &profile)
	if err != nil {
		return nil, fmt.Errorf("Error while loading user profile: %s", err)
	}

	link, err := profile.Links.GetLinkByRel("keys")
	if err != nil {
		return nil, errors.New("You are not authorized to manage SSH public keys")
	}

	err = link.Post(c.client, &key, &resultkey)
	if err != nil {
		return nil, err
	}

	return &resultkey, err
}

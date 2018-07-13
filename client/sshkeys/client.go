package sshkeys

import (
	"github.com/mittwald/spacectl/client/lowlevel"
	"log"
)

type SSHKeyClient interface {
	List() ([]SSHKey, error)
	Get(keyID string) (*SSHKey, error)
	Add(key []byte, algorithm, comment string) (*SSHKey, error)
	Delete(keyID string) error
}

func NewSSHKeyClient(client *lowlevel.SpacesLowlevelClient, logger *log.Logger) SSHKeyClient {
	return &sshKeyClient{client, logger}
}

type sshKeyClient struct {
	client *lowlevel.SpacesLowlevelClient
	logger *log.Logger
}

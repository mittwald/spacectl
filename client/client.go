package client

import (
	"fmt"
	"github.com/mittwald/spacectl/client/backups"
	"github.com/mittwald/spacectl/client/invites"
	"github.com/mittwald/spacectl/client/lowlevel"
	"github.com/mittwald/spacectl/client/payment"
	"github.com/mittwald/spacectl/client/software"
	"github.com/mittwald/spacectl/client/spaces"
	"github.com/mittwald/spacectl/client/sshkeys"
	"github.com/mittwald/spacectl/client/teams"
	"io/ioutil"
	"log"
)

type SpacesClient interface {
	Teams() teams.TeamsClient
	Invites() invites.InvitesClient
	Spaces() spaces.SpacesClient
	Backups() backups.BackupClient
	Recoveries() backups.RecoveryClient
	SSHKeys() sshkeys.SSHKeyClient
	Applications() software.SoftwareClient
	Databases() software.SoftwareClient
	Payment() payment.Client
}

type SpacesClientConfig struct {
	Token     string
	APIServer string
	Logger    *log.Logger
}

type spacesClient struct {
	client *lowlevel.SpacesLowlevelClient
	logger *log.Logger
}

func NewSpacesClient(config SpacesClientConfig) (SpacesClient, error) {
	if config.Logger == nil {
		config.Logger = log.New(ioutil.Discard, "", 0)
	}

	lowlevelClient, err := lowlevel.NewSpacesLowlevelClient(config.Token, config.APIServer, config.Logger)
	if err != nil {
		return nil, fmt.Errorf("could not create SPACES client: %s", err)
	}

	return &spacesClient{
		lowlevelClient,
		config.Logger,
	}, nil
}

func (c *spacesClient) Teams() teams.TeamsClient {
	return teams.NewTeamsClient(c.client)
}

func (c *spacesClient) Invites() invites.InvitesClient {
	return invites.NewInvitesClient(c.client)
}

func (c *spacesClient) Spaces() spaces.SpacesClient {
	return spaces.NewSpacesClient(c.client, c.logger)
}

func (c *spacesClient) Backups() backups.BackupClient {
	return backups.NewBackupClient(c.client, c.logger)
}

func (c *spacesClient) Recoveries() backups.RecoveryClient {
	return backups.NewRecoveryClient(c.client, c.logger)
}

func (c *spacesClient) SSHKeys() sshkeys.SSHKeyClient {
	return sshkeys.NewSSHKeyClient(c.client, c.logger)
}

func (c *spacesClient) Applications() software.SoftwareClient {
	return software.NewSoftwareClient(c.client, "applications")
}

func (c *spacesClient) Databases() software.SoftwareClient {
	return software.NewSoftwareClient(c.client, "databases")
}

func (c *spacesClient) Payment() payment.Client {
	return payment.NewClient(c.client)
}

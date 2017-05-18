package invites

import "github.com/mittwald/spacectl/client/lowlevel"

type InvitesClient interface {
	List() ([]Invite, error)
	//Accept(inviteID string) (error)
	//Decline(inviteID string) (error)
	//Revoke(inviteID string) (error)
}

func NewInvitesClient(client *lowlevel.SpacesLowlevelClient) (InvitesClient) {
	return &invitesClient{client}
}

type invitesClient struct {
	client *lowlevel.SpacesLowlevelClient
}
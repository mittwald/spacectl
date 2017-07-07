package spaces

import "github.com/mittwald/spacectl/client/lowlevel"

type SpacesClient interface {
	List() ([]Space, error)
	ListByTeam(teamID string) ([]Space, error)
	//Create(string, string) (Space, error)
	//InviteByEmail(teamID string, email string, message string, role string) (Invite, error)
	//ListMembers(teamID string) ([]Membership, error)
}

func NewSpacesClient(client *lowlevel.SpacesLowlevelClient) (SpacesClient) {
	return &spacesClient{client}
}

type spacesClient struct {
	client *lowlevel.SpacesLowlevelClient
}
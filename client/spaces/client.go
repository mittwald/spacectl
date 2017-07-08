package spaces

import (
	"github.com/mittwald/spacectl/client/lowlevel"
	"log"
)

type SpacesClient interface {
	List() ([]Space, error)
	ListByTeam(teamID string) ([]Space, error)
	Declare(teamID string, declaration *SpaceDeclaration) (*Space, error)
	//Create(string, string) (Space, error)
	//InviteByEmail(teamID string, email string, message string, role string) (Invite, error)
	//ListMembers(teamID string) ([]Membership, error)
}

func NewSpacesClient(client *lowlevel.SpacesLowlevelClient, logger *log.Logger) (SpacesClient) {
	return &spacesClient{client, logger}
}

type spacesClient struct {
	client *lowlevel.SpacesLowlevelClient
	logger *log.Logger
}
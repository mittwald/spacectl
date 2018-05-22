package spaces

import (
	"github.com/mittwald/spacectl/client/lowlevel"
	"log"
)

type SpacesClient interface {
	List() ([]Space, error)
	ListByTeam(teamID string) ([]Space, error)
	Declare(teamID string, declaration *SpaceDeclaration) (*Space, error)
	GetByID(spaceID string) (*Space, error)
	GetByTeamAndName(teamIDOrName string, spaceIDOrName string) (*Space, error)
	Delete(spaceID string) (error)
	UpdateApplication(spaceID, stage, targetStage, version string) (*ApplicationUpdate, error)
	ListApplicationUpdatesByStage(spaceID, stage string) ([]ApplicationUpdate, error)
	ListApplicationUpdatesBySpace(spaceID string) ([]ApplicationUpdate, error)
}

func NewSpacesClient(client *lowlevel.SpacesLowlevelClient, logger *log.Logger) (SpacesClient) {
	return &spacesClient{client, logger}
}

type spacesClient struct {
	client *lowlevel.SpacesLowlevelClient
	logger *log.Logger
}
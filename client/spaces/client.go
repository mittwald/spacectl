package spaces

import (
	"log"

	"github.com/mittwald/spacectl/client/lowlevel"
)

type SpacesClient interface {
	List() ([]Space, error)
	ListByTeam(teamID string) ([]Space, error)
	Declare(teamID string, declaration *SpaceDeclaration) (*Space, error)
	GetByID(spaceID string) (*Space, error)
	GetByTeamAndName(teamIDOrName string, spaceIDOrName string) (*Space, error)
	Delete(spaceID string) error
	UpdateApplication(spaceID, stage, targetStage, version string) (*ApplicationUpdate, error)
	ListApplicationUpdatesByStage(spaceID, stage string) ([]ApplicationUpdate, error)
	ListApplicationUpdatesBySpace(spaceID string) ([]ApplicationUpdate, error)
	GetPaymentLink(spaceID string) (*SpacePaymentLink, error)
	ConnectWithPaymentProfile(spaceID string, paymentProfileID string, planID string) (*SpacePaymentLink, error)
	ListCaughtEmails(spaceID string) (CaughtEmailList, error)
	GetComputeMetrics(spaceID string, scope string) (ComputeMetricPointList, error)
	ListVirtualHostsByStage(spaceID, stage string) (VirtualHostList, error)
	UpdateVirtualHost(spaceID, stage string, vhost VirtualHost) (*VirtualHost, error)
	DeleteVirtualHost(spaceID, stage, hostname string) error
	GetStageProtection(spaceID, stage string) (*StageProtection, error)
	CreateStageProtection(spaceID, stage string, inputProtection StageProtection) (*StageProtection, error)
	DeleteStageProtection(spaceID, stage string) error
	ListDatabaseUsersByStage(spaceID, stage string) (DatabaseUserList, error)
	GetDatabaseUser(spaceID, stage, username string) (*DatabaseUser, error)
	CreateDatabaseUser(spaceID, stage string, dbUser DatabaseUserInput) (*DatabaseUser, error)
	DeleteDatabaseUser(spaceID, stage, name, dbType string) error
}

func NewSpacesClient(client *lowlevel.SpacesLowlevelClient, logger *log.Logger) SpacesClient {
	return &spacesClient{client, logger}
}

type spacesClient struct {
	client *lowlevel.SpacesLowlevelClient
	logger *log.Logger
}

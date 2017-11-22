package backups

import (
	"fmt"
	"github.com/mittwald/spacectl/client/errors"
	"github.com/mittwald/spacectl/client/spaces"
)

type createBackupRequest struct {
	Keep        bool   `json:"keep"`
	Description string `json:"description,omitempty"`
}

func (c *backupClient) Create(spaceID string, stageName string, keep bool, description string) (*Backup, error) {
	stageURL := fmt.Sprintf("/spaces/%s/stages/%s", spaceID, stageName)
	stage := spaces.Stage{}

	err := c.client.Get(stageURL, &stage)
	if err != nil {
		return nil, errors.ErrNested{Inner: err, Msg: "could not access stage"}
	}

	backupsLink, err := stage.Links.GetLinkByRel("backups")
	if err != nil {
		return nil, errors.ErrUnauthorized{Inner: err, Msg: "backups for this stage are not accessible"}
	}

	res := Backup{}
	req := createBackupRequest{
		Keep:        keep,
		Description: description,
	}

	err = backupsLink.Post(c.client, req, &res)
	if err != nil {
		return nil, errors.ErrNested{Inner: err, Msg: "could not create backup"}
	}

	return &res, nil
}

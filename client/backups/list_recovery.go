package backups

import (
	"fmt"
	"github.com/mittwald/spacectl/client/errors"
	"github.com/mittwald/spacectl/client/spaces"
)

func (c *recoveryClient) ListForSpace(spaceID string) ([]Recovery, error) {
	spaceURL := fmt.Sprintf("/spaces/%s", spaceID)
	space := spaces.Space{}

	err := c.client.Get(spaceURL, &space)
	if err != nil {
		return nil, errors.ErrNested{Inner: err, Msg: "could not load space"}
	}

	recoveryLink, err := space.Links.GetLinkByRel("recoveries")
	recoveries := make([]Recovery, 0)

	if err != nil {
		return nil, errors.ErrUnauthorized{Inner: err, Msg: "recoveries are not accessible for this backup"}
	}

	err = recoveryLink.Get(c.client, &recoveries)
	if err != nil {
		return nil, errors.ErrNested{Inner: err, Msg: "could not load recoveries"}
	}

	return recoveries, nil
}

func (c *recoveryClient) ListForBackupID(backupID string) ([]Recovery, error) {
	backupURL := fmt.Sprintf("/backups/%s", backupID)
	backup := Backup{}

	err := c.client.Get(backupURL, &backup)
	if err != nil {
		return nil, errors.ErrNested{Inner: err, Msg: "could not load backup"}
	}

	return c.ListForBackup(&backup)
}

func (c *recoveryClient) ListForBackup(backup *Backup) ([]Recovery, error) {
	recoveryLink, err := backup.Links.GetLinkByRel("recoveries")
	recoveries := make([]Recovery, 0)

	if err != nil {
		return nil, errors.ErrUnauthorized{Inner: err, Msg: "recoveries are not accessible for this backup"}
	}

	err = recoveryLink.Get(c.client, &recoveries)
	if err != nil {
		return nil, errors.ErrNested{Inner: err, Msg: "could not load recoveries"}
	}

	return recoveries, nil
}

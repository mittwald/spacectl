package backups

import (
	"fmt"
	"github.com/mittwald/spacectl/client/errors"
	"github.com/mittwald/spacectl/client/spaces"
	"net/url"
	"time"
)

func (o *BackupListOptions) toQuery() url.Values {
	if o == nil {
		return url.Values{}
	}

	q := url.Values{}

	if o.OnlyKeep {
		q.Set("keep", "1")
	}

	if !o.Since.IsZero() {
		q.Set("since", o.Since.Format(time.RFC3339))
	}

	return q
}

func (c *backupClient) ListForSpace(spaceID string, opts *BackupListOptions) ([]Backup, error) {
	space := spaces.Space{}
	spaceURL := fmt.Sprintf("/spaces/%s", spaceID)

	err := c.client.Get(spaceURL, &space)
	if err != nil {
		return nil, fmt.Errorf("could not load space %s: %s", spaceID, err)
	}

	backupLink, err := space.Links.GetLinkByRel("backups")
	if err != nil {
		return nil, errors.ErrUnauthorized{Inner: err, Msg: "backups are not accessible for this Space"}
	}

	backups := make([]Backup, 0)
	err = backupLink.GetWithQuery(opts.toQuery(), c.client, &backups)
	if err != nil {
		return nil, err
	}

	return backups, nil
}

func (c *backupClient) ListForStage(spaceID string, stageName string, opts *BackupListOptions) ([]Backup, error) {
	stage := spaces.Stage{}
	stageURL := fmt.Sprintf("/spaces/%s/stages/%s", spaceID, stageName)

	err := c.client.Get(stageURL, &stage)
	if err != nil {
		return nil, fmt.Errorf("could not load stage %s/%s: %s", spaceID, stageName, err)
	}

	backupLink, err := stage.Links.GetLinkByRel("backups")
	if err != nil {
		return nil, errors.ErrUnauthorized{Inner: err, Msg: "backups are not accessible for this Space"}
	}

	backups := make([]Backup, 0)
	err = backupLink.GetWithQuery(opts.toQuery(), c.client, &backups)
	if err != nil {
		return nil, err
	}

	return backups, nil
}

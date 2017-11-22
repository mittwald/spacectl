package backups

import "errors"

func (c *backupClient) Create(spaceID string, stage string, keep bool, description string) (*Backup, error) {
	// TODO
	return nil, errors.New("not implemented")
}
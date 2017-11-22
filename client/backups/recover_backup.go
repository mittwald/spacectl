package backups

import "errors"

func (c *backupClient) Recover(backupID string, files RecoverySpec, databases RecoverySpec) (*Recovery, error) {
	// TODO
	return nil, errors.New("not implemented")
}
package backups

import (
	"fmt"
)

func (c *backupClient) Delete(backupID string) error {
	backupURL := fmt.Sprintf("/backups/%s", backupID)
	return c.client.Delete(backupURL, nil)
}
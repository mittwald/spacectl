package backups

import "fmt"

func (c *backupClient) Get(backupID string) (*Backup, error) {
	url := fmt.Sprintf("/backups/%s", backupID)
	backup := Backup{}

	err := c.client.Get(url, &backup)
	if err != nil {
		return nil, err
	}

	return &backup, err
}

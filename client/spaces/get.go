package spaces

import (
	"fmt"
)

func (c *spacesClient) GetByID(spaceID string) (*Space, error) {
	var space Space

	err := c.client.Get("/spaces/" + spaceID, &space)
	if err != nil {
		return nil, fmt.Errorf("could not load space: %s", err)
	}

	return &space, nil
}

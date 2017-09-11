package teams

import (
	"fmt"
)

func (t *teamsClient) Get(idOrLabel string) (*Team, error) {
	var team Team

	url := fmt.Sprintf("/v1/teams/%s", idOrLabel)
	err := t.client.Get(url, &team)
	if err != nil {
		return nil, fmt.Errorf("could not get team: %s", err)
	}

	return &team, nil
}

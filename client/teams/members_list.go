package teams

import "fmt"

func (t *teamsClient) ListMembers(teamID string) ([]Membership, error) {
	var memberships []Membership

	url := fmt.Sprintf("/teams/%s/members", teamID)
	err := t.client.Get(url, &memberships)
	if err != nil {
		return nil, err
	}

	return memberships, nil
}
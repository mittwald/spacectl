package teams

import "fmt"

func (t *teamsClient) ListMembers(teamID string) ([]Membership, error) {
	var memberships []Membership
	var team Team

	err := t.client.Get("/teams/"+teamID, &team)
	if err != nil {
		return nil, err
	}

	link, err := team.Links.GetLinkByRel("members")
	if err != nil {
		return nil, err
	}

	err = link.Get(t.client, &memberships)
	if err != nil {
		return nil, fmt.Errorf("could not load team members: %s", err)
	}

	return memberships, nil
}

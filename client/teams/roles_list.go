package teams

import "fmt"

func (t *teamsClient) ListRoles(teamID string) ([]TeamRole, error) {
	var roles []TeamRole
	var team Team

	err := t.client.Get("/teams/"+teamID, &team)
	if err != nil {
		return nil, err
	}

	link, err := team.Links.GetLinkByRel("roles")
	if err != nil {
		return nil, err
	}

	err = link.Get(t.client, &roles)
	if err != nil {
		return nil, fmt.Errorf("could not load team roles: %s", err)
	}

	return roles, nil
}

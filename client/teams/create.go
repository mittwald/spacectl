package teams

import "fmt"

type createTeamRequest struct {
	Name    string `json:"name"`
	DNSName string `json:"dnsLabel"`
}

func (t *teamsClient) Create(name string, dnsLabel string) (Team, error) {
	var team Team

	req := createTeamRequest{name, dnsLabel}

	err := t.client.Post("/teams", &req, &team)
	if err != nil {
		return team, fmt.Errorf("could not create team: %s", err)
	}

	return team, nil
}

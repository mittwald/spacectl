package teams

import "fmt"

type createTeamRequest struct {
	Name string `json:"name"`
}

func (t *teamsClient) Create(name string) (*Team, error) {
	var team Team

	req := createTeamRequest{name}

	err := t.client.Post("/teams", &req, &team)
	if err != nil {
		return nil, fmt.Errorf("could not create team: %s", err)
	}

	return &team, nil
}
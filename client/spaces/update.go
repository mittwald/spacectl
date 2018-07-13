package spaces

import "fmt"

func (c *spacesClient) ListApplicationUpdatesByStage(spaceID, stageName string) ([]ApplicationUpdate, error) {
	var stage Stage
	var updates []ApplicationUpdate

	err := c.client.Get("/spaces/"+spaceID+"/stages/"+stageName, &stage)
	if err != nil {
		return nil, fmt.Errorf("could not load stage: %s", err)
	}

	link, err := stage.Links.GetLinkByRel("applicationUpdates")
	if err != nil {
		return nil, fmt.Errorf("could not access updates: %s", err)
	}

	err = link.Get(c.client, &updates)
	if err != nil {
		return nil, fmt.Errorf("could not load application updates: %s", err)
	}

	return updates, nil
}

func (c *spacesClient) ListApplicationUpdatesBySpace(spaceID string) ([]ApplicationUpdate, error) {
	var space Space
	var updates []ApplicationUpdate

	err := c.client.Get("/spaces/"+spaceID, &space)
	if err != nil {
		return nil, fmt.Errorf("could not load stage: %s", err)
	}

	link, err := space.Links.GetLinkByRel("applicationUpdates")
	if err != nil {
		return nil, fmt.Errorf("could not access updates: %s", err)
	}

	err = link.Get(c.client, &updates)
	if err != nil {
		return nil, fmt.Errorf("could not load application updates: %s", err)
	}

	return updates, nil
}

func (c *spacesClient) UpdateApplication(spaceID, stageName, targetStageName, version string) (*ApplicationUpdate, error) {
	var stage Stage
	var update ApplicationUpdate

	updateInput := ApplicationUpdateInput{
		VersionConstraint: version,
	}

	if targetStageName != "" {
		updateInput.TargetStage = &StageRef{
			Name: targetStageName,
		}
	}

	err := c.client.Get("/spaces/"+spaceID+"/stages/"+stageName, &stage)
	if err != nil {
		return nil, fmt.Errorf("could not load stage: %s", err)
	}

	link, err := stage.Actions.GetLinkByRel("applicationUpdate")
	if err != nil {
		return nil, fmt.Errorf("could not access updates: %s", err)
	}

	err = link.Post(c.client, updateInput, &update)
	if err != nil {
		return nil, fmt.Errorf("could not request application updates: %s", err)
	}

	return &update, nil
}

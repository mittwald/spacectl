package software

import "fmt"

func (t *softwareClient) List() ([]Software, error) {
	var software []Software

	err := t.client.Get("/softwares/"+t.group, &software)
	if err != nil {
		return nil, fmt.Errorf("could not load software: %s", err)
	}

	return software, nil
}

func (t *softwareClient) ListWithVersions() ([]Software, error) {
	softwareList, err := t.List()
	if err != nil {
		return nil, err
	}

	results := make([]Software, len(softwareList))

	for j := range softwareList {
		sw, err := t.Get(softwareList[j].Identifier)
		if err != nil {
			return nil, err
		}

		results[j] = *sw
	}

	return results, nil
}

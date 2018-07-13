package software

import (
	"fmt"
	"net/url"
)

func (t *softwareClient) Get(softwareID string) (*Software, error) {
	var software Software

	softwareURL := fmt.Sprintf("/softwares/%s/%s", url.PathEscape(t.group), url.PathEscape(softwareID))
	err := t.client.Get(softwareURL, &software)

	if err != nil {
		return nil, fmt.Errorf("could not load software: %s", err)
	}

	return &software, nil
}

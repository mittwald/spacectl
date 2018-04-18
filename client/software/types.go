package software

import "github.com/mittwald/spacectl/client/lowlevel"

type ContainerImage struct {
	Name string `json:"name"`
	HREF string `json:"href,omitempty"`
}

type Version struct {
	Number         string         `json:"number"`
	HREF           string         `json:"href,omitempty"`
	ContainerImage ContainerImage `json:"containerImage"`
}

type Software struct {
	Links      lowlevel.LinkList `json:"_links,omitempty"`
	Identifier string            `json:"identifier"`
	Name       string            `json:"name"`
	Versions   []Version         `json:"versions"`
}

func (s *Software) LatestVersion() *Version {
	if len(s.Versions) == 0 {
		return nil
	}

	return &s.Versions[len(s.Versions) - 1]
}
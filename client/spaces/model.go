package spaces

type SpaceName struct {
	DNSName           string `json:"dnsName"`
	HumanReadableName string `json:"humanReadableName"`
}

type SoftwareRef struct {
	ID   string `json:"id"`
	HREF string `json:"href"`
}

type VersionRef struct {
	Number string `json:"number"`
}

type TeamRef struct {
	Name     string `json:"name"`
	DNSLabel string `json:"dnsLabel"`
}

type Stage struct {
	Name              string      `json:"name"`
	Application       SoftwareRef `json:"application"`
	Version           VersionRef  `json:"version"`
	VersionConstraint string      `json:"versionConstraint"`
	DNSNames          []string    `json:"dnsNames"`
}

type Space struct {
	ID       string    `json:"id"`
	Name     SpaceName `json:"name"`
	Status   string    `json:"status"`
	DNSNames []string  `json:"dnsNames"`
	Stages   []Stage   `json:"stages"`
	Team     TeamRef   `json:"team"`
}

func (s Space) StagesCount() int {
	return len(s.Stages)
}

func (s Space) StagesNames() []string {
	names := make([]string, len(s.Stages))
	for i := range s.Stages {
		names[i] = s.Stages[i].Name
	}
	return names
}

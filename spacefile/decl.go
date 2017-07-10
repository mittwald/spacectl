package spacefile

import (
	"github.com/mittwald/spacectl/client/spaces"
)

func (s *SpaceDef) ToSpaceDeclaration() (*spaces.SpaceDeclaration, error) {
	stages := make([]spaces.StageDeclaration, len(s.Stages))

	for i := range s.Stages {
		st := &s.Stages[i]
		app := st.Applications[0]

		stageDecl := spaces.StageDeclaration{
			Name: st.Name,
			VersionConstraint: app.Version,
			Application: spaces.SoftwareRef{
				ID: app.Identifier,
			},
			UserData: app.UserData,
		}

		stages[i] = stageDecl
	}

	decl := spaces.SpaceDeclaration{
		Name: spaces.SpaceName{
			DNSName: s.DNSLabel,
			HumanReadableName: s.Name,
		},
		Stages: stages,
	}

	return &decl, nil
}
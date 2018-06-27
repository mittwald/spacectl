package spacefile

import (
	"github.com/mittwald/spacectl/client/spaces"
)

// ToSpaceDeclaration converts the SpaceDef object used in the Spacefile
// to a SpaceDeclaration used for the Spaces API calls
func (s *SpaceDef) ToSpaceDeclaration() (*spaces.SpaceDeclaration, error) {
	var err error
	stages := make([]spaces.StageDeclaration, len(s.Stages))

	for i := range s.Stages {
		st := &s.Stages[i]
		app := st.Application()

		appDecl := spaces.SoftwareRef{
			ID: app.Identifier,
		}

		cronjobDecls := make([]spaces.Cronjob, len(st.Cronjobs))
		for i := range st.Cronjobs {
			cronjobDecls[i], err = st.Cronjobs[i].ToDeclaration()
			if err != nil {
				return nil, err
			}
		}

		databaseDecls := make([]spaces.SoftwareVersionRef, len(st.Databases))
		for i := range st.Databases {
			databaseDecls[i] = spaces.SoftwareVersionRef{
				Software:          spaces.SoftwareRef{ID: st.Databases[i].Identifier},
				VersionConstraint: st.Databases[i].Version,
				UserData:          st.Databases[i].UserData,
			}
		}

		stageDecl := spaces.StageDeclaration{
			Name:              st.Name,
			VersionConstraint: app.Version,
			Application:       appDecl,
			Databases:         databaseDecls,
			UserData:          app.UserData,
			Cronjobs:          cronjobDecls,
		}

		stages[i] = stageDecl
	}

	decl := spaces.SpaceDeclaration{
		Name: spaces.SpaceName{
			DNSName:           s.DNSLabel,
			HumanReadableName: s.Name,
		},
		Stages: stages,
	}

	return &decl, nil
}

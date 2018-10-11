package spacefile

import (
	"github.com/mittwald/spacectl/client"
	"github.com/mittwald/spacectl/client/errors"
	"github.com/mittwald/spacectl/client/spaces"
)

// ToSpaceDeclaration converts the SpaceDef object used in the Spacefile
// to a SpaceDeclaration used for the Spaces API calls
func (s *SpaceDef) ToSpaceDeclaration() (*spaces.SpaceDeclaration, error) {
	stages := make([]spaces.StageDeclaration, len(s.Stages))

	for i := range s.Stages {
		st := &s.Stages[i]
		app := st.Application()

		appDecl := spaces.SoftwareRef{
			ID: app.Identifier,
		}

		cronjobDecls := make([]spaces.Cronjob, len(st.Cronjobs))
		for i := range st.Cronjobs {
			cronjobDecls[i] = st.Cronjobs[i].ToDeclaration()
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

// FromSpaceDeclaration converts a Space Declaration returned by the API
// to a SpaceDef object that can be returned as a spacefile
// it also makes additional api requests to gather information not contained in
// the SpaceDeclaration
func FromSpace(decl *spaces.Space, api client.SpacesClient) (*SpaceDef, error) {
	stages := make([]StageDef, len(decl.Stages))

	for i := range decl.Stages {
		st := decl.Stages[i]
		app := st.Application

		appDefs := make(SoftwareDefList, 1) // there can only be one
		appDefs[0] = SoftwareDef{
			Identifier: app.Software.ID,
			Version:    app.VersionConstraint,
			UserData:   st.UserData,
		}

		cronjobDefs := make(CronjobDefList, len(st.Cronjobs))
		for i := range st.Cronjobs {
			cronjobDefs[i] = CronjobFromDeclaration(&st.Cronjobs[i])
		}

		databaseDefs := make(SoftwareDefList, len(st.Databases))
		for i := range st.Databases {
			databaseDefs[i] = SoftwareDef{
				Identifier: st.Databases[i].Software.ID,
				Version:    st.Databases[i].VersionConstraint,
				UserData:   st.Databases[i].UserData,
			}
		}

		protection, err := api.Spaces().GetStageProtection(decl.ID, st.Name)
		if err != nil {
			return nil, errors.ErrNested{Inner: err, Msg: "could not get stage protection"}
		}

		stageDef := StageDef{
			Name:         st.Name,
			Applications: appDefs,
			Cronjobs:     cronjobDefs,
			Databases:    databaseDefs,
			Protection:   protection.ProtectionType,
		}

		stages[i] = stageDef
	}

	def := SpaceDef{
		TeamID:   decl.Team.DNSLabel,
		DNSLabel: decl.Name.DNSName,
		Name:     decl.Name.HumanReadableName,
		Stages:   stages,
	}

	return &def, nil
}

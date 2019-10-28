package spacefile

import (
	"fmt"
	"github.com/cloudfoundry/bytefmt"
	"github.com/mittwald/spacectl/client"
	"github.com/mittwald/spacectl/client/errors"
	"github.com/mittwald/spacectl/client/payment"
	"github.com/mittwald/spacectl/client/spaces"
	"math"
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
		for j := range st.Databases {
			databaseDecls[j] = spaces.SoftwareVersionRef{
				Software:          spaces.SoftwareRef{ID: st.Databases[j].Identifier},
				VersionConstraint: st.Databases[j].Version,
				UserData:          st.Databases[j].UserData,
			}

			if st.Databases[j].Storage.Size != "" {
				sizeBytes, err := bytefmt.ToBytes(st.Databases[j].Storage.Size)
				if err != nil {
					return nil, fmt.Errorf("cannot parse storage size %s (stage %d, database %d): %s", st.Databases[j].Storage.Size, i, j, err.Error())
				}

				databaseDecls[j].Storage = &spaces.SoftwareStorage{
					SizeGB: uint64(math.Ceil(float64(sizeBytes) / (1 << 30))),
				}
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

	prep := payment.SpaceResourcePreprovisioningInput{}
	opts := payment.SpaceOptionInput{}

	bytes, err := s.StorageBytes()
	if err != nil {
		return nil, err
	}
	if bytes > 0 {
		prep.Storage = &payment.SpaceResourcePreprovisioningInputItem{Quantity: bytes}
	}

	if stages := len(s.Stages) - s.CountOnDemandStages(); stages > 0 {
		prep.Stages = &payment.SpaceResourcePreprovisioningInputItem{Quantity: uint64(stages)}
	}

	if pods := s.Resource("scaling"); pods != nil {
		cnt, ok := pods.Quantity.(int)
		if !ok {
			return nil, fmt.Errorf("scaling quantity must be int, is %T", pods.Quantity)
		}
		prep.Scaling = &payment.SpaceResourcePreprovisioningInputItem{Quantity: uint64(cnt)}
	}

	if opt := s.Option("backupIntervalMinutes"); opt != nil {
		val, ok := opt.Value.(int)
		if !ok {
			return nil, fmt.Errorf("backup interval must be int, is %T", opt.Value)
		}
		opts.BackupIntervalMinutes = uint64(val)
	}

	decl := spaces.SpaceDeclaration{
		Name: spaces.SpaceName{
			DNSName:           s.DNSLabel,
			HumanReadableName: s.Name,
		},
		Stages: stages,
		PaymentLink: spaces.SpacePaymentLinkInput{
			Plan: payment.PlanReferenceInput{
				ID: s.Payment.PlanID,
			},
			PaymentProfile: payment.PaymentProfileReferenceInput{
				ID: s.Payment.PaymentProfileID,
			},
			Preprovisionings: &prep,
			Options:          &opts,
		},
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

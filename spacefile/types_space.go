package spacefile

import (
	"code.cloudfoundry.org/bytefmt"
	"errors"
	"fmt"

	"github.com/spf13/viper"

	"github.com/hashicorp/go-multierror"
)

type SpaceDef struct {
	DNSLabel  string        `hcl:",key"`
	Name      string        `hcl:"name"`
	TeamID    string        `hcl:"team"`
	Payment   PaymentDef    `hcl:"payment"`
	Resources []ResourceDef `hcl:"resource"`
	Stages    []StageDef    `hcl:"stage"`

	stagesByName    map[string]*StageDef
	resourcesByName map[string]*ResourceDef
}

func (d *SpaceDef) Validate(offline bool) error {
	var err *multierror.Error

	if len(d.DNSLabel) == 0 {
		err = multierror.Append(err, errors.New("empty Space name"))
	}

	if len(d.Name) == 0 {
		d.Name = d.DNSLabel
	}

	if len(d.Stages) == 0 {
		err = multierror.Append(err, fmt.Errorf("space \"%s\" should contain at least one stage", d.DNSLabel))
	}

	if len(d.TeamID) == 0 {
		if len(viper.GetString("teamID")) > 0 {
			d.TeamID = viper.GetString("teamID")
		} else {
			err = multierror.Append(err, errors.New("empty Team name"))
		}
	}

	for i := range d.Resources {
		switch d.Resources[i].Resource {
		case "storage":
			switch q := d.Resources[i].Quantity.(type) {
			case string:
				_, bErr := bytefmt.ToBytes(q)
				if bErr != nil {
					err = multierror.Append(err, fmt.Errorf("resource 'storage' must contain a valid byte value: %s", bErr))
				}
			case int:
			case int64:
			}
		default:
			_, ok := d.Resources[i].Quantity.(int)
			if !ok {
				err = multierror.Append(err, fmt.Errorf("quantity for resource '%s' must be int, is %T", d.Resources[i].Resource, d.Resources[i].Quantity))
			}
		}
	}

	for i := range d.Stages {
		err = multierror.Append(err, d.Stages[i].Validate(offline))
	}

	paymentErr := d.Payment.Validate(offline)
	if paymentErr != nil {
		err = multierror.Append(err, paymentErr)
	}

	return err.ErrorOrNil()
}

func (d *SpaceDef) resolveReferences() error {
	var err *multierror.Error

	d.stagesByName = make(map[string]*StageDef)
	d.resourcesByName = make(map[string]*ResourceDef)

	for i := range d.Stages {
		if _, ok := d.stagesByName[d.Stages[i].Name]; ok {
			err = multierror.Append(err, fmt.Errorf("duplicate stage declared: '%s'", d.Stages[i].Name))
		}

		d.stagesByName[d.Stages[i].Name] = &d.Stages[i]
	}

	for i := range d.Stages {
		if d.Stages[i].Inherit == "" {
			continue
		}

		parent, ok := d.stagesByName[d.Stages[i].Inherit]
		if !ok {
			err = multierror.Append(err, fmt.Errorf("stage '%s' in Space '%s' inherits non-existent stage '%s'", d.Stages[i].Name, d.DNSLabel, d.Stages[i].Inherit))
		} else {
			d.Stages[i].inheritStage = parent
		}
	}

	for i := range d.Stages {
		err = multierror.Append(err, d.Stages[i].resolveInheritance(0))
		err = multierror.Append(err, d.Stages[i].resolveUserData())
	}

	for i := range d.Resources {
		d.resourcesByName[d.Resources[i].Resource] = &d.Resources[i]
	}

	return err.ErrorOrNil()
}

func (d *SpaceDef) GetStageByName(stageName string) *StageDef {
	for _, stage := range d.Stages {
		if stage.Name == stageName {
			return &stage
		}
	}
	return nil
}

func (d *SpaceDef) CountOnDemandStages() int {
	c := 0

	for i := range d.Stages {
		if d.Stages[i].OnDemand {
			c ++
		}
	}

	return c
}

func (d *SpaceDef) Resource(r string) *ResourceDef {
	for i := range d.Resources {
		if d.Resources[i].Resource == r {
			return &d.Resources[i]
		}
	}

	return nil
}

func (d *SpaceDef) StorageBytes() (uint64, error) {
	r := d.Resource("storage")
	if r == nil {
		return 0, nil
	}

	switch q := r.Quantity.(type) {
	case string:
		b, err := bytefmt.ToBytes(q)
		if err != nil {
			return 0, err
		}

		return b, nil
	case int:
	case int64:
	case uint:
	case uint64:
		return uint64(q), nil
	}

	return 0, fmt.Errorf("unknown type for storage quantity: %T", r.Quantity)
}

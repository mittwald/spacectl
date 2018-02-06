package spacefile

import (
	"errors"
	"fmt"
	"github.com/hashicorp/go-multierror"
)

type SpaceDef struct {
	DNSLabel string     `hcl:",key"`
	Name     string     `hcl:"name"`
	TeamID   string     `hcl:"team"`
	Stages   []StageDef `hcl:"stage"`

	stagesByName map[string]*StageDef
}

func (d *SpaceDef) Validate() error {
	var err *multierror.Error

	if len(d.DNSLabel) == 0 {
		err = multierror.Append(err, errors.New("Empty Space name"))
	}

	if len(d.Stages) == 0 {
		err = multierror.Append(err, fmt.Errorf("Space \"%s\" should contain at least one stage", d.DNSLabel))
	}

	for i := range d.Stages {
		err = multierror.Append(err, d.Stages[i].Validate())
	}

	return err.ErrorOrNil()
}

func (d *SpaceDef) resolveReferences() error {
	var err *multierror.Error

	d.stagesByName = make(map[string]*StageDef)

	for i := range d.Stages {
		if _, ok := d.stagesByName[d.Stages[i].Name]; ok {
			err = multierror.Append(err, fmt.Errorf("Duplicate stage declared: '%s'", d.Stages[i].Name))
		}

		d.stagesByName[d.Stages[i].Name] = &d.Stages[i]
	}

	for i := range d.Stages {
		if d.Stages[i].Inherit == "" {
			continue
		}

		parent, ok := d.stagesByName[d.Stages[i].Inherit]
		if !ok {
			err = multierror.Append(err, fmt.Errorf("Stage '%s' in Space '%s' inherits non-existent stage '%s'", d.Stages[i].Name, d.DNSLabel, d.Stages[i].Inherit))
		} else {
			d.Stages[i].inheritStage = parent
		}
	}

	for i := range d.Stages {
		err = multierror.Append(err, d.Stages[i].resolveInheritance(0))
		err = multierror.Append(err, d.Stages[i].resolveUserData())
	}

	return err.ErrorOrNil()
}
package spacefile

import (
	"errors"
	"fmt"
	"github.com/hashicorp/go-multierror"
	"github.com/imdario/mergo"
)

const DefaultFilename = ".spacefile.hcl"

type SpaceDef struct {
	DNSLabel string     `hcl:",key"`
	Name     string     `hcl:"name"`
	TeamID   string     `hcl:"team"`
	Stages   []StageDef `hcl:"stage"`

	stagesByName map[string]*StageDef
}

type SoftwareDef struct {
	Identifier string      `hcl:",key"`
	Version    string      `hcl:"version"`
	UserData   interface{} `hcl:"userData"`
}

type StageDef struct {
	Name         string        `hcl:",key"`
	Inherit      string        `hcl:"inherit"`
	Applications []SoftwareDef `hcl:"application"`
	Databases    []SoftwareDef `hcl:"database"`

	inheritStage *StageDef
}

type Spacefile struct {
	Version string     `hcl:"version"`
	Spaces  []SpaceDef `hcl:"space"`
}

func (f *Spacefile) Validate() error {
	var err *multierror.Error

	if f.Version != "1" {
		err = multierror.Append(err, fmt.Errorf("Unsupported version: %s", f.Version))
	}

	if len(f.Spaces) == 0 {
		err = multierror.Append(err, errors.New("Spacefile does not contain a space definition"))
	}

	if len(f.Spaces) > 1 {
		err = multierror.Append(err, errors.New("Spacefile should not contain more than one space definition"))
	}

	for i := range f.Spaces {
		err = multierror.Append(err, f.Spaces[i].Validate())
	}

	return err.ErrorOrNil()
}

func (f *Spacefile) resolveReferences() error {
	var err *multierror.Error

	for i := range f.Spaces {
		err = multierror.Append(err, f.Spaces[i].resolveReferences())
	}

	return err.ErrorOrNil()
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

func (d *StageDef) Validate() error {
	var err *multierror.Error

	if len(d.Applications) > 1 {
		err = multierror.Append(err, fmt.Errorf("Stage '%s' shoud not contain more than one application", d.Name))
	}

	return err.ErrorOrNil()
}

func (d *StageDef) resolveUserData() error {
	var mErr *multierror.Error
	var err error

	for i := range d.Applications {
		d.Applications[i].UserData, err = unfuckHCL(d.Applications[i].UserData, "")
		mErr = multierror.Append(mErr, err)
	}

	return mErr.ErrorOrNil()
}

func (d *StageDef) resolveInheritance(level int) error {
	if level > 4 {
		return fmt.Errorf("Could not resolve stage dependencies after %d levels. Please check that there is no cyclic inheritance", level)
	}

	if d.inheritStage == nil {
		return nil
	}

	err := d.inheritStage.resolveInheritance(level + 1)
	if err != nil {
		return err
	}

	originalName := d.Name

	err = mergo.MergeWithOverwrite(d, d.inheritStage)
	if err != nil {
		return err
	}

	d.Name = originalName
	d.inheritStage = nil

	return nil
}

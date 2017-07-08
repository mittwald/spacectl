package spacefile

import (
	"fmt"
	"github.com/hashicorp/go-multierror"
	"bytes"
	"encoding/gob"
	"github.com/imdario/mergo"
)

type SpaceDef struct {
	DNSLabel string `hcl:",key"`
	Name string `hcl:"name"`
	TeamID string `hcl:"team"`
	Stages []StageDef `hcl:"stage"`

	stagesByName map[string]*StageDef
}

type SoftwareDef struct {
	Identifier string `hcl:",key"`
	Version string `hcl:"version"`
	UserData interface{}
}

type StageDef struct {
	Name string `hcl:",key"`
	Inherit string `hcl:"inherit"`
	Applications []SoftwareDef `hcl:"application"`
	Databases []SoftwareDef `hcl:"database"`

	inheritStage *StageDef
}

type Spacefile struct {
	Spaces []SpaceDef `hcl:"space"`
}

func (f *Spacefile) resolveReferences() error {
	var err *multierror.Error

	for i := range f.Spaces {
		err = multierror.Append(err, f.Spaces[i].resolveStageInheritance())
	}

	return err.ErrorOrNil()
}

func (d *SpaceDef) resolveStageInheritance() error {
	var err *multierror.Error

	d.stagesByName = make(map[string]*StageDef)

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
			err = multierror.Append(err, fmt.Errorf("stage '%s' inherits non-existent stage '%s'", d.Stages[i].Name, d.Stages[i].Inherit))
		} else {
			d.Stages[i].inheritStage = parent
		}
	}

	for i := range d.Stages {
		err = multierror.Append(err, d.Stages[i].resolveInheritance(0))
	}

	return err.ErrorOrNil()
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

	err = mergo.Merge(d, d.inheritStage)
	if err != nil {
		return err
	}

	d.inheritStage = nil

	/*
	myAppsByName := make(map[string]SoftwareDef)
	for k := range d.Applications {
		myAppsByName[d.Applications[k].Identifier] = d.Applications[k]
	}

	var err error

	d.Applications = make([]SoftwareDef, 0)

	for k, a := range d.inheritStage.Applications {
		dataCopy, copyErr := deepCopy(a.UserData)
		if copyErr != nil {
			err = multierror.Append(err, copyErr)
		} else {
			newApp := SoftwareDef{Identifier: a.Identifier, Version: a.Version, UserData: dataCopy}

			if overwrite, ok := myAppsByName[a.Identifier]; ok {
				if overwrite.Version != "" {
					newApp.Version = overwrite.Version
				}
			}

			d.Applications = append(d.Applications, newApp)
		}
	}*/

	return nil
}

func deepCopy(obj interface{}) (interface{}, error) {
	var mod bytes.Buffer
	enc := gob.NewEncoder(&mod)
	dec := gob.NewDecoder(&mod)

	err := enc.Encode(obj)
	if err != nil {
		return nil, err
	}

	var cpy interface{}
	err = dec.Decode(&cpy)
	if err != nil {
		return nil, err
	}

	return cpy, nil
}
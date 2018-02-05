package spacefile

import (
	"errors"
	"fmt"

	"github.com/hashicorp/go-multierror"
)

const DefaultFilename = ".spacefile.hcl"

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

package spacefile

import (
	"io/ioutil"
	"fmt"
	"github.com/hashicorp/hcl"
	"github.com/hashicorp/go-multierror"
	"os"
)

func ParseSpacefile(filename string) (*Spacefile, error) {
	contents, err := ioutil.ReadFile(filename)

	if os.IsNotExist(err) {
		return nil, fmt.Errorf("The Spacefile at %s does not exist, yet. Create one using the \"spacectl space init\" command.", filename)
	}

	if err != nil {
		return nil, fmt.Errorf("Could not read Spacefile at %s: %s", filename, err)
	}

	obj := Spacefile{}

	err = hcl.Decode(&obj, string(contents))
	if err != nil {
		return nil, SyntaxError{filename, err}
	}

	var mErr *multierror.Error

	mErr = multierror.Append(mErr, obj.resolveReferences())
	mErr = multierror.Append(mErr, obj.Validate())

	if mErr.ErrorOrNil() != nil {
		return nil, mErr
	}

	return &obj, nil
}
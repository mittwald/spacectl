package spacefile

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/hcl"
)

func ParseSpacefile(filename string, offline bool) (*Spacefile, error) {
	contents, err := ioutil.ReadFile(filename)

	if os.IsNotExist(err) {
		return nil, ErrSpacefileNotFound{filename}
	}

	if err != nil {
		return nil, fmt.Errorf("could not read Spacefile at %s: %s", filename, err)
	}

	obj := Spacefile{}

	err = hcl.Decode(&obj, string(contents))
	if err != nil {
		return nil, SyntaxError{filename, err}
	}

	var mErr *multierror.Error

	mErr = multierror.Append(mErr, obj.resolveReferences())
	mErr = multierror.Append(mErr, obj.Validate(offline))

	if mErr.ErrorOrNil() != nil {
		return nil, mErr
	}

	return &obj, nil
}

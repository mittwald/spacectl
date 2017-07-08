package spacefile

import (
	"io/ioutil"
	"fmt"
	"github.com/hashicorp/hcl"
)

func ParseSpacefile(filename string) (*Spacefile, error) {
	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("could not read Spacefile at %s: %s", filename, err)
	}

	obj := Spacefile{}

	err = hcl.Decode(&obj, string(contents))
	if err != nil {
		return nil, SyntaxError{filename, err}
	}

	err = obj.resolveReferences()
	if err != nil {
		return nil, err
	}

	return &obj, nil
}
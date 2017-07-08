package spacefile

import (
	"github.com/hashicorp/go-multierror"
	"fmt"
)

func unfuckHCL(in interface{}, path string) (interface{}, error) {
	var mErr *multierror.Error
	var err error

	fmt.Printf("%v (type=%T)\n", in, in)

	switch mapped := in.(type) {
	case map[string]interface{}:
		for key := range mapped {
			mapped[key], err = unfuckHCL(mapped[key], path + "." + key)
			mErr = multierror.Append(mErr, err)
		}
	case []interface{}:
		if len(mapped) > 1 {
			mErr = multierror.Append(mErr, fmt.Errorf("more than 1 element for %s", path))
		}

		in = mapped[0]
	case []map[string]interface{}:
		if len(mapped) > 1 {
			mErr = multierror.Append(mErr, fmt.Errorf("more than 1 element for %s", path))
		}

		in = mapped[0]
	}

	return in, mErr.ErrorOrNil()
}

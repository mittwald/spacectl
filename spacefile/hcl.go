package spacefile

import (
	"fmt"

	"github.com/hashicorp/go-multierror"
)

func unfuckHCL(in interface{}, path string) (interface{}, error) {
	var mErr *multierror.Error
	var err error

	switch mapped := in.(type) {
	case map[string]interface{}:
		for key := range mapped {
			mapped[key], err = unfuckHCL(mapped[key], path+"."+key)
			mErr = multierror.Append(mErr, err)
		}
	case []interface{}:
		if len(mapped) > 1 {
			mErr = multierror.Append(mErr, fmt.Errorf("more than 1 element for %s", path))
		}

		in = mapped[0]
	case []map[string]interface{}:
		for key := range mapped {
			for subKey := range mapped[key] {
				mapped[0][subKey], err = unfuckHCL(mapped[key][subKey], path+".0."+subKey)
			}
		}

		in, err = unfuckHCL(mapped[0], path+".0")
		mErr = multierror.Append(mErr, err)
	}

	return in, mErr.ErrorOrNil()
}

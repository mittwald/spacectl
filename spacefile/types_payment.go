package spacefile

import (
	"errors"
	"github.com/hashicorp/go-multierror"
)

type PaymentDef struct {
	PaymentProfileID string `hcl:"paymentProfile"`
	PlanID           string `hcl:"plan"`
}

func (d *PaymentDef) Validate(offline bool) error {
	var err *multierror.Error

	if d.PaymentProfileID == "" {
		err = multierror.Append(err, errors.New("'paymentProfile' must not be empty"))
	}

	if d.PlanID == "" {
		err = multierror.Append(err, errors.New("'planID' must not be empty"))
	}

	return err.ErrorOrNil()
}
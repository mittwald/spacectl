package spaces

import (
	"code.cloudfoundry.org/bytefmt"
	"fmt"
	"github.com/mittwald/spacectl/client/lowlevel"
	"github.com/mittwald/spacectl/client/payment"
	"strings"
)

type ConnectOption func(in *SpacePaymentLinkInput) error

func (c *spacesClient) ConnectWithPaymentProfile(spaceID string, paymentProfileID string, planID string, opts ...ConnectOption) (*SpacePaymentLink, error) {
	var space Space
	var paymentLink SpacePaymentLink

	err := c.client.Get("/spaces/"+spaceID, &space)
	if err != nil {
		return nil, fmt.Errorf("could not load stage: %s", err)
	}

	link, err := space.Links.GetLinkByRel("paymentlink")
	if err != nil {
		return nil, fmt.Errorf("could not access payment connection: %s", err.Error())
	}

	preprov := payment.SpaceResourcePreprovisioningInput{}

	input := SpacePaymentLinkInput{
		Plan:             payment.PlanReferenceInput{ID: planID},
		PaymentProfile:   payment.PaymentProfileReferenceInput{ID: paymentProfileID},
		Preprovisionings: &preprov,
	}

	for i := range opts {
		if err := opts[i](&input); err != nil {
			return nil, err
		}
	}

	err = link.Put(c.client, &input, &paymentLink)
	if err != nil {
		return nil, fmt.Errorf("could not connect payment profile: %s", err.Error())
	}

	return &paymentLink, nil
}

func WithStorage(storageBytes uint64) ConnectOption {
	return func(i *SpacePaymentLinkInput) error {
		if i.Preprovisionings == nil {
			i.Preprovisionings = &payment.SpaceResourcePreprovisioningInput{}
		}

		i.Preprovisionings.Storage = &payment.SpaceResourcePreprovisioningInputItem{
			Quantity: storageBytes,
		}
		return nil
	}
}

func WithStorageStr(storage string) ConnectOption {
	return func(i *SpacePaymentLinkInput) error {
		if i.Preprovisionings == nil {
			i.Preprovisionings = &payment.SpaceResourcePreprovisioningInput{}
		}

		if strings.TrimRight(storage, "TGMKIB") == "0" {
			i.Preprovisionings.Storage = &payment.SpaceResourcePreprovisioningInputItem{
				Quantity: 0,
			}
			return nil
		}

		b, err := bytefmt.ToBytes(storage)
		if err != nil {
			return err
		}

		i.Preprovisionings.Storage = &payment.SpaceResourcePreprovisioningInputItem{
			Quantity: b,
		}
		return nil
	}
}

func WithStages(stages uint64) ConnectOption {
	return func(i *SpacePaymentLinkInput) error {
		i.Preprovisionings.Stages = &payment.SpaceResourcePreprovisioningInputItem{
			Quantity: stages,
		}
		return nil
	}
}

func WithPods(pods uint64) ConnectOption {
	return func(i *SpacePaymentLinkInput) error {
		i.Preprovisionings.Scaling = &payment.SpaceResourcePreprovisioningInputItem{
			Quantity: pods,
		}
		return nil
	}
}

func WithoutTestingPeriod() ConnectOption {
	return func(i *SpacePaymentLinkInput) error {
		i.SkipTestingPeriod = false
		return nil
	}
}

func WithBackupIntervalMinutes(interval uint64) ConnectOption {
	return func(i *SpacePaymentLinkInput) error {
		i.Options.BackupIntervalMinutes = interval
		return nil
	}
}

func (c *spacesClient) GetPaymentLink(spaceID string) (*SpacePaymentLink, error) {
	var space Space
	var paymentLink SpacePaymentLink

	err := c.client.Get("/spaces/"+spaceID, &space)
	if err != nil {
		return nil, fmt.Errorf("could not load stage: %s", err)
	}

	link, err := space.Links.GetLinkByRel("paymentlink")
	if err != nil {
		return nil, fmt.Errorf("could not access payment connection: %s", err.Error())
	}

	err = link.Get(c.client, &paymentLink)
	if err != nil {
		statusErr, ok := err.(lowlevel.ErrUnexpectedStatusCode)
		if ok && statusErr.StatusCode == 404 {
			return nil, nil
		}

		return nil, fmt.Errorf("could not load payment link: %s", err.Error())
	}

	return &paymentLink, nil
}

package spaces

import (
	"fmt"
	"github.com/mittwald/spacectl/client/lowlevel"
	"github.com/mittwald/spacectl/client/payment"
)

func (c *spacesClient) ConnectWithPaymentProfile(spaceID string, paymentProfileID string, planID string) (*SpacePaymentLink, error) {
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

	input := SpacePaymentLinkInput{
		Plan:           payment.PlanReferenceInput{ID: planID},
		PaymentProfile: payment.PaymentProfileReferenceInput{ID: paymentProfileID},
	}

	err = link.Put(c.client, &input, &paymentLink)
	if err != nil {
		return nil, fmt.Errorf("could not connect payment profile: %s", err.Error())
	}

	return &paymentLink, nil
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

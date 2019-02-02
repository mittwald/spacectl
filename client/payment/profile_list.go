package payment

import "fmt"

func (c *paymentClient) ListPaymentProfiles() ([]PaymentProfile, error) {
	var profiles []PaymentProfile

	err := c.client.Get("/paymentprofiles", &profiles)
	if err != nil {
		return nil, fmt.Errorf("could not load profiles: %s", err)
	}

	return profiles, nil
}

package payment

import "fmt"

func (c *paymentClient) ConnectProfile(customerNumber string, password string) (*PaymentProfile, error) {
	var profile PaymentProfile

	req := PaymentProfileConnectRequest{
		CustomerNumber: customerNumber,
		Password:       password,
	}

	err := c.client.Post("/paymentprofiles/actions/connect", req, &profile)
	if err != nil {
		return nil, fmt.Errorf("could not load profiles: %s", err)
	}

	return &profile, nil
}

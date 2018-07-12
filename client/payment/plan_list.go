package payment

import "fmt"

func (c *paymentClient) ListPlans() ([]Plan, error) {
	var plans []Plan

	err := c.client.Get("/plans", &plans)
	if err != nil {
		return nil, fmt.Errorf("could not load plans: %s", err)
	}

	return plans, nil
}

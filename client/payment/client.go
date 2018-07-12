package payment

import "github.com/mittwald/spacectl/client/lowlevel"

type PaymentClient interface {
	ListPlans() ([]Plan, error)
}

func NewPaymentClient(client *lowlevel.SpacesLowlevelClient) (PaymentClient) {
	return &paymentClient{client}
}

type paymentClient struct {
	client *lowlevel.SpacesLowlevelClient
}
package payment

import "github.com/mittwald/spacectl/client/lowlevel"

type Client interface {
	ListPlans() ([]Plan, error)
	ListPaymentProfiles() ([]PaymentProfile, error)
	ConnectProfile(customerNumber string, password string) (*PaymentProfile, error)
}

func NewClient(client *lowlevel.SpacesLowlevelClient) Client {
	return &paymentClient{client}
}

type paymentClient struct {
	client *lowlevel.SpacesLowlevelClient
}

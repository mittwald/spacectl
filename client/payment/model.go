package payment

import "time"

type BasePrice struct {
	Value    float32 `json:"value"`
	Currency string  `json:"currency"`
	Unit     string  `json:"unit,omitempty"`
}

type Plan struct {
	ID        string     `json:"id"`
	Name      string     `json:"name"`
	BasePrice *BasePrice `json:"basePrice"`
}

type PlanReferenceInput struct {
	ID string `json:"id"`
}

type PaymentProfileContactAddress struct {
	Street      string `json:"street"`
	HouseNumber string `json:"houseNumber"`
	City        string `json:"city"`
	ZIP         string `json:"zip"`
	Country     string `json:"country"`
}

type PaymentProfileContact struct {
	FirstName    string                       `json:"firstName"`
	LastName     string                       `json:"lastName"`
	Title        string                       `json:"title,omitempty"`
	Salutation   string                       `json:"salutation"`
	Company      string                       `json:"company,omitempty"`
	Address      PaymentProfileContactAddress `json:"address"`
	EmailAddress string                       `json:"emailAddress"`
}

type PaymentProfileContactWithPhone struct {
	FirstName    string                       `json:"firstName"`
	LastName     string                       `json:"lastName"`
	Title        string                       `json:"title,omitempty"`
	Salutation   string                       `json:"salutation"`
	Company      string                       `json:"company,omitempty"`
	Address      PaymentProfileContactAddress `json:"address"`
	EmailAddress string                       `json:"emailAddress"`
	Phone        string                       `json:"phone"`
}

type PaymentProfileInvoiceSettings struct {
	InvoiceRecipient          *PaymentProfileContact `json:"invoiceRecipient,omitempty"`
	AdditionalEmailRecipients []string               `json:"additionalEmailRecipients,omitempty"`
	PrintedInvoices           bool                   `json:"printedInvoices"`
	VatID                     string                 `json:"vatID,omitempty"`
}

type PaymentOption struct {
	Type          string `json:"type"`
	IBAN          string `json:"iban,omitempty"`
	BIC           string `json:"bic,omitempty"`
	AccountHolder string `json:"accountHolder,omitempty"`
}

type PaymentProfile struct {
	ID              string                         `json:"id"`
	CreatedAt       time.Time                      `json:"createdAt"`
	ModifiedAt      time.Time                      `json:"modifiedAt"`
	ContractPartner PaymentProfileContactWithPhone `json:"contractPartner"`
	InvoiceSettings PaymentProfileInvoiceSettings  `json:"invoiceSettings"`
	Payment         PaymentOption                  `json:"payment"`
}

type PaymentProfileReferenceInput struct {
	ID string `json:"id"`
}

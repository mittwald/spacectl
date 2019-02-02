package view

import (
	"fmt"
	"github.com/gosuri/uitable"
	"github.com/mittwald/spacectl/client/payment"
	"github.com/mittwald/spacectl/cmd/helper"
	"io"
	"time"
)

type PaymentProfileView struct {
	payment.PaymentProfile
}

func (v *PaymentProfileView) ContractPartnerSummary() string {
	if v.ContractPartner.Company == "" {
		return fmt.Sprintf("%s %s %s", v.ContractPartner.Salutation, v.ContractPartner.FirstName, v.ContractPartner.LastName)
	}

	return fmt.Sprintf(
		"%s %s %s, %s",
		v.ContractPartner.Salutation,
		v.ContractPartner.FirstName,
		v.ContractPartner.LastName,
		v.ContractPartner.Company,
	)
}

type PaymentProfileListView interface {
	List(profiles []payment.PaymentProfile, out io.Writer)
}

type TabularPaymentProfileListView struct{}

func (t TabularPaymentProfileListView) List(profiles []payment.PaymentProfile, out io.Writer) {
	if len(profiles) == 0 {
		fmt.Fprintln(out, "No payment profiles found.")
		return
	}

	table := uitable.New()
	table.MaxColWidth = 50
	table.AddRow("ID", "CONTRACT PARTNER", "BILLING TYPE", "CREATED")

	for _, profile := range profiles {
		v := PaymentProfileView{profile}
		since := helper.HumanReadableDateDiff(time.Now(), profile.CreatedAt)

		table.AddRow(
			profile.ID,
			v.ContractPartnerSummary(),
			profile.Payment.Type,
			since + " ago",
		)
	}

	fmt.Fprintln(out, table)
}

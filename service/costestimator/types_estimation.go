package costestimator

import "golang.org/x/text/currency"

type EstimationLineItem struct {
	Subject     string
	Quantity    Quantity
	MonthlyCost currency.Amount
}

type Estimation struct {
	MonthlyTotalCost currency.Amount

	LineItems []EstimationLineItem

	OnDemandCharges struct {
		Storage BasePrice
		Stages  BasePrice
	}
}

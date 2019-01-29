package costestimator

import (
	"golang.org/x/text/currency"
)

type Quantity struct {
	Value int
	Unit  Unit
}

type BasePrice struct {
	Amount currency.Amount
	Unit   Unit
}

func (q *Quantity) String() string {
	return q.Unit.Format(q.Value)
}
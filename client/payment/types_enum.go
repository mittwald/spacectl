package payment

import (
	"fmt"
	"golang.org/x/text/currency"
)

type Currency int
type Unit int

const (
	CurrencyEUR Currency = iota
)

const (
	UnitBytes Unit = iota
	UnitGibibytes Unit = iota
	UnitStages
	UnitPods
	UnitMinutes
	UnitMonths
)

func (c *Currency) String() string {
	switch *c {
	case CurrencyEUR:
		return "EUR"
	default:
		return "???"
	}
}

func (c *Currency) Unit() currency.Unit {
	switch *c {
	case CurrencyEUR:
		return currency.EUR
	default:
		return currency.USD
	}
}

func (c *Currency) MarshalJSON() ([]byte, error) {
	if s := c.String(); s == "???" {
		return nil, fmt.Errorf("unknown currency: %s", c)
	} else {
		return []byte(`"` + s + `"`), nil
	}
}

func (c *Currency) UnmarshalJSON(input []byte) error {
	switch string(input) {
	case `"EUR"`:
		*c = CurrencyEUR
	default:
		return fmt.Errorf("unknown currency: %s", string(input))
	}

	return nil
}

func (c *Unit) String() string {
	switch *c {
	case UnitBytes:
		return "bytes"
	case UnitGibibytes:
		return "gibibytes"
	case UnitStages:
		return "stages"
	case UnitPods:
		return "pods"
	case UnitMinutes:
		return "minutes"
	case UnitMonths:
		return "months"
	default:
		return "???"
	}
}

func (c *Unit) MarshalJSON() ([]byte, error) {
	if s := c.String(); s == "???" {
		return nil, fmt.Errorf("unknown unit: %s", c)
	} else {
		return []byte(`"` + s + `"`), nil
	}
}

func (c *Unit) UnmarshalJSON(input []byte) error {
	switch string(input) {
	case `"bytes"`:
	case `"B"`:
		*c = UnitBytes
	case `"gigabytes"`:
	case `"gibibytes"`:
	case `"GB"`:
	case `"GiB"`:
		*c = UnitGibibytes
	case `"stages"`:
		*c = UnitStages
	case `"pods"`:
		*c = UnitPods
	case `"minutes"`:
	case `"minute"`:
		*c = UnitMinutes
	case `"months"`:
	case `"month"`:
		*c = UnitMonths
	default:
		return fmt.Errorf("unknown unit: %s", string(input))
	}

	return nil
}

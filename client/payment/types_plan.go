package payment

import (
	"fmt"
	"golang.org/x/text/currency"
)

type BasePrice struct {
	Value    float64  `json:"value"`
	Currency Currency `json:"currency"`
	Unit     Unit     `json:"unit,omitempty"`
	Interval string   `json:"interval,omitempty"`
}

func (p BasePrice) CurrencyAmount() currency.Amount {
	return p.Currency.Unit().Amount(p.Value)
}

func (p BasePrice) MustConvertUnit(to Unit) BasePrice {
	n, err := p.ConvertUnit(to)
	if err != nil {
		panic(err)
	}

	return n
}

func (p BasePrice) ConvertUnit(to Unit) (BasePrice, error) {
	switch to {
	case UnitGibibytes:
		switch p.Unit {
		case UnitBytes:
			return BasePrice{
				Value:    p.Value * (1 << 30),
				Currency: p.Currency,
				Unit:     UnitGibibytes,
				Interval: p.Interval,
			}, nil
		case UnitGibibytes:
			return p, nil
		}
	case UnitBytes:
		switch p.Unit {
		case UnitBytes:
			return p, nil
		case UnitGibibytes:
			return BasePrice{
				Value:    p.Value / (1 << 30),
				Currency: p.Currency,
				Unit:     UnitBytes,
				Interval: p.Interval,
			}, nil
		}
	}

	return BasePrice{}, fmt.Errorf("unsupported type conversion: %s to %s", p.Unit.String(), to.String())
}

type Plan struct {
	ID        string       `json:"id"`
	Name      string       `json:"name"`
	BasePrice *BasePrice   `json:"basePrice"`
	Features  PlanFeatures `json:"features"`
}

type PlanFeatures struct {
	Storage PlanFeature        `json:"storage"`
	Stages  PlanFeature        `json:"stages"`
	Backups PlanFeatureBackup  `json:"backups"`
	Scaling PlanFeatureScaling `json:"scaling"`
}

type PlanFeatureBackupOption struct {
	MinimumInterval Quantity  `json:"minimumInterval"`
	BasePrice       BasePrice `json:"basePrice"`
}

type PlanFeatureBackup struct {
	MinimumInterval Quantity                  `json:"minimumInterval"`
	Options         []PlanFeatureBackupOption `json:"options,omitempty"`
}

type Exceedances struct {
	PreProvision Exceedance `json:"preprovision"`
	OnDemand     Exceedance `json:"ondemand"`
}

type ExceedancesWithoutOnDemand struct {
	PreProvision Exceedance `json:"preprovision"`
}

type PlanFeatureScaling struct {
	Minimum    Quantity                   `json:"minimum"`
	Maximum    Quantity                   `json:"maximum"`
	Exceedance ExceedancesWithoutOnDemand `json:"exceedance"`
}

type PlanFeature struct {
	Included   Quantity    `json:"included"`
	Exceedance Exceedances `json:"exceedance"`
}

type Exceedance struct {
	Quantity  Quantity  `json:"quantity"`
	BasePrice BasePrice `json:"basePrice"`
}

type Quantity struct {
	Value uint64 `json:"value"`
	Unit  Unit   `json:"unit"`
}

type PlanReferenceInput struct {
	ID string `json:"id"`
}

package costestimator

import (
	"fmt"
	"github.com/mittwald/spacectl/client/payment"
)

type estimator struct {
	plans []payment.Plan
}

func New(plans []payment.Plan) Estimator {
	return &estimator{
		plans,
	}
}

func (e *estimator) Estimate(params Params) (*Estimation, error) {
	res := Estimation{}

	plan, err := e.getPlan(params.PlanID)
	if err != nil {
		return nil, err
	}

	totalCurrency := plan.BasePrice.Currency.Unit()
	totalCost := plan.BasePrice.Value

	res.LineItems = make([]EstimationLineItem, 1)

	mainLineItem := EstimationLineItem{
		Quantity: Quantity{Value: 1, Unit: NewUnit("month")},
		Subject: "Basic cloud infrastructure",
		MonthlyCost: plan.BasePrice.CurrencyAmount(),
	}

	res.LineItems[0] = mainLineItem

	if params.Storage > 0 {
		storageBasePricePerByte, err := plan.Features.Storage.Exceedance.PreProvision.BasePrice.ConvertUnit(payment.UnitBytes)
		if err != nil {
			return nil, err
		}

		storageCost := storageBasePricePerByte.Value * float64(params.Storage)
		totalCost += storageCost

		storageLineItem := EstimationLineItem{
			Quantity: Quantity{Value: int(params.Storage >> 30), Unit: NewUnit("GB")},
			Subject: "Storage costs",
			MonthlyCost: totalCurrency.Amount(storageCost),
		}

		res.LineItems = append(res.LineItems, storageLineItem)
	}

	if params.Stages > 0 {
		preStages := params.Stages - params.StagesOnDemand
		preStageCost := plan.Features.Stages.Exceedance.PreProvision.BasePrice.Value * float64(preStages)
		totalCost += preStageCost

		res.LineItems = append(res.LineItems, EstimationLineItem{
			Quantity: Quantity{Value: preStages, Unit: NewUnit("stages")},
			Subject: "Preprovisioned stages (billed monthly in-advance)",
			MonthlyCost: totalCurrency.Amount(preStageCost),
		})

		if params.StagesOnDemand > 0 {
			onDemandStageCost := plan.Features.Stages.Exceedance.OnDemand.BasePrice.Value * float64(params.StagesOnDemand)
			totalCost += onDemandStageCost

			res.LineItems = append(res.LineItems, EstimationLineItem{
				Quantity: Quantity{Value: preStages, Unit: NewUnit("stages")},
				Subject: "On-demand stages (billed retroactively per-usage)",
				MonthlyCost: totalCurrency.Amount(onDemandStageCost),
			})
		}
	}

	res.MonthlyTotalCost = totalCurrency.Amount(totalCost)
	res.OnDemandCharges.Stages.Amount = plan.Features.Stages.Exceedance.OnDemand.BasePrice.CurrencyAmount()
	res.OnDemandCharges.Stages.Unit = NewUnit(plan.Features.Stages.Exceedance.OnDemand.BasePrice.Unit.String())
	res.OnDemandCharges.Storage.Amount = plan.Features.Storage.Exceedance.OnDemand.BasePrice.MustConvertUnit(payment.UnitGibibytes).CurrencyAmount()
	res.OnDemandCharges.Storage.Unit = NewUnit("GB")

	return &res, nil
}

func (e *estimator) getPlan(planID string) (*payment.Plan, error) {
	for i := range e.plans {
		if e.plans[i].ID == planID {
			return &e.plans[i], nil
		}
	}

	return nil, fmt.Errorf("plan does not exist: %s", planID)
}

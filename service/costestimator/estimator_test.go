package costestimator

import (
	myassert "github.com/mittwald/spacectl/assert"
	"github.com/mittwald/spacectl/client/payment"
	"github.com/stretchr/testify/assert"
	"testing"
)

var testPlanOne = payment.Plan{
	ID: "spaces.test/v1",
	Name: "Test-Space",
	BasePrice: &payment.BasePrice{
		Unit: payment.UnitMonths,
		Value: 99.00,
		Interval: "month",
		Currency: payment.CurrencyEUR,
	},
	Features: payment.PlanFeatures{
		Storage: payment.PlanFeature{
			Included: payment.Quantity{Value: 0, Unit: payment.UnitBytes},
			Exceedance: payment.Exceedances{
				PreProvision: payment.Exceedance{
					Quantity: payment.Quantity{Value: 1, Unit: payment.UnitGibibytes},
					BasePrice: payment.BasePrice{
						Unit: payment.UnitGibibytes,
						Value: .3,
						Currency: payment.CurrencyEUR,
						Interval: "month",
					},
				},
				OnDemand: payment.Exceedance{
					Quantity: payment.Quantity{Value: 1, Unit: payment.UnitGibibytes},
					BasePrice: payment.BasePrice{
						Unit: payment.UnitGibibytes,
						Value: .70,
						Currency: payment.CurrencyEUR,
						Interval: "month",
					},
				},
			},
		},
		Stages: payment.PlanFeature{
			Included: payment.Quantity{Value: 0, Unit: payment.UnitStages},
			Exceedance: payment.Exceedances{
				PreProvision: payment.Exceedance{
					Quantity: payment.Quantity{Value: 1, Unit: payment.UnitStages},
					BasePrice: payment.BasePrice{
						Unit: payment.UnitStages,
						Value: 40.00,
						Currency: payment.CurrencyEUR,
						Interval: "month",
					},
				},
				OnDemand: payment.Exceedance{
					Quantity: payment.Quantity{Value: 1, Unit: payment.UnitStages},
					BasePrice: payment.BasePrice{
						Unit: payment.UnitStages,
						Value: 180.00,
						Currency: payment.CurrencyEUR,
						Interval: "month",
					},
				},
			},
		},
		Backups: payment.PlanFeatureBackup{
			MinimumInterval: payment.Quantity{Value: 1, Unit: payment.UnitMinutes},
		},
		Scaling: payment.PlanFeatureScaling{
			Minimum: payment.Quantity{Value: 1, Unit: payment.UnitPods},
			Maximum: payment.Quantity{Value: 1, Unit: payment.UnitPods},
			Exceedance: payment.ExceedancesWithoutOnDemand{
				PreProvision: payment.Exceedance{
					Quantity: payment.Quantity{Value: 1, Unit: payment.UnitPods},
					BasePrice: payment.BasePrice{
						Unit: payment.UnitPods,
						Value: 40.00,
						Currency: payment.CurrencyEUR,
						Interval: "month",
					},
				},
			},
		},
	},
}

func TestCanGenerateEstimationForPlanOnly(t *testing.T) {
	e := New([]payment.Plan{testPlanOne})
	p := Params{
		PlanID: testPlanOne.ID,
	}

	est, err := e.Estimate(p)
	if err != nil {
		t.Fatal(err)
	}

	myassert.CurrencyEquals(t, "EUR 99.00", est.MonthlyTotalCost)
	assert.Len(t, est.LineItems, 1)
	assert.Equal(t, est.LineItems[0].Quantity.Value, 1)
	myassert.CurrencyEquals(t, "EUR 99.00", est.LineItems[0].MonthlyCost)
}

func TestCanGenerateEstimationWithStorage(t *testing.T) {
	e := New([]payment.Plan{testPlanOne})
	p := Params{
		PlanID: testPlanOne.ID,
		Storage: 40 << 30,
	}

	est, err := e.Estimate(p)
	if err != nil {
		t.Fatal(err)
	}

	myassert.CurrencyEquals(t, "EUR 111.00", est.MonthlyTotalCost)
	assert.Len(t, est.LineItems, 2)
	assert.Equal(t, est.LineItems[0].Quantity.Value, 1)
	assert.Equal(t, est.LineItems[1].Quantity.Value, 40)

	myassert.CurrencyEquals(t, "EUR 99.00", est.LineItems[0].MonthlyCost)
	myassert.CurrencyEquals(t, "EUR 12.00", est.LineItems[1].MonthlyCost)
}

func TestCanGenerateEstimationWithStages(t *testing.T) {
	e := New([]payment.Plan{testPlanOne})
	p := Params{
		PlanID: testPlanOne.ID,
		Stages: 1,
	}

	est, err := e.Estimate(p)
	if err != nil {
		t.Fatal(err)
	}

	myassert.CurrencyEquals(t, "EUR 139.00", est.MonthlyTotalCost)
	assert.Len(t, est.LineItems, 2)
	assert.Equal(t, est.LineItems[0].Quantity.Value, 1)
	assert.Equal(t, est.LineItems[1].Quantity.Value, 1)

	myassert.CurrencyEquals(t, "EUR 99.00", est.LineItems[0].MonthlyCost)
	myassert.CurrencyEquals(t, "EUR 40.00", est.LineItems[1].MonthlyCost)
}

func TestCanGenerateEstimationWithMixedPreprovisionedAndOndemandStages(t *testing.T) {
	e := New([]payment.Plan{testPlanOne})
	p := Params{
		PlanID: testPlanOne.ID,
		Stages: 2,
		StagesOnDemand: 1,
	}

	est, err := e.Estimate(p)
	if err != nil {
		t.Fatal(err)
	}

	myassert.CurrencyEquals(t, "EUR 319.00", est.MonthlyTotalCost)
	assert.Len(t, est.LineItems, 3)
	assert.Equal(t, est.LineItems[0].Quantity.Value, 1)
	assert.Equal(t, est.LineItems[1].Quantity.Value, 1)
	assert.Equal(t, est.LineItems[2].Quantity.Value, 1)

	myassert.CurrencyEquals(t, "EUR 99.00", est.LineItems[0].MonthlyCost)
	myassert.CurrencyEquals(t, "EUR 40.00", est.LineItems[1].MonthlyCost)
	myassert.CurrencyEquals(t, "EUR 180.00", est.LineItems[2].MonthlyCost)
}

func TestEstimationContainsOnDemandCharges(t *testing.T) {
	e := New([]payment.Plan{testPlanOne})
	p := Params{
		PlanID: testPlanOne.ID,
	}

	est, err := e.Estimate(p)
	if err != nil {
		t.Fatal(err)
	}

	myassert.CurrencyEquals(t, "EUR 0.70", est.OnDemandCharges.Storage.Amount)
	myassert.CurrencyEquals(t, "EUR 180.00", est.OnDemandCharges.Stages.Amount)
}
package costestimator

type Params struct {
	PlanID         string
	Storage        uint64
	Stages         int
	StagesOnDemand int
	Scaling        int
}

type Estimator interface {
	Estimate(params Params) (*Estimation, error)
}

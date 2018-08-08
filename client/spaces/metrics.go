package spaces

import (
	"fmt"
	"net/url"
)

const (
	ScopeHour  = "hour"
	ScopeToday = "today"
	ScopeWeek  = "week"
	ScopeMonth = "month"
	ScopeYear  = "year"
)

func (c *spacesClient) GetComputeMetrics(spaceID string, scope string) (ComputeMetricPointList, error) {
	var space Space
	var metrics []ComputeMetricPoint

	err := c.client.Get("/spaces/"+spaceID, &space)
	if err != nil {
		return nil, fmt.Errorf("could not load stage: %s", err)
	}

	link, err := space.Links.GetLinkByRel("computeMetrics")
	if err != nil {
		return nil, fmt.Errorf("could not access metrics: %s", err)
	}

	query := url.Values{}
	query.Set("scope", scope)

	err = link.GetWithQuery(query, c.client, &metrics)
	if err != nil {
		return nil, fmt.Errorf("could not load metrics: %s", err)
	}

	return metrics, nil
}

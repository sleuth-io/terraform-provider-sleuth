package gqlclient

import (
	"errors"
	"fmt"

	"github.com/shurcooL/graphql"
)

// GetMetricImpactSource - Returns error impact source
func (c *Client) GetMetricImpactSource(projectSlug *string, slug *string) (*MetricImpactSource, error) {
	var query struct {
		Project struct {
			ImpactSources []struct {
				Type         graphql.String
				ImpactSource MetricImpactSource `graphql:"... on MetricImpactSource"`
			} `graphql:"impactSources(impactSourceSlug: $impactSourceSlug)"`
		} `graphql:"project(projectSlug: $projectSlug)"`
	}
	variables := map[string]interface{}{
		"projectSlug":      graphql.ID(*projectSlug),
		"impactSourceSlug": graphql.ID(*slug),
	}

	err := c.doQuery(&query, variables)

	if err != nil {
		return nil, err
	}

	for _, src := range query.Project.ImpactSources {
		if src.Type == "METRIC" {
			return &src.ImpactSource, nil
		}
	}
	return nil, nil
}

// CreateMetricImpactSource - Creates a environment
func (c *Client) CreateMetricImpactSource(input CreateMetricImpactSourceMutationInput) (*MetricImpactSource, error) {

	var m struct {
		CreateMetricImpactSource struct {
			ImpactSource MetricImpactSource
			Errors       ErrorsType
		} `graphql:"createMetricImpactSource(input: $input)"`
	}
	variables := map[string]interface{}{
		"input": input,
	}

	err := c.doMutate(&m, variables)

	if err != nil {
		return nil, err
	}

	if len(m.CreateMetricImpactSource.Errors) > 0 {
		return nil, errors.New(fmt.Sprintf("%s %+v", "Errors creating impact source: ", m.CreateMetricImpactSource.Errors))
	}
	return &m.CreateMetricImpactSource.ImpactSource, nil
}

// UpdateMetricImpactSource - Updates a environment
func (c *Client) UpdateMetricImpactSource(input UpdateMetricImpactSourceMutationInput) (*MetricImpactSource, error) {

	var m struct {
		UpdateMetricImpactSource struct {
			ImpactSource MetricImpactSource
			Errors       ErrorsType
		} `graphql:"updateMetricImpactSource(input: $input)"`
	}
	variables := map[string]interface{}{
		"input": input,
	}

	err := c.doMutate(&m, variables)

	if err != nil {
		return nil, err
	}

	if len(m.UpdateMetricImpactSource.Errors) > 0 {
		return nil, errors.New("Errors updating impact source")
	}

	return &m.UpdateMetricImpactSource.ImpactSource, nil
}

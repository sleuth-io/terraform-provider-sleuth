package gqlclient

import (
	"context"
	"errors"
	"fmt"

	"github.com/shurcooL/graphql"
)

// GetMetricImpactSource - Returns error impact source
func (c *Client) GetMetricImpactSource(ctx context.Context, projectSlug *string, slug *string) (*MetricImpactSource, error) {
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

	err := c.doQuery(ctx, &query, variables)

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
func (c *Client) CreateMetricImpactSource(ctx context.Context, input CreateMetricImpactSourceMutationInput) (*MetricImpactSource, error) {

	var m struct {
		CreateMetricImpactSource struct {
			ImpactSource MetricImpactSource
			Errors       ErrorsType
		} `graphql:"createMetricImpactSource(input: $input)"`
	}
	variables := map[string]interface{}{
		"input": input,
	}

	err := c.doMutate(ctx, &m, variables)

	if err != nil {
		return nil, err
	}

	if len(m.CreateMetricImpactSource.Errors) > 0 {
		return nil, fmt.Errorf("errors creating impact source: %+v", m.CreateMetricImpactSource.Errors)
	}
	return &m.CreateMetricImpactSource.ImpactSource, nil
}

// UpdateMetricImpactSource - Updates a environment
func (c *Client) UpdateMetricImpactSource(ctx context.Context, input UpdateMetricImpactSourceMutationInput) (*MetricImpactSource, error) {

	var m struct {
		UpdateMetricImpactSource struct {
			ImpactSource MetricImpactSource
			Errors       ErrorsType
		} `graphql:"updateMetricImpactSource(input: $input)"`
	}
	variables := map[string]interface{}{
		"input": input,
	}

	err := c.doMutate(ctx, &m, variables)

	if err != nil {
		return nil, err
	}

	if len(m.UpdateMetricImpactSource.Errors) > 0 {
		return nil, errors.New("Errors updating impact source")
	}

	return &m.UpdateMetricImpactSource.ImpactSource, nil
}

package gqlclient

import (
	"errors"
	"fmt"

	"github.com/shurcooL/graphql"
)

// GetErrorImpactSource - Returns error impact source
func (c *Client) GetErrorImpactSource(projectSlug *string, slug *string) (*ErrorImpactSource, error) {
	var query struct {
		Project struct {
			ImpactSources []struct {
				Type         graphql.String
				ImpactSource ErrorImpactSource `graphql:"... on ErrorImpactSource"`
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
		if src.Type == "ERROR" {
			return &src.ImpactSource, nil
		}
	}
	return nil, nil
}

// CreateErrorImpactSource - Creates a environment
func (c *Client) CreateErrorImpactSource(input CreateErrorImpactSourceMutationInput) (*ErrorImpactSource, error) {

	var m struct {
		CreateErrorImpactSource struct {
			ImpactSource ErrorImpactSource
			Errors       ErrorsType
		} `graphql:"createErrorImpactSource(input: $input)"`
	}
	variables := map[string]interface{}{
		"input": input,
	}

	err := c.doMutate(&m, variables)

	if err != nil {
		return nil, err
	}

	if len(m.CreateErrorImpactSource.Errors) > 0 {
		return nil, fmt.Errorf("errors creating impact source: %+v", m.CreateErrorImpactSource.Errors)
	}
	return &m.CreateErrorImpactSource.ImpactSource, nil
}

// UpdateErrorImpactSource - Updates a environment
func (c *Client) UpdateErrorImpactSource(input UpdateErrorImpactSourceMutationInput) (*ErrorImpactSource, error) {

	var m struct {
		UpdateErrorImpactSource struct {
			ImpactSource ErrorImpactSource
			Errors       ErrorsType
		} `graphql:"updateErrorImpactSource(input: $input)"`
	}
	variables := map[string]interface{}{
		"input": input,
	}

	err := c.doMutate(&m, variables)

	if err != nil {
		return nil, err
	}

	if len(m.UpdateErrorImpactSource.Errors) > 0 {
		return nil, errors.New("Errors updating impact source")
	}

	return &m.UpdateErrorImpactSource.ImpactSource, nil
}

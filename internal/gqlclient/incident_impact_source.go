package gqlclient

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/shurcooL/graphql"
)

// GetIncidentImpactSource returns incident impact source
func (c *Client) GetIncidentImpactSource(ctx context.Context, projectSlug, slug string) (*IncidentImpactSource, error) {
	var query struct {
		Project struct {
			ImpactSources []struct {
				ImpactSource IncidentImpactSource `graphql:"... on IncidentImpactSource" json:"impactSource"`
			} `graphql:"impactSources(impactSourceSlug: $impactSourceSlug)"`
		} `graphql:"project(projectSlug: $projectSlug)"`
	}
	variables := map[string]interface{}{
		"projectSlug":      graphql.ID(projectSlug),
		"impactSourceSlug": graphql.ID(slug),
	}

	err := c.doQuery(ctx, &query, variables)
	tflog.Error(ctx, fmt.Sprintf("Error when calling GraphQL, %+v", err))
	if err != nil {
		return nil, err
	}

	if len(query.Project.ImpactSources) > 1 {
		tflog.Warn(ctx, "More than one impact source found", map[string]interface{}{"slug": slug, "projectSlug": projectSlug})
	}
	return &query.Project.ImpactSources[0].ImpactSource, nil

}

func (c *Client) CreateIncidentImpactSource(ctx context.Context, input IncidentImpactSourceInputType) (*IncidentImpactSource, error) {
	var m struct {
		CreateIncidentImpactSource struct {
			ImpactSource IncidentImpactSource
			Errors       ErrorsType
		} `graphql:"createIncidentImpactSource(input: $input)"`
	}

	variables := map[string]interface{}{
		"input": input,
	}

	err := c.doMutate(ctx, &m, variables)

	if err != nil {
		return nil, err
	}

	if len(m.CreateIncidentImpactSource.Errors) > 0 {
		tflog.Error(ctx, fmt.Sprintf("%+v", m.CreateIncidentImpactSource.Errors))
		return nil, fmt.Errorf("Error creating incident impact source:  %+v", m.CreateIncidentImpactSource.Errors)
	}

	return &m.CreateIncidentImpactSource.ImpactSource, nil
}

func (c *Client) UpdateIncidentImpactSource(ctx context.Context, input IncidentImpactSourceInputUpdateType) (*IncidentImpactSource, error) {
	var m struct {
		UpdateIncidentImpactSource struct {
			ImpactSource IncidentImpactSource
			Errors       ErrorsType
		} `graphql:"updateIncidentImpactSource(input: $input)"`
	}

	variables := map[string]interface{}{
		"input": input,
	}

	err := c.doMutate(ctx, &m, variables)

	if err != nil {
		return nil, err
	}

	if len(m.UpdateIncidentImpactSource.Errors) > 0 {
		tflog.Error(ctx, fmt.Sprintf("%+v", m.UpdateIncidentImpactSource.Errors))
		return nil, fmt.Errorf("Error creating incident impact source:  %+v", m.UpdateIncidentImpactSource.Errors)
	}

	return &m.UpdateIncidentImpactSource.ImpactSource, nil
}

func (c *Client) DeleteIncidentImpactSource(ctx context.Context, input IncidentImpactSourceDeleteInputType) (bool, error) {
	tflog.Info(ctx, fmt.Sprintf("Deleting impact source %+v", input))
	var m struct {
		DeleteIncidentImpactSource struct {
			Success bool
		} `graphql:"deleteIncidentImpactSource(input: $input)"`
	}

	variables := map[string]interface{}{
		"input": input,
	}

	err := c.doMutate(ctx, &m, variables)

	if err != nil {
		return false, err
	}

	return m.DeleteIncidentImpactSource.Success, nil
}

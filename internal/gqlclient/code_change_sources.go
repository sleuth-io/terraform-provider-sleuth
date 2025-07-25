package gqlclient

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/shurcooL/graphql"
)

// ErrNoCodeChangeSourceFound indicates that a code change source with the given
// project slug & slug does not exist.
var ErrNoCodeChangeSourceFound = errors.New("code change source not found")

func (c *Client) GetCodeChangeSource(ctx context.Context, projectSlug *string, slug *string) (*CodeChangeSource, error) {
	var query struct {
		Project struct {
			ChangeSources []struct {
				Type         graphql.String
				ChangeSource CodeChangeSource `graphql:"... on CodeChangeSource"`
			}
		} `graphql:"project(projectSlug: $projectSlug)"`
	}
	variables := map[string]interface{}{
		"projectSlug": graphql.ID(*projectSlug),
	}

	err := c.doQuery(ctx, &query, variables)

	if err != nil {
		return nil, err
	}

	for _, src := range query.Project.ChangeSources {
		if src.Type == "CODE" {
			if src.ChangeSource.Slug == *slug {
				src.ChangeSource.Repository.Provider = strings.ToUpper(src.ChangeSource.Repository.Provider)
				// TODO: this should not be done here but we want to be consistent for now
				for idx, buildMapping := range src.ChangeSource.DeployTrackingBuildMappings {
					src.ChangeSource.DeployTrackingBuildMappings[idx].Provider = strings.ToUpper(buildMapping.Provider)
				}
				return &src.ChangeSource, nil
			}
		}
	}
	return nil, ErrNoCodeChangeSourceFound
}

func (c *Client) CreateCodeChangeSource(ctx context.Context, input CreateCodeChangeSourceMutationInput) (*CodeChangeSource, error) {

	var m struct {
		CreateCodeChangeSource struct {
			ChangeSource CodeChangeSource
			Errors       ErrorsType
		} `graphql:"createCodeChangeSource(input: $input)"`
	}
	input.DeployTrackingType = strings.ToUpper(input.DeployTrackingType)
	variables := map[string]interface{}{
		"input": input,
	}

	err := c.doMutate(ctx, &m, variables)

	if err != nil {
		return nil, err
	}

	if len(m.CreateCodeChangeSource.Errors) > 0 {
		return nil, fmt.Errorf("errors creating change source: %+v", m.CreateCodeChangeSource.Errors)
	}
	return &m.CreateCodeChangeSource.ChangeSource, nil
}

func (c *Client) UpdateCodeChangeSource(ctx context.Context, input UpdateCodeChangeSourceMutationInput) (*CodeChangeSource, error) {

	var m struct {
		UpdateCodeChangeSource struct {
			ChangeSource CodeChangeSource
			Errors       ErrorsType
		} `graphql:"updateCodeChangeSource(input: $input)"`
	}
	input.DeployTrackingType = strings.ToUpper(input.DeployTrackingType)
	variables := map[string]interface{}{
		"input": input,
	}

	err := c.doMutate(ctx, &m, variables)

	if err != nil {
		return nil, err
	}

	if len(m.UpdateCodeChangeSource.Errors) > 0 {
		return nil, fmt.Errorf("Errors updating code change source: %+v", m.UpdateCodeChangeSource.Errors)
	}
	return &m.UpdateCodeChangeSource.ChangeSource, nil
}

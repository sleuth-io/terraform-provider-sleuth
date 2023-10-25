package gqlclient

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/shurcooL/graphql"
)

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

	err := c.doQuery(&query, variables)

	if err != nil {
		return nil, err
	}

	for _, src := range query.Project.ChangeSources {
		if src.Type == "CODE" {
			if src.ChangeSource.Slug == *slug {
				src.ChangeSource.Repository.Provider = strings.ToUpper(src.ChangeSource.Repository.Provider)
				// Sort mappings based on order passed to have consistent results in state
				sort.Slice(src.ChangeSource.DeployTrackingBuildMappings, func(i, j int) bool {
					return src.ChangeSource.DeployTrackingBuildMappings[i].Order < src.ChangeSource.DeployTrackingBuildMappings[j].Order
				})
				// TODO: this should not be done here but we want to be consistent for now
				for idx, buildMapping := range src.ChangeSource.DeployTrackingBuildMappings {
					src.ChangeSource.DeployTrackingBuildMappings[idx].Provider = strings.ToUpper(buildMapping.Provider)
				}
				return &src.ChangeSource, nil
			}
		}
	}

	return nil, nil
}

func (c *Client) CreateCodeChangeSource(input CreateCodeChangeSourceMutationInput) (*CodeChangeSource, error) {

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

	err := c.doMutate(&m, variables)

	if err != nil {
		return nil, err
	}

	if len(m.CreateCodeChangeSource.Errors) > 0 {
		return nil, errors.New(fmt.Sprintf("%s %+v", "Errors creating change source: ", m.CreateCodeChangeSource.Errors))
	}
	return &m.CreateCodeChangeSource.ChangeSource, nil
}

func (c *Client) UpdateCodeChangeSource(input UpdateCodeChangeSourceMutationInput) (*CodeChangeSource, error) {

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

	err := c.doMutate(&m, variables)

	if err != nil {
		return nil, err
	}

	if len(m.UpdateCodeChangeSource.Errors) > 0 {
		return nil, errors.New("Errors updating code change source")
	}
	return &m.UpdateCodeChangeSource.ChangeSource, nil
}

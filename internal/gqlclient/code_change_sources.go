package gqlclient

import (
	"errors"
	"fmt"
	"github.com/shurcooL/graphql"
	"strings"
)

func (c *Client) GetCodeChangeSource(projectSlug *string, slug *string) (*CodeChangeSource, error) {
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
				for idx, mapping := range src.ChangeSource.EnvironmentMappings {
					src.ChangeSource.EnvironmentMappings[idx].EnvironmentSlug = fmt.Sprintf("%s/%s", *projectSlug, mapping.EnvironmentSlug)
				}
				return &src.ChangeSource, nil
			}
		}
	}
	return nil, errors.New("Code change source not found")
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
	for idx, mapping := range m.CreateCodeChangeSource.ChangeSource.EnvironmentMappings {
		m.CreateCodeChangeSource.ChangeSource.EnvironmentMappings[idx].EnvironmentSlug = fmt.Sprintf("%s/%s", input.ProjectSlug, mapping.EnvironmentSlug)
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
	for idx, mapping := range m.UpdateCodeChangeSource.ChangeSource.EnvironmentMappings {
		m.UpdateCodeChangeSource.ChangeSource.EnvironmentMappings[idx].EnvironmentSlug = fmt.Sprintf("%s/%s", input.ProjectSlug, mapping.EnvironmentSlug)
	}
	return &m.UpdateCodeChangeSource.ChangeSource, nil
}

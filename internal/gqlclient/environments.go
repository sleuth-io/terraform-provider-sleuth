package gqlclient

import (
	"errors"
	"github.com/shurcooL/graphql"
)


func (c *Client) GetEnvironmentByName(projectSlug *string, name *string) (*Environment, error) {
	var query struct {
		Project struct {
			Environments []Environment
		}`graphql:"project(orgSlug: $orgSlug, projectSlug: $projectSlug)"`
	}
	variables := map[string]interface{}{
		"orgSlug":   graphql.ID(c.OrgSlug),
		"projectSlug": graphql.ID(*projectSlug),
	}

	err := c.doQuery(&query, variables)

	if err != nil {
		return nil, err
	}

	for _, env := range query.Project.Environments {
		if env.Name == *name {
			return &env, nil
		}
	}
	return nil, errors.New("Not found")
}


// GetEnvironment - Returns environment
func (c *Client) GetEnvironment(projectSlug *string, slug *string) (*Environment, error) {
	var query struct {
		Project struct {
			Environments []Environment
		}`graphql:"project(orgSlug: $orgSlug, projectSlug: $projectSlug)"`
	}
	variables := map[string]interface{}{
		"orgSlug":   graphql.ID(c.OrgSlug),
		"projectSlug": graphql.ID(*projectSlug),
	}

	err := c.doQuery(&query, variables)

	if err != nil {
		return nil, err
	}

	for _, env := range query.Project.Environments {
		if env.Slug == *slug {
			return &env, nil
		}
	}
	return nil, errors.New("Not found")
}

// CreateEnvironment - Creates a environment
func (c *Client) CreateEnvironment(projectSlug *string, input CreateEnvironmentMutationInput) (*Environment, error) {

	var m struct {
		CreateEnvironment struct {
			Environment Environment
			Errors ErrorsType
		} `graphql:"createEnvironment(orgSlug: $orgSlug, projectSlug: $projectSlug, input: $input)"`
	}
	variables := map[string]interface{}{
		"orgSlug": graphql.ID(c.OrgSlug),
		"projectSlug": graphql.ID(*projectSlug),
		"input":   input,
	}

	err := c.doMutate(&m, variables)

	if err != nil {
		return nil, err
	}

	if len(m.CreateEnvironment.Errors) > 0 {
		return nil, errors.New("Errors creating environment")
	}
	return &m.CreateEnvironment.Environment, nil
}

// UpdateEnvironment - Updates a environment
func (c *Client) UpdateEnvironment(projectSlug *string, slug *string, input UpdateEnvironmentMutationInput) (*Environment, error) {

	var m struct {
		UpdateEnvironment struct {
			Environment Environment
			Errors ErrorsType
		} `graphql:"updateEnvironment(orgSlug: $orgSlug, projectSlug: $projectSlug, slug: $slug, input: $input)"`
	}
	variables := map[string]interface{}{
		"orgSlug":   graphql.ID(c.OrgSlug),
		"input":   input,
		"projectSlug": graphql.ID(*projectSlug),
		"slug": graphql.ID(*slug),
	}

	err := c.doMutate(&m, variables)

	if err != nil {
		return nil, err
	}

	if len(m.UpdateEnvironment.Errors) > 0 {
		return nil, errors.New("Errors updating environment")
	}

	return &m.UpdateEnvironment.Environment, nil
}


// DeleteEnvironment - Deletes a environment
func (c *Client) DeleteEnvironment(projectSlug *string, slug *string) error {

	var m struct {
		DeleteEnvironment struct {
			Success graphql.Boolean
		} `graphql:"deleteEnvironment(slug: $slug, projectSlug: $projectSlug, orgSlug: $orgSlug)"`
	}
	variables := map[string]interface{}{
		"orgSlug":   graphql.ID(c.OrgSlug),
		"projectSlug": 	graphql.ID(*projectSlug),
		"slug":   graphql.ID(*slug),
	}

	err := c.doMutate(&m, variables)

	if err != nil {
		return err
	}

	if !m.DeleteEnvironment.Success {
		return errors.New("Missing")
	} else {
		return nil
	}
}

package gqlclient

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/shurcooL/graphql"
)

var ErrNotFound = errors.New("Resource was not found")

func (c *Client) GetEnvironmentByName(ctx context.Context, projectSlug *string, name *string) (*Environment, error) {
	var query struct {
		Project struct {
			Environments []Environment
		} `graphql:"project(projectSlug: $projectSlug)"`
	}
	variables := map[string]interface{}{
		"projectSlug": graphql.ID(*projectSlug),
	}

	err := c.doQuery(ctx, &query, variables)

	if err != nil {
		return nil, err
	}

	for _, env := range query.Project.Environments {
		if env.Name == *name {
			return &env, nil
		}
	}
	return nil, ErrNotFound
}

// GetEnvironment - Returns environment
func (c *Client) GetEnvironment(ctx context.Context, projectSlug *string, slug *string) (*Environment, error) {
	var query struct {
		Project struct {
			Environments []Environment
		} `graphql:"project(projectSlug: $projectSlug)"`
	}
	variables := map[string]interface{}{
		"projectSlug": graphql.ID(*projectSlug),
	}

	err := c.doQuery(ctx, &query, variables)

	if err != nil {
		return nil, err
	}

	for _, env := range query.Project.Environments {
		if env.Slug == *slug {
			return &env, nil
		}
	}
	return nil, nil
}

// CreateEnvironment - Creates a environment
func (c *Client) CreateEnvironment(ctx context.Context, input CreateEnvironmentMutationInput) (*Environment, error) {

	var m struct {
		CreateEnvironment struct {
			Environment Environment
			Errors      ErrorsType
		} `graphql:"createEnvironment(input: $input)"`
	}
	variables := map[string]interface{}{
		"input": input,
	}

	err := c.doMutate(ctx, &m, variables)

	if err != nil {
		return nil, err
	}

	if len(m.CreateEnvironment.Errors) > 0 {
		return nil, fmt.Errorf("errors creating environment: %+v", m.CreateEnvironment.Errors)
	}
	return &m.CreateEnvironment.Environment, nil
}

// UpdateEnvironment - Updates a environment
func (c *Client) UpdateEnvironment(ctx context.Context, input UpdateEnvironmentMutationInput) (*Environment, error) {

	var m struct {
		UpdateEnvironment struct {
			Environment Environment
			Errors      ErrorsType
		} `graphql:"updateEnvironment(input: $input)"`
	}
	variables := map[string]interface{}{
		"input": input,
	}

	err := c.doMutate(ctx, &m, variables)

	if err != nil {
		return nil, err
	}

	if len(m.UpdateEnvironment.Errors) > 0 {
		return nil, errors.New("Errors updating environment")
	}

	return &m.UpdateEnvironment.Environment, nil
}

// DeleteEnvironment - Deletes a environment
func (c *Client) DeleteEnvironment(ctx context.Context, projectSlug *string, slug *string) error {
	var m struct {
		DeleteEnvironment struct {
			Success graphql.Boolean
			Errors  ErrorsType
		} `graphql:"deleteEnvironment(input: $input)"`
	}
	variables := map[string]interface{}{
		"input": DeleteEnvironmentMutationInput{ProjectSlug: *projectSlug, Slug: *slug},
	}

	err := c.doMutate(ctx, &m, variables)

	if err != nil {
		return err
	}

	tflog.Info(ctx, "DeleteEnvironment result", map[string]interface{}{"errors": fmt.Sprintf("%+v", m.DeleteEnvironment.Errors)})

	for _, err := range m.DeleteEnvironment.Errors {
		// If the error is that the environment is the default environment, then we can ignore it
		if err.Field == "slug" && err.Messages[0] == "You can not delete your default environment" {
			return nil
		}
	}

	if !m.DeleteEnvironment.Success {
		return fmt.Errorf("deleting environment was not successful: %+v", m.DeleteEnvironment.Errors)
	} else {
		return nil
	}
}

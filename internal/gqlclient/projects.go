package gqlclient

import (
	"errors"
	"github.com/shurcooL/graphql"
)

//// GetProjects - Returns list of projects
//func (c *Client) GetProjects() ([]Project, error) {
//	var query struct {
//		Projects []struct {
//			Name graphql.String
//			Slug graphql.String
//		} `graphql:"projects(orgSlug: $orgSlug)"`
//	}
//	variables := map[string]interface{}{
//		"orgSlug":   graphql.ID(c.OrgSlug),
//	}
//
//	err := c.doQuery(&query, variables)
//
//	if err != nil {
//		return nil, err
//	}
//
//	projects := []Project{}
//	for _, element  := range query.Projects{
//		project := Project{Name: string(element.Slug)}
//		projects = append(projects, project)
//	}
//
//	return projects, nil
//}

// GetProject - Returns project
func (c *Client) GetProject(slug *string) (*Project, error) {
	var query struct {
		Project Project `graphql:"project(projectSlug: $projectSlug)"`
	}
	variables := map[string]interface{}{
		"projectSlug": graphql.ID(*slug),
	}

	err := c.doQuery(&query, variables)

	if err != nil {
		return nil, err
	}

	return &query.Project, nil
}

// CreateProject - Creates a project
func (c *Client) CreateProject(input CreateProjectMutationInput) (*Project, error) {

	var m struct {
		CreateProject struct {
			Project Project
			Errors  ErrorsType
		} `graphql:"createProject(input: $input)"`
	}
	variables := map[string]interface{}{
		"input": input,
	}

	err := c.doMutate(&m, variables)

	if err != nil {
		return nil, err
	}

	if len(m.CreateProject.Errors) > 0 {
		return nil, errors.New("Errors creating project")
	}
	return &m.CreateProject.Project, nil
}

// UpdateProject - Updates a project
func (c *Client) UpdateProject(slug *string, input UpdateProjectMutationInput) (*Project, error) {

	var m struct {
		UpdateProject struct {
			Project Project
			Errors  ErrorsType
		} `graphql:"updateProject(input: $input)"`
	}
	variables := map[string]interface{}{
		"input": input,
	}

	err := c.doMutate(&m, variables)

	if err != nil {
		return nil, err
	}

	if len(m.UpdateProject.Errors) > 0 {
		return nil, errors.New("Errors updating project")
	}

	return &m.UpdateProject.Project, nil
}

// DeleteProject - Deletes a project
func (c *Client) DeleteProject(slug *string) error {

	var m struct {
		DeleteProject struct {
			Success graphql.Boolean
		} `graphql:"deleteProject(input: $input)"`
	}
	variables := map[string]interface{}{
		"input": DeleteProjectMutationInput{Slug: *slug},
	}

	err := c.doMutate(&m, variables)

	if err != nil {
		return err
	}

	if !m.DeleteProject.Success {
		return errors.New("Missing")
	} else {
		return nil
	}
}

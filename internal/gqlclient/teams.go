package gqlclient

import (
	"context"
	"fmt"
	"strings"

	"github.com/shurcooL/graphql"
)

// GetTeam - Returns team
func (c *Client) GetTeam(ctx context.Context, slug *string) (*Team, error) {
	var query struct {
		Team Team `graphql:"team(teamSlug: $teamSlug)"`
	}
	variables := map[string]interface{}{
		"teamSlug": graphql.ID(*slug),
	}
	err := c.doQuery(ctx, &query, variables)
	if err != nil {
		if strings.HasSuffix(err.Error(), "not found") {
			return nil, nil
		}
		return nil, err
	}
	return &query.Team, nil
}

// CreateTeam - Creates a team
func (c *Client) CreateTeam(ctx context.Context, input CreateTeamMutationInput) (*Team, error) {
	var m struct {
		CreateTeam struct {
			Team   Team
			Errors ErrorsType
		} `graphql:"createTeam(input: $input)"`
	}
	variables := map[string]interface{}{
		"input": input,
	}

	err := c.doMutate(ctx, &m, variables)
	if err != nil {
		return nil, err
	}

	if len(m.CreateTeam.Errors) > 0 {
		return nil, fmt.Errorf("%+v", m.CreateTeam.Errors)
	}
	return &m.CreateTeam.Team, nil
}

// UpdateTeam - Updates a team
func (c *Client) UpdateTeam(ctx context.Context, slug *string, input UpdateTeamMutationInput) (*Team, error) {
	var m struct {
		UpdateTeam struct {
			Team   Team
			Errors ErrorsType
		} `graphql:"updateTeam(input: $input)"`
	}
	variables := map[string]interface{}{
		"input": input,
	}

	err := c.doMutate(ctx, &m, variables)
	if err != nil {
		return nil, err
	}

	if len(m.UpdateTeam.Errors) > 0 {
		return nil, fmt.Errorf("%+v", m.UpdateTeam.Errors)
	}
	return &m.UpdateTeam.Team, nil
}

// DeleteTeam - Deletes a team
func (c *Client) DeleteTeam(ctx context.Context, slug *string) error {
	var m struct {
		DeleteTeam struct {
			Success graphql.Boolean
		} `graphql:"deleteTeam(input: $input)"`
	}
	variables := map[string]interface{}{
		"input": DeleteTeamMutationInput{Slug: *slug},
	}

	err := c.doMutate(ctx, &m, variables)
	if err != nil {
		return err
	}

	if !m.DeleteTeam.Success {
		return fmt.Errorf("Missing")
	} else {
		return nil
	}
}

// AddTeamMembers - Adds members to a team
func (c *Client) AddTeamMembers(ctx context.Context, input AddTeamMembersMutationInput) error {
	var m struct {
		AddTeamMembers struct {
			Success bool
			Errors  ErrorsType
		} `graphql:"addTeamMembers(input: $input)"`
	}
	variables := map[string]interface{}{
		"input": input,
	}

	err := c.doMutate(ctx, &m, variables)
	if err != nil {
		return err
	}
	if !m.AddTeamMembers.Success {
		return fmt.Errorf("errors adding team members: %+v", m.AddTeamMembers.Errors)
	}
	return nil
}

// RemoveTeamMembers - Removes members from a team
func (c *Client) RemoveTeamMembers(ctx context.Context, input RemoveTeamMembersMutationInput) error {
	var m struct {
		RemoveTeamMembers struct {
			Success bool
			Errors  ErrorsType
		} `graphql:"removeTeamMembers(input: $input)"`
	}
	variables := map[string]interface{}{
		"input": input,
	}

	err := c.doMutate(ctx, &m, variables)
	if err != nil {
		return err
	}
	if !m.RemoveTeamMembers.Success {
		return fmt.Errorf("errors removing team members: %+v", m.RemoveTeamMembers.Errors)
	}
	return nil
}

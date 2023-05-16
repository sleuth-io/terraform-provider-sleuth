package gqlclient

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/shurcooL/graphql"
	"sort"
	"strconv"
)

type Int64Slice []int64

func (x Int64Slice) Len() int           { return len(x) }
func (x Int64Slice) Less(i, j int) bool { return x[i] < x[j] }
func (x Int64Slice) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }

type GQLClient interface {
	GetTeam(context.Context, string) (*Team, error)
}

func (c *Client) GetTeam(ctx context.Context, slug string) (*Team, error) {
	tflog.Info(ctx, "Fetching from remote", map[string]any{"slug": slug})
	var query struct {
		Team TeamMutations `graphql:"team(teamSlug: $teamSlug)"`
	}
	variables := map[string]interface{}{
		"teamSlug": graphql.ID(slug),
	}
	err := c.doQuery(&query, variables)
	if err != nil {
		return nil, err
	}
	var memberIDs Int64Slice
	for _, val := range query.Team.Members.Objects {
		mID, _ := strconv.Atoi(val.ID)
		memberIDs = append(memberIDs, int64(mID))
	}
	// sort IDs to keep them in the same order at all times
	sort.Sort(memberIDs)

	t := Team{
		query.Team.Name,
		query.Team.Slug,
		memberIDs,
	}

	return &t, nil
}

func (c *Client) CreateTeam(ctx context.Context, input CreateTeamMutationInput) (*Team, error) {
	tflog.Info(ctx, "Create team", map[string]any{"input": fmt.Sprintf("%+v", input)})
	var m struct {
		CreateTeam struct {
			Team TeamMutations
		} `graphql:"createTeam(input: $input)"`
	}
	variables := map[string]interface{}{
		"input": input,
	}

	err := c.doMutate(&m, variables)
	tflog.Info(ctx, "Create team result", map[string]any{"result": fmt.Sprintf("%+v", m)})
	if err != nil {
		return nil, err
	}
	var memberIDs Int64Slice
	for _, val := range m.CreateTeam.Team.Members.Objects {
		mID, _ := strconv.Atoi(val.ID)
		memberIDs = append(memberIDs, int64(mID))
	}

	sort.Sort(memberIDs)

	t := Team{
		m.CreateTeam.Team.Name,
		m.CreateTeam.Team.Slug,
		memberIDs,
	}
	return &t, nil
}

func (c *Client) UpdateTeam(ctx context.Context, input UpdateTeamMutationInput) (*Team, error) {
	tflog.Info(ctx, "Update team", map[string]any{"input": fmt.Sprintf("%+v", input)})
	var m struct {
		UpdateTeam struct {
			Team TeamMutations
		} `graphql:"updateTeam(input: $input)"`
	}
	variables := map[string]interface{}{
		"input": input,
	}

	err := c.doMutate(&m, variables)
	tflog.Info(ctx, "Update team result", map[string]any{"result": fmt.Sprintf("%+v", m)})
	if err != nil {
		return nil, err
	}
	var memberIDs Int64Slice
	for _, val := range m.UpdateTeam.Team.Members.Objects {
		mID, _ := strconv.Atoi(val.ID)
		memberIDs = append(memberIDs, int64(mID))
	}

	sort.Sort(memberIDs)

	t := Team{
		m.UpdateTeam.Team.Name,
		m.UpdateTeam.Team.Slug,
		memberIDs,
	}
	return &t, nil
}

func (c *Client) DeleteTeam(ctx context.Context, input DeleteTeamMutationInput) error {
	tflog.Info(ctx, "Delete team", map[string]any{"input": fmt.Sprintf("%+v", input)})
	var m struct {
		DeleteTeam struct {
			Success bool `json:"success"`
		} `graphql:"deleteTeam(input: $input)"`
	}
	variables := map[string]interface{}{
		"input": input,
	}

	err := c.doMutate(&m, variables)
	tflog.Info(ctx, "Delete team result", map[string]any{"result": fmt.Sprintf("%+v", m)})
	if err != nil {
		return err
	}
	return nil
}

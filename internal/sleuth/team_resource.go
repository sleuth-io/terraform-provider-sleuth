package sleuth

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/shurcooL/graphql"
	"github.com/sleuth-io/terraform-provider-sleuth/internal/gqlclient"
)

var (
	_ resource.Resource                = &teamResource{}
	_ resource.ResourceWithConfigure   = &teamResource{}
	_ resource.ResourceWithImportState = &teamResource{}
)

type teamResourceModel struct {
	ID         types.String `tfsdk:"id"`
	Name       types.String `tfsdk:"name"`
	Slug       types.String `tfsdk:"slug"`
	ParentSlug types.String `tfsdk:"parent_slug"`
	Members    types.List   `tfsdk:"members"`
}

type teamResource struct {
	c *gqlclient.Client
}

func NewTeamResource() resource.Resource {
	return &teamResource{}
}

func (t *teamResource) Schema(_ context.Context, _ resource.SchemaRequest, res *resource.SchemaResponse) {
	res.Schema = schema.Schema{
		Description: "Team resource manages Sleuth teams and subteams.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Team name",
				Required:            true,
			},
			"slug": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Team slug",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"parent_slug": schema.StringAttribute{
				MarkdownDescription: "Parent team slug (for subteams)",
				Optional:            true,
				Computed:            true,
			},
			"members": schema.ListAttribute{
				Description: "List of user emails to be members of the team.",
				ElementType: basetypes.StringType{},
				Optional:    true,
			},
		},
	}
}

func (t *teamResource) Configure(ctx context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	t.c = req.ProviderData.(*gqlclient.Client)
}

func (t *teamResource) Metadata(ctx context.Context, req resource.MetadataRequest, res *resource.MetadataResponse) {
	res.TypeName = req.ProviderTypeName + "_team"
}

func (t *teamResource) Create(ctx context.Context, req resource.CreateRequest, res *resource.CreateResponse) {
	var plan teamResourceModel
	diags := req.Plan.Get(ctx, &plan)
	res.Diagnostics.Append(diags...)
	if res.Diagnostics.HasError() {
		return
	}

	var parent *string
	if !plan.ParentSlug.IsNull() && plan.ParentSlug.ValueString() != "" {
		s := plan.ParentSlug.ValueString()
		parent = &s
	}

	input := gqlclient.CreateTeamMutationInput{
		Name:   plan.Name.ValueString(),
		Parent: parent,
	}

	team, err := t.c.CreateTeam(ctx, input)
	if err != nil {
		res.Diagnostics.AddError(
			"Error creating team",
			err.Error(),
		)
		return
	}

	// Handle members
	if !plan.Members.IsNull() && plan.Members.Elements() != nil && len(plan.Members.Elements()) > 0 {
		var emails []string
		for _, v := range plan.Members.Elements() {
			email := v.(types.String).ValueString()
			emails = append(emails, email)
		}
		users, err := getUsersByEmails(t.c, emails)
		if err != nil {
			res.Diagnostics.AddError("Error resolving user emails", err.Error())
			return
		}
		var userIDs []string
		for _, u := range users {
			userIDs = append(userIDs, u.ID)
		}
		if len(userIDs) > 0 {
			addInput := gqlclient.AddTeamMembersMutationInput{
				Slug:    team.Slug,
				Members: userIDs,
			}
			err = t.c.AddTeamMembers(ctx, addInput)
			if err != nil {
				res.Diagnostics.AddError("Error adding team members", err.Error())
				return
			}
		}
	}

	// Fetch actual members from API for state
	emails, err := getTeamMemberEmails(t.c, team.Slug)
	if err != nil {
		res.Diagnostics.AddError("Error fetching team members after create", err.Error())
		return
	}
	membersNull := plan.Members.IsNull()
	state := getNewStateFromTeam(team, emails, membersNull)
	res.Diagnostics.Append(res.State.Set(ctx, state)...)
}

func (t *teamResource) Read(ctx context.Context, req resource.ReadRequest, res *resource.ReadResponse) {
	var state teamResourceModel
	diags := req.State.Get(ctx, &state)
	res.Diagnostics.Append(diags...)
	if res.Diagnostics.HasError() {
		return
	}

	slug := state.Slug.ValueString()
	team, err := t.c.GetTeam(ctx, &slug)
	if err != nil {
		res.Diagnostics.AddError("Error reading team", err.Error())
		return
	}
	if team == nil {
		res.State.RemoveResource(ctx)
		return
	}
	emails, err := getTeamMemberEmails(t.c, team.Slug)
	if err != nil {
		res.Diagnostics.AddError("Error fetching team members", err.Error())
		return
	}
	membersNull := state.Members.IsNull()
	newState := getNewStateFromTeam(team, emails, membersNull)
	res.Diagnostics.Append(res.State.Set(ctx, newState)...)
}

func (t *teamResource) Update(ctx context.Context, req resource.UpdateRequest, res *resource.UpdateResponse) {
	var plan, state teamResourceModel
	diags := req.Plan.Get(ctx, &plan)
	res.Diagnostics.Append(diags...)
	diags = req.State.Get(ctx, &state)
	res.Diagnostics.Append(diags...)
	if res.Diagnostics.HasError() {
		return
	}

	slug := state.Slug.ValueString()

	updateNeeded := plan.Name.ValueString() != state.Name.ValueString() || plan.ParentSlug.ValueString() != state.ParentSlug.ValueString()

	var updatedTeam *gqlclient.Team
	var err error
	if updateNeeded {
		var parent *string
		if !plan.ParentSlug.IsNull() && plan.ParentSlug.ValueString() != "" {
			s := plan.ParentSlug.ValueString()
			parent = &s
		}
		var name *string
		if plan.Name.ValueString() != "" {
			n := plan.Name.ValueString()
			name = &n
		}
		input := gqlclient.UpdateTeamMutationInput{
			Slug:   slug,
			Name:   name,
			Parent: parent,
		}
		updatedTeam, err = t.c.UpdateTeam(ctx, &slug, input)
		if err != nil {
			res.Diagnostics.AddError("Error updating team", err.Error())
			return
		}
	}

	// Handle members
	var oldEmails, newEmails []string
	if !state.Members.IsNull() && state.Members.Elements() != nil {
		for _, v := range state.Members.Elements() {
			oldEmails = append(oldEmails, v.(types.String).ValueString())
		}
	}
	if !plan.Members.IsNull() && plan.Members.Elements() != nil {
		for _, v := range plan.Members.Elements() {
			newEmails = append(newEmails, v.(types.String).ValueString())
		}
	}

	// Compute additions and removals
	oldSet := make(map[string]struct{})
	newSet := make(map[string]struct{})
	for _, e := range oldEmails {
		oldSet[e] = struct{}{}
	}
	for _, e := range newEmails {
		newSet[e] = struct{}{}
	}
	var toAdd, toRemove []string
	for e := range newSet {
		if _, found := oldSet[e]; !found {
			toAdd = append(toAdd, e)
		}
	}
	for e := range oldSet {
		if _, found := newSet[e]; !found {
			toRemove = append(toRemove, e)
		}
	}

	if len(toAdd) > 0 {
		users, err := getUsersByEmails(t.c, toAdd)
		if err != nil {
			res.Diagnostics.AddError("Error resolving user emails to add", err.Error())
			return
		}
		var userIDs []string
		for _, u := range users {
			userIDs = append(userIDs, u.ID)
		}
		if len(userIDs) > 0 {
			addInput := gqlclient.AddTeamMembersMutationInput{
				Slug:    state.Slug.ValueString(),
				Members: userIDs,
			}
			err = t.c.AddTeamMembers(ctx, addInput)
			if err != nil {
				res.Diagnostics.AddError("Error adding team members", err.Error())
				return
			}
		}
	}
	if len(toRemove) > 0 {
		users, err := getUsersByEmails(t.c, toRemove)
		if err != nil {
			res.Diagnostics.AddError("Error resolving user emails to remove", err.Error())
			return
		}
		var userIDs []string
		for _, u := range users {
			userIDs = append(userIDs, u.ID)
		}
		if len(userIDs) > 0 {
			removeInput := gqlclient.RemoveTeamMembersMutationInput{
				Slug:    state.Slug.ValueString(),
				Members: userIDs,
			}
			err = t.c.RemoveTeamMembers(ctx, removeInput)
			if err != nil {
				res.Diagnostics.AddError("Error removing team members", err.Error())
				return
			}
		}
	}

	// Set the new state
	if updatedTeam == nil {
		updatedTeam, err = t.c.GetTeam(ctx, &slug)
		if err != nil {
			res.Diagnostics.AddError("Error reading team after update", err.Error())
			return
		}
	}
	emails, err := getTeamMemberEmails(t.c, updatedTeam.Slug)
	if err != nil {
		res.Diagnostics.AddError("Error fetching team members after update", err.Error())
		return
	}
	membersNull := plan.Members.IsNull()
	newState := getNewStateFromTeam(updatedTeam, emails, membersNull)
	if updatedTeam != nil {
		newState.ID = types.StringValue(updatedTeam.ID)
		newState.Name = types.StringValue(updatedTeam.Name)
		newState.Slug = types.StringValue(updatedTeam.Slug)
	}
	res.Diagnostics.Append(res.State.Set(ctx, newState)...)
}

func (t *teamResource) Delete(ctx context.Context, req resource.DeleteRequest, res *resource.DeleteResponse) {
	var state teamResourceModel
	diags := req.State.Get(ctx, &state)
	res.Diagnostics.Append(diags...)
	if res.Diagnostics.HasError() {
		return
	}

	slug := state.Slug.ValueString()
	err := t.c.DeleteTeam(ctx, &slug)
	if err != nil {
		res.Diagnostics.AddError("Error deleting team", err.Error())
		return
	}
}

func (t *teamResource) ImportState(ctx context.Context, req resource.ImportStateRequest, res *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("slug"), req, res)
}

// getUsersByEmails returns users matching the given emails
func getUsersByEmails(c *gqlclient.Client, emails []string) ([]gqlclient.User, error) {
	if len(emails) == 0 {
		return nil, nil
	}
	var query struct {
		Organization struct {
			Users struct {
				Objects []struct {
					ID    string `graphql:"id"`
					Email string `graphql:"email"`
				} `graphql:"objects"`
			} `graphql:"users(term: $term, emails: $emails, page: $page, pageSize: $pageSize)"`
		} `graphql:"organization(orgSlug: $orgSlug)"`
	}
	var emailVars []graphql.String
	for _, e := range emails {
		emailVars = append(emailVars, graphql.String(e))
	}
	variables := map[string]interface{}{
		"orgSlug":  graphql.ID(c.OrgSlug),
		"term":     graphql.String(""),
		"emails":   emailVars,
		"page":     graphql.Int(1),
		"pageSize": graphql.Int(50),
	}
	if err := c.GQLClient.Query(context.Background(), &query, variables); err != nil {
		return nil, err
	}
	var users []gqlclient.User
	for _, u := range query.Organization.Users.Objects {
		users = append(users, gqlclient.User{ID: u.ID, Email: u.Email})
	}
	return users, nil
}

// Add helper functions for state/model conversion
func getNewStateFromTeam(team *gqlclient.Team, memberEmails []string, membersNull bool) teamResourceModel {
	parentSlug := types.StringNull()
	if team.Parent != nil && team.Parent.Slug != "" {
		parentSlug = types.StringValue(team.Parent.Slug)
	}
	var membersList basetypes.ListValue
	if membersNull {
		membersList = types.ListNull(types.StringType)
	} else {
		var elems []attr.Value
		for _, email := range memberEmails {
			elems = append(elems, types.StringValue(email))
		}
		membersList, _ = types.ListValue(types.StringType, elems)
	}
	return teamResourceModel{
		ID:         types.StringValue(team.ID),
		Name:       types.StringValue(team.Name),
		Slug:       types.StringValue(team.Slug),
		ParentSlug: parentSlug,
		Members:    membersList,
	}
}

func getTeamMemberEmails(c *gqlclient.Client, slug string) ([]string, error) {
	var query struct {
		Team struct {
			Members struct {
				Objects []gqlclient.User `graphql:"objects"`
			} `graphql:"members(page: $page, pageSize: $pageSize)"`
		} `graphql:"team(teamSlug: $teamSlug)"`
	}
	variables := map[string]interface{}{
		"teamSlug": graphql.ID(slug),
		"page":     graphql.Int(1),
		"pageSize": graphql.Int(50),
	}
	if err := c.GQLClient.Query(context.Background(), &query, variables); err != nil {
		return nil, err
	}
	var emails []string
	for _, member := range query.Team.Members.Objects {
		emails = append(emails, member.Email)
	}
	return emails, nil
}

package sleuth

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/sleuth-io/terraform-provider-sleuth/internal/gqlclient"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &teamResource{}
	_ resource.ResourceWithConfigure = &teamResource{}
)

type teamResourceModel struct {
	Name    types.String `tfsdk:"name"`
	Members types.List   `tfsdk:"members"`
	Slug    types.String `tfsdk:"slug"`
}

type teamResourceValues struct {
	Name    types.String `tfsdk:"name"`
	Members types.List   `tfsdk:"members"`
	Slug    types.String `tfsdk:"slug"`
}

// NewTeamResource is a helper function to simplify the provider implementation.
func NewTeamResource() resource.Resource {
	return &teamResource{}
}

// teamResource is the resource implementation.
type teamResource struct {
	c *gqlclient.Client
}

// Metadata returns the resource type name.
func (r *teamResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_team"
}

// Schema defines the schema for the resource.
func (r *teamResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a team.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "name of the team",
				Required:    true,
			},
			"members": schema.ListAttribute{
				ElementType: types.Int64Type,
				Description: "Members",
				Required:    true,
			},
			"slug": schema.StringAttribute{
				Description: "schema",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

// Create a new resource
func (r *teamResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan teamResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Error(ctx, fmt.Sprintf("plan %+v \n", plan))
	var memberIDs []int64
	for _, el := range plan.Members.Elements() {
		x := el.(types.Int64)
		memberIDs = append(memberIDs, x.ValueInt64())
	}
	tflog.Error(ctx, fmt.Sprintf("member ids %+v", memberIDs))
	teamMutation := gqlclient.TeamM{Team: gqlclient.Team{Name: plan.Name.ValueString()}, Members: memberIDs}

	team, err := r.c.CreateTeam(ctx, gqlclient.CreateTeamMutationInput{teamMutation})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating team",
			"Could not create team, unexpected error: "+err.Error(),
		)
		return
	}

	plan.Slug = types.StringValue(team.Slug)

	tflog.Error(ctx, fmt.Sprintf("Format in plan %+v\n %+v \n", team, plan))

	val := []attr.Value{}
	for _, m := range team.Members {
		v, diags := types.Int64Value(m).ToInt64Value(ctx)
		resp.Diagnostics.Append(diags...)

		val = append(val, v)
	}
	members, diags := types.ListValue(types.Int64Type, val)
	resp.Diagnostics.Append(diags...)
	trm := teamResourceModel{
		Name:    types.StringValue(team.Name),
		Slug:    types.StringValue(team.Slug),
		Members: members,
	}
	diags = resp.State.Set(ctx, trm)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information
func (r *teamResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	ctx = tflog.SetField(ctx, "resource", "team")

	var state teamResourceModel
	diags := req.State.Get(ctx, &state)

	tflog.Info(ctx, "Refreshing Team resource", map[string]any{"state": fmt.Sprintf("%+v", state)})
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	team, err := r.c.GetTeam(ctx, state.Slug.ValueString())
	if err != nil {
		tflog.Error(ctx, fmt.Sprintf("Error obtaining team: %+v", err))
		resp.Diagnostics.AddError(
			"Error Reading Team",
			"Could not read Team Slug "+state.Slug.ValueString()+": "+err.Error(),
		)
		return
	}

	val := []attr.Value{}
	for _, m := range team.Members {
		v, diags := types.Int64Value(m).ToInt64Value(ctx)
		resp.Diagnostics.Append(diags...)

		val = append(val, v)
	}
	members, _ := types.ListValue(types.Int64Type, val)

	trm := teamResourceModel{
		Name:    types.StringValue(team.Name),
		Slug:    types.StringValue(team.Slug),
		Members: members,
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &trm)
	resp.Diagnostics.Append(diags...)
	return
}

func (r *teamResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	ctx = tflog.SetField(ctx, "resource", "team")

	// Retrieve values from plan
	var plan teamResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	var state teamResourceModel
	diags = req.Plan.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	tflog.Info(ctx, "Updating team", map[string]any{"slug": state.Slug})

	if resp.Diagnostics.HasError() {
		return
	}

	memberIDs := transformMemberIDAttrValueToInt(plan.Members.Elements())

	teamMutation := gqlclient.TeamM{
		Team: gqlclient.Team{
			Name: plan.Name.ValueString(),
			Slug: state.Slug.ValueString(),
		},
		Members: memberIDs,
	}

	team, err := r.c.UpdateTeam(ctx, gqlclient.UpdateTeamMutationInput{TeamM: teamMutation})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating team",
			"Could not create team, unexpected error: "+err.Error(),
		)
		return
	}

	plan.Slug = types.StringValue(team.Slug)

	members, diags := transformMemberIDIntToAttrValue(team.Members)

	trm := teamResourceModel{
		Name:    types.StringValue(team.Name),
		Slug:    types.StringValue(team.Slug),
		Members: members,
	}
	diags = resp.State.Set(ctx, trm)
	resp.Diagnostics.Append(diags...)
	return
}

func (r *teamResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	ctx = tflog.SetField(ctx, "resource", "team")

	var plan teamResourceModel
	if resp.Diagnostics.HasError() {
		return
	}

	req.State.Get(ctx, &plan)
	tflog.Info(ctx, "Deleting team", map[string]any{"slug": plan.Slug})

	teamMutation := gqlclient.DeleteTeamMutationInput{
		Slug: plan.Slug.ValueString(),
	}

	err := r.c.DeleteTeam(ctx, teamMutation)
	if err != nil {
		tflog.Info(ctx, "Error deleting plan", map[string]any{"err": err})
		resp.Diagnostics.AddError(
			"Error deleting team",
			"Could not create team, unexpected error: "+err.Error(),
		)
		return
	}
}

// Configure adds the provider configured client to the resource.
func (r *teamResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.c = req.ProviderData.(*gqlclient.Client)
}

func transformMemberIDAttrValueToInt(members []attr.Value) []int64 {
	var memberIDs []int64
	for _, el := range members {
		x := el.(types.Int64)
		memberIDs = append(memberIDs, x.ValueInt64())
	}
	return memberIDs
}

func transformMemberIDIntToAttrValue(members []int64) (basetypes.ListValue, diag.Diagnostics) {
	val := []attr.Value{}
	for _, m := range members {
		v, _ := types.Int64Value(m).ToInt64Value(nil) // Context is not used
		val = append(val, v)
	}
	return types.ListValue(types.Int64Type, val)
}

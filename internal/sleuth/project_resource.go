package sleuth

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/sleuth-io/terraform-provider-sleuth/internal/gqlclient"
)

var (
	_ resource.Resource                = &projectResource{}
	_ resource.ResourceWithConfigure   = &projectResource{}
	_ resource.ResourceWithImportState = &projectResource{}
)

type projectResourceModel struct {
	ID                        types.String `tfsdk:"id"`
	Name                      types.String `tfsdk:"name"`
	Slug                      types.String `tfsdk:"slug"`
	Description               types.String `tfsdk:"description"`
	IssueTrackerProviderType  types.String `tfsdk:"issue_tracker_provider_type"`
	BuildProvider             types.String `tfsdk:"build_provider"`
	ChangeFailureRateBoundary types.String `tfsdk:"change_failure_rate_boundary"`
	ImpactSensitivity         types.String `tfsdk:"impact_sensitivity"`
	FailureSensitivity        types.Int64  `tfsdk:"failure_sensitivity"`
}

type projectResource struct {
	c *gqlclient.Client
}

func NewProjectResource() resource.Resource {
	return &projectResource{}
}

func (p *projectResource) Schema(_ context.Context, _ resource.SchemaRequest, res *resource.SchemaResponse) {
	res.Schema = schema.Schema{
		Description: "Project resource manages Sleuth project.",
		Attributes: map[string]schema.Attribute{
			// Added due to tests
			"id": schema.StringAttribute{
				Computed: true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Project name",
				Required:            true,
			},
			"slug": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Project slug",
			},
			"description": schema.StringAttribute{
				DeprecationMessage:  "Project description will be removed",
				MarkdownDescription: "Project description",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"issue_tracker_provider_type": schema.StringAttribute{
				MarkdownDescription: "Where to find issues linked to by changes",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("SOURCE_PROVIDER"),
			},
			"build_provider": schema.StringAttribute{
				MarkdownDescription: "Where to find builds related to changes",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("NONE"),
			},
			"change_failure_rate_boundary": schema.StringAttribute{
				MarkdownDescription: "The health rating at which point it will be considered a failure",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("UNHEALTHY"),
			},
			"impact_sensitivity": schema.StringAttribute{
				MarkdownDescription: "How many impact measures Sleuth takes into account when auto-determining a deploys health.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("NORMAL"),
			},
			"failure_sensitivity": schema.Int64Attribute{
				MarkdownDescription: "The amount of time (in seconds) a deploy must spend in a failure status (Unhealthy, Incident, etc.) before it is determined a failure. Setting this value to a longer time means that less deploys will be classified.",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(420),
			},
		},
	}
}

func (p *projectResource) Configure(ctx context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	p.c = req.ProviderData.(*gqlclient.Client)
}

func (p *projectResource) Metadata(ctx context.Context, req resource.MetadataRequest, res *resource.MetadataResponse) {
	res.TypeName = req.ProviderTypeName + "_project"
}

func (p *projectResource) Create(ctx context.Context, req resource.CreateRequest, res *resource.CreateResponse) {
	ctx = tflog.SetField(ctx, "resource", "project")
	ctx = tflog.SetField(ctx, "operation", "create")

	var plan projectResourceModel
	diags := req.Plan.Get(ctx, &plan)
	res.Diagnostics.Append(diags...)

	tflog.Info(ctx, "Creating Project resource", map[string]any{"plan": plan})

	inputFields := getMutableProjectStruct(plan)

	input := gqlclient.CreateProjectMutationInput{MutableProject: &inputFields}

	proj, err := p.c.CreateProject(input)
	if err != nil {
		tflog.Error(ctx, "Error creating Project", map[string]any{"error": err.Error()})
		res.Diagnostics.AddError(
			"Error creating Project",
			fmt.Sprintf("Could not create project, unexpected error: %+v", err.Error()),
		)
		return
	}

	tflog.Info(ctx, "Created Project", map[string]any{"project": proj})

	state := getNewStateFromProject(proj)

	diags = res.State.Set(ctx, state)
	res.Diagnostics.Append(diags...)
}

func (p *projectResource) Read(ctx context.Context, req resource.ReadRequest, res *resource.ReadResponse) {
	ctx = tflog.SetField(ctx, "resource", "project")
	ctx = tflog.SetField(ctx, "operation", "read")

	var state projectResourceModel
	diags := req.State.Get(ctx, &state)

	tflog.Info(ctx, "Refreshing Project resource", map[string]any{"state": fmt.Sprintf("%+v", state)})
	res.Diagnostics.Append(diags...)

	if res.Diagnostics.HasError() {
		return
	}

	proj, err := p.c.GetProject(state.Slug.ValueStringPointer())
	if err != nil {
		tflog.Error(ctx, fmt.Sprintf("Error obtaining project: %+v", err))
		res.Diagnostics.AddError(
			"Error Reading Project",
			fmt.Sprintf("Could not read Project Slug %+s, %+v", state.Slug.ValueString(), err.Error()),
		)
		return
	}
	if proj == nil {
		return
	}

	newState := getNewStateFromProject(proj)

	diags = res.State.Set(ctx, &newState)
	res.Diagnostics.Append(diags...)
	return
}

func (p *projectResource) Update(ctx context.Context, req resource.UpdateRequest, res *resource.UpdateResponse) {
	ctx = tflog.SetField(ctx, "resource", "project")
	ctx = tflog.SetField(ctx, "operation", "update")

	var plan projectResourceModel
	diags := req.Plan.Get(ctx, &plan)
	res.Diagnostics.Append(diags...)

	var state projectResourceModel
	diags = req.State.Get(ctx, &state)
	res.Diagnostics.Append(diags...)

	tflog.Info(ctx, "Updating Project resource", map[string]any{"plan": plan, "state": state})

	inputFields := getMutableProjectStruct(plan)

	input := gqlclient.UpdateProjectMutationInput{Slug: state.Slug.ValueString(), MutableProject: &inputFields}

	proj, err := p.c.UpdateProject(state.Slug.ValueStringPointer(), input)
	tflog.Error(ctx, fmt.Sprintf("PRoj: %+v %+v", proj, err))
	if err != nil {
		diags.AddError("Error updating project", err.Error())
		return
	}

	tflog.Info(ctx, "Updated Project", map[string]any{"project": proj})

	newState := getNewStateFromProject(proj)

	diags = res.State.Set(ctx, newState)
	res.Diagnostics.Append(diags...)
}

func (p *projectResource) Delete(ctx context.Context, req resource.DeleteRequest, res *resource.DeleteResponse) {
	tflog.SetField(ctx, "resource", "project")
	tflog.SetField(ctx, "operation", "delete")

	var state projectResourceModel
	req.State.Get(ctx, &state)

	tflog.Info(ctx, "Deleting Project resource", map[string]any{"state": fmt.Sprintf("%+v", state)})

	err := p.c.DeleteProject(state.Slug.ValueStringPointer())
	if err != nil {
		// Ignore missing as the project gets deleted when the last env gets deleted
		if err.Error() != "Missing" {
			tflog.Error(ctx, "Unexpected error deleting project", map[string]any{"error": err.Error()})
			res.Diagnostics.AddError("Unexpected error deleting project", err.Error())
			return
		}
	}

	res.State.RemoveResource(ctx)
}

func (p *projectResource) ImportState(ctx context.Context, req resource.ImportStateRequest, res *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("slug"), req, res)
}

func getNewStateFromProject(proj *gqlclient.Project) projectResourceModel {
	return projectResourceModel{
		ID:                        types.StringValue(proj.Slug),
		Name:                      types.StringValue(proj.Name),
		Slug:                      types.StringValue(proj.Slug),
		Description:               types.StringValue(proj.Description),
		IssueTrackerProviderType:  types.StringValue(proj.IssueTrackerProvider),
		BuildProvider:             types.StringValue(proj.BuildProvider),
		ChangeFailureRateBoundary: types.StringValue(proj.ChangeFailureRateBoundary),
		ImpactSensitivity:         types.StringValue(proj.ImpactSensitivity),
		FailureSensitivity:        types.Int64Value(int64(proj.FailureSensitivity)),
	}
}

func getMutableProjectStruct(plan projectResourceModel) gqlclient.MutableProject {
	return gqlclient.MutableProject{
		Name:                      plan.Name.ValueString(),
		Description:               plan.Description.ValueString(),
		IssueTrackerProvider:      plan.IssueTrackerProviderType.ValueString(),
		BuildProvider:             plan.BuildProvider.ValueString(),
		ChangeFailureRateBoundary: plan.ChangeFailureRateBoundary.ValueString(),
		ImpactSensitivity:         plan.ImpactSensitivity.ValueString(),
		FailureSensitivity:        int(plan.FailureSensitivity.ValueInt64()),
	}
}

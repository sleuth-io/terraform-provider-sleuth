package sleuth

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/sleuth-io/terraform-provider-sleuth/internal/gqlclient"
)

var (
	_ resource.Resource                = &errorImpactSourceResource{}
	_ resource.ResourceWithConfigure   = &errorImpactSourceResource{}
	_ resource.ResourceWithImportState = &errorImpactSourceResource{}
)

type errorImpactResourceModel struct {
	ID   types.String `tfsdk:"id"`
	Slug types.String `tfsdk:"slug"`

	ProjectSlug     types.String `tfsdk:"project_slug"`
	EnvironmentSlug types.String `tfsdk:"environment_slug"`

	Name                       types.String  `tfsdk:"name"`
	ProviderType               types.String  `tfsdk:"provider_type"`
	ErrorOrgKey                types.String  `tfsdk:"error_org_key"`
	ErrorProjectKey            types.String  `tfsdk:"error_project_key"`
	ErrorEnvironment           types.String  `tfsdk:"error_environment"`
	ManuallySetHealthThreshold types.Float64 `tfsdk:"manually_set_health_threshold"`
	IntegrationSlug            types.String  `tfsdk:"integration_slug"`
}

type errorImpactSourceResource struct {
	c *gqlclient.Client
}

func NewErrorImpactSourceResource() resource.Resource {
	return &errorImpactSourceResource{}
}

func (eisr *errorImpactSourceResource) Schema(_ context.Context, _ resource.SchemaRequest, res *resource.SchemaResponse) {
	res.Schema = schema.Schema{
		MarkdownDescription: "Sleuth error impact source.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"slug": schema.StringAttribute{
				Computed: true,
			},
			"project_slug": schema.StringAttribute{
				MarkdownDescription: "The slug of the project that this error impact source belongs to.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(), // ForceNew replacement
				},
			},
			"environment_slug": schema.StringAttribute{
				MarkdownDescription: "The slug of the environment that this error impact source belongs to.",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Error impact source name",
				Required:            true,
			},
			"provider_type": schema.StringAttribute{
				MarkdownDescription: "Integration provider type",
				Required:            true,
			},
			"error_org_key": schema.StringAttribute{
				MarkdownDescription: "The organization key of the integration provider",
				Required:            true,
			},
			"error_project_key": schema.StringAttribute{
				MarkdownDescription: "The project key of the integration provider",
				Required:            true,
			},
			"error_environment": schema.StringAttribute{
				MarkdownDescription: "The environment of the integration provider",
				Required:            true,
			},
			"manually_set_health_threshold": schema.Float64Attribute{
				MarkdownDescription: "The manually set threshold to start marking failed values",
				Optional:            true,
			},
			"integration_slug": schema.StringAttribute{
				MarkdownDescription: "The integration slug",
				Optional:            true,
				Computed:            true,
			},
		},
	}
}

func (eisr *errorImpactSourceResource) Configure(ctx context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	eisr.c = req.ProviderData.(*gqlclient.Client)
}

func (eisr *errorImpactSourceResource) Metadata(ctx context.Context, req resource.MetadataRequest, res *resource.MetadataResponse) {
	res.TypeName = req.ProviderTypeName + "_error_impact_source"
}

func (eisr *errorImpactSourceResource) Create(ctx context.Context, req resource.CreateRequest, res *resource.CreateResponse) {
	ctx = tflog.SetField(ctx, "resource", "error_impact_source")
	ctx = tflog.SetField(ctx, "operation", "create")

	var plan errorImpactResourceModel
	diags := req.Plan.Get(ctx, &plan)
	res.Diagnostics.Append(diags...)

	tflog.Info(ctx, "Creating ErrorImpactSource resource", map[string]any{"name": plan.Name.ValueString(), "projectSlug": plan.ProjectSlug.ValueString()})

	if res.Diagnostics.HasError() {
		tflog.Error(ctx, "Error getting ErrorImpactSource plan", map[string]any{"diagnostics": res.Diagnostics})
		return
	}

	projectSlug := plan.ProjectSlug.ValueString()
	inputFields := getMutableErrorImpactSourceStruct(plan)

	input := gqlclient.CreateErrorImpactSourceMutationInput{
		ProjectSlug:              projectSlug,
		MutableErrorImpactSource: inputFields,
	}

	eis, err := eisr.c.CreateErrorImpactSource(input)
	if err != nil {
		tflog.Error(ctx, "Error creating ErrorImpactSource", map[string]any{"error": err.Error()})
		res.Diagnostics.AddError(
			"Error creating ErrorImpactSource",
			fmt.Sprintf("Could not create error impact source, unexpected error: %+v", err.Error()),
		)
		return
	}

	state := getNewStateFromErrorImpactSource(eis, projectSlug)
	res.Diagnostics.Append(diags...)
	diags = res.State.Set(ctx, state)
	res.Diagnostics.Append(diags...)
	tflog.Info(ctx, fmt.Sprintf("Successfully created ErrorImpactSource"), map[string]any{"diags": res.Diagnostics})
}

func (eisr *errorImpactSourceResource) Read(ctx context.Context, req resource.ReadRequest, res *resource.ReadResponse) {
	ctx = tflog.SetField(ctx, "resource", "error_impact_source")
	ctx = tflog.SetField(ctx, "operation", "read")

	var state errorImpactResourceModel
	diags := req.State.Get(ctx, &state)
	res.Diagnostics.Append(diags...)

	tflog.Info(ctx, "Reading ErrorImpactSource resource", map[string]any{"state": state})
	projectSlug := state.ProjectSlug.ValueString()
	slug := state.Slug.ValueString()

	eis, err := eisr.c.GetErrorImpactSource(&projectSlug, &slug)
	if err != nil {
		tflog.Error(ctx, "Error reading ErrorImpactSource", map[string]any{"error": err.Error()})
		res.Diagnostics.AddError(
			"Error reading ErrorImpactSource",
			fmt.Sprintf("Could not read error impact source, unexpected error: %+v", err.Error()),
		)
		return
	}
	newState := getNewStateFromErrorImpactSource(eis, projectSlug)
	res.Diagnostics.Append(diags...)

	diags = res.State.Set(ctx, newState)
	res.Diagnostics.Append(diags...)
}

func (eisr *errorImpactSourceResource) Update(ctx context.Context, req resource.UpdateRequest, res *resource.UpdateResponse) {
	ctx = tflog.SetField(ctx, "resource", "error_impact_source")
	ctx = tflog.SetField(ctx, "operation", "update")

	var state errorImpactResourceModel
	diags := req.State.Get(ctx, &state)
	res.Diagnostics.Append(diags...)

	var plan errorImpactResourceModel
	diags = req.Plan.Get(ctx, &plan)
	res.Diagnostics.Append(diags...)

	tflog.Info(ctx, "Updating ErrorImpactSource resource", map[string]any{"plan": plan})

	if res.Diagnostics.HasError() {
		tflog.Error(ctx, "Error getting ErrorImpactSource plan", map[string]any{"diagnostics": res.Diagnostics})
		return
	}

	projectSlug := plan.ProjectSlug.ValueString()
	inputFields := getMutableErrorImpactSourceStruct(plan)

	input := gqlclient.UpdateErrorImpactSourceMutationInput{
		ProjectSlug:              projectSlug,
		Slug:                     state.Slug.ValueString(),
		MutableErrorImpactSource: inputFields,
	}

	eis, err := eisr.c.UpdateErrorImpactSource(input)
	if err != nil {
		tflog.Error(ctx, "Error updating ErrorImpactSource", map[string]any{"error": err.Error()})
		res.Diagnostics.AddError(
			"Error updating ErrorImpactSource",
			fmt.Sprintf("Could not update error impact source, unexpected error: %+v", err.Error()),
		)
		return
	}

	newState := getNewStateFromErrorImpactSource(eis, projectSlug)
	res.Diagnostics.Append(diags...)

	diags = res.State.Set(ctx, newState)
	res.Diagnostics.Append(diags...)
	tflog.Info(ctx, fmt.Sprintf("Successfully created ErrorImpactSource"), map[string]any{"diags": res.Diagnostics})

}

func (eisr *errorImpactSourceResource) Delete(ctx context.Context, req resource.DeleteRequest, res *resource.DeleteResponse) {
	ctx = tflog.SetField(ctx, "resource", "error_impact_source")
	ctx = tflog.SetField(ctx, "operation", "delete")

	var state errorImpactResourceModel
	diags := req.State.Get(ctx, &state)
	res.Diagnostics.Append(diags...)

	tflog.Info(ctx, "Deleting ErrorImpactSource resource", map[string]any{"state": state})
	projectSlug := state.ProjectSlug.ValueStringPointer()
	slug := state.Slug.ValueStringPointer()

	err := eisr.c.DeleteImpactSource(projectSlug, slug)
	if err != nil {
		tflog.Error(ctx, "Error deleting ErrorImpactSource", map[string]any{"error": err.Error()})
		res.Diagnostics.AddError(
			"Error deleting ErrorImpactSource",
			fmt.Sprintf("Could not delete error impact source, unexpected error: %+v", err.Error()),
		)
		return
	}
}

func (eisr *errorImpactSourceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, res *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("slug"), req, res)
}

func getNewStateFromErrorImpactSource(eis *gqlclient.ErrorImpactSource, projectSlug string) errorImpactResourceModel {
	return errorImpactResourceModel{
		ID:                         types.StringValue(eis.Slug),
		Slug:                       types.StringValue(eis.Slug),
		ProjectSlug:                types.StringValue(projectSlug),
		EnvironmentSlug:            types.StringValue(eis.Environment.Slug),
		Name:                       types.StringValue(eis.Name),
		ProviderType:               types.StringValue(strings.ToUpper(eis.Provider)),
		ErrorOrgKey:                types.StringValue(eis.ErrorOrgKey),
		ErrorProjectKey:            types.StringValue(eis.ErrorProjectKey),
		ErrorEnvironment:           types.StringValue(eis.ErrorEnvironment),
		ManuallySetHealthThreshold: types.Float64PointerValue(eis.ManuallySetHealthThreshold),
		IntegrationSlug:            types.StringValue(eis.IntegrationAuthSlug),
	}
}

func getMutableErrorImpactSourceStruct(plan errorImpactResourceModel) *gqlclient.MutableErrorImpactSource {

	return &gqlclient.MutableErrorImpactSource{
		EnvironmentSlug:            plan.EnvironmentSlug.ValueString(),
		Name:                       plan.Name.ValueString(),
		Provider:                   plan.ProviderType.ValueString(),
		ErrorOrgKey:                plan.ErrorOrgKey.ValueString(),
		ErrorProjectKey:            plan.ErrorProjectKey.ValueString(),
		ErrorEnvironment:           plan.ErrorEnvironment.ValueString(),
		ManuallySetHealthThreshold: plan.ManuallySetHealthThreshold.ValueFloat64(),
		IntegrationSlug:            plan.IntegrationSlug.ValueString(),
	}
}

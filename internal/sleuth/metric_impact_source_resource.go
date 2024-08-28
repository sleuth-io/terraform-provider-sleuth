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
	_ resource.Resource                = &metricImpactSourceResource{}
	_ resource.ResourceWithConfigure   = &metricImpactSourceResource{}
	_ resource.ResourceWithImportState = &metricImpactSourceResource{}
)

type metricImpactResourceModel struct {
	ID   types.String `tfsdk:"id"`
	Slug types.String `tfsdk:"slug"`

	ProjectSlug types.String `tfsdk:"project_slug"`
	EnvSlug     types.String `tfsdk:"environment_slug"`

	Name                       types.String  `tfsdk:"name"`
	ProviderType               types.String  `tfsdk:"provider_type"`
	IntegrationSlug            types.String  `tfsdk:"integration_slug"`
	Query                      types.String  `tfsdk:"query"`
	LessIsBetter               types.Bool    `tfsdk:"less_is_better"`
	ManuallySetHealthThreshold types.Float64 `tfsdk:"manually_set_health_threshold"`
}

type metricImpactSourceResource struct {
	c *gqlclient.Client
}

func NewMetricImpactSourceResource() resource.Resource {
	return &metricImpactSourceResource{}
}

func (misr *metricImpactSourceResource) Schema(_ context.Context, _ resource.SchemaRequest, res *resource.SchemaResponse) {
	res.Schema = schema.Schema{
		MarkdownDescription: "Sleuth metric impact source.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"slug": schema.StringAttribute{
				Computed: true,
			},
			"project_slug": schema.StringAttribute{
				MarkdownDescription: "The slug of the project that this metric impact source belongs to.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(), // ForceNew replacement
				},
			},
			"environment_slug": schema.StringAttribute{
				MarkdownDescription: "The slug of the environment that this metric impact source belongs to.",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Impact source name",
				Required:            true,
			},

			"provider_type": schema.StringAttribute{
				MarkdownDescription: "Integration provider type",
				Required:            true,
			},
			"integration_slug": schema.StringAttribute{
				MarkdownDescription: "Integration slug is generated automatically when an integration is set up in Sleuth. By default, it matches the `provider_type`. Any value specified in the integration's `Description label` field gets appended to the `integration_slug`, spaces replaced with dashes, e.g. `cloudwatch-test`",
				Optional:            true,
				Computed:            true,
			},
			"query": schema.StringAttribute{
				MarkdownDescription: "The metric query",
				Required:            true,
			},
			"less_is_better": schema.BoolAttribute{
				MarkdownDescription: "Whether smaller values are better or not",
				Optional:            true,
				Computed:            true,
			},
			"manually_set_health_threshold": schema.Float64Attribute{
				MarkdownDescription: "The manually set threshold to start marking failed values",
				Optional:            true,
			},
		},
	}
}

func (misr *metricImpactSourceResource) Configure(ctx context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	misr.c = req.ProviderData.(*gqlclient.Client)
}

func (misr *metricImpactSourceResource) Metadata(ctx context.Context, req resource.MetadataRequest, res *resource.MetadataResponse) {
	res.TypeName = req.ProviderTypeName + "_metric_impact_source"
}

func (misr *metricImpactSourceResource) Create(ctx context.Context, req resource.CreateRequest, res *resource.CreateResponse) {
	ctx = tflog.SetField(ctx, "resource", "metric_impact_source")
	ctx = tflog.SetField(ctx, "operation", "create")

	var plan metricImpactResourceModel
	diags := req.Plan.Get(ctx, &plan)
	res.Diagnostics.Append(diags...)

	tflog.Info(ctx, "Creating MetricImpactSource resource", map[string]any{"name": plan.Name.ValueString(), "projectSlug": plan.ProjectSlug.ValueString()})

	if res.Diagnostics.HasError() {
		tflog.Error(ctx, "Error getting MetricImpactSource plan", map[string]any{"diagnostics": res.Diagnostics})
		return
	}

	projectSlug := plan.ProjectSlug.ValueString()
	inputFields := getMutableMetricImpactSourceStruct(plan)

	input := gqlclient.CreateMetricImpactSourceMutationInput{
		ProjectSlug:               projectSlug,
		MutableMetricImpactSource: inputFields,
	}

	mis, err := misr.c.CreateMetricImpactSource(input)
	if err != nil {
		tflog.Error(ctx, "Error creating MetricImpactSource", map[string]any{"error": err.Error()})
		res.Diagnostics.AddError(
			"Error creating MetricImpactSource",
			fmt.Sprintf("Could not create metric impact source, unexpected error: %+v", err.Error()),
		)
		return
	}

	state := getNewStateFromMetricImpactSource(mis, projectSlug)
	res.Diagnostics.Append(diags...)
	diags = res.State.Set(ctx, state)
	res.Diagnostics.Append(diags...)
	tflog.Info(ctx, "Successfully created MetricImpactSource", map[string]any{"diags": res.Diagnostics})
}

func (misr *metricImpactSourceResource) Read(ctx context.Context, req resource.ReadRequest, res *resource.ReadResponse) {
	ctx = tflog.SetField(ctx, "resource", "metric_impact_source")
	ctx = tflog.SetField(ctx, "operation", "read")

	var state metricImpactResourceModel
	diags := req.State.Get(ctx, &state)
	res.Diagnostics.Append(diags...)

	tflog.Info(ctx, "Reading MetricImpactSource resource", map[string]any{"state": state})
	projectSlug := state.ProjectSlug.ValueString()
	slug := state.Slug.ValueString()

	// when importing a resource, only ID will be set with project slug & slug
	if projectSlug == "" {
		id := state.ID.ValueString()
		splits := strings.Split(id, "/")
		projectSlug = splits[0]
		slug = splits[1]
	}

	ccs, err := misr.c.GetMetricImpactSource(&projectSlug, &slug)
	if err != nil {
		tflog.Error(ctx, "Error reading MetricImpactSource", map[string]any{"error": err.Error()})
		res.Diagnostics.AddError(
			"Error reading MetricImpactSource",
			fmt.Sprintf("Could not read metric impact source, unexpected error: %+v", err.Error()),
		)
		return
	}
	newState := getNewStateFromMetricImpactSource(ccs, projectSlug)
	res.Diagnostics.Append(diags...)

	diags = res.State.Set(ctx, newState)
	res.Diagnostics.Append(diags...)

}

func (misr *metricImpactSourceResource) Update(ctx context.Context, req resource.UpdateRequest, res *resource.UpdateResponse) {
	ctx = tflog.SetField(ctx, "resource", "metric_impact_source")
	ctx = tflog.SetField(ctx, "operation", "update")

	var state metricImpactResourceModel
	diags := req.State.Get(ctx, &state)
	res.Diagnostics.Append(diags...)

	var plan metricImpactResourceModel
	diags = req.Plan.Get(ctx, &plan)
	res.Diagnostics.Append(diags...)

	tflog.Info(ctx, "Updating MetricImpactSource resource", map[string]any{"plan": plan})

	projectSlug := plan.ProjectSlug.ValueString()
	inputFields := getMutableMetricImpactSourceStruct(plan)

	input := gqlclient.UpdateMetricImpactSourceMutationInput{
		ProjectSlug:               projectSlug,
		Slug:                      state.Slug.ValueString(),
		MutableMetricImpactSource: inputFields,
	}

	ccs, err := misr.c.UpdateMetricImpactSource(input)
	if err != nil {
		tflog.Error(ctx, "Error updating MetricImpactSource", map[string]any{"error": err.Error()})
		res.Diagnostics.AddError(
			"Error updating MetricImpactSource",
			fmt.Sprintf("Could not update metric impact source, unexpected error: %+v", err.Error()),
		)
		return
	}

	newState := getNewStateFromMetricImpactSource(ccs, projectSlug)

	diags = res.State.Set(ctx, newState)
	res.Diagnostics.Append(diags...)
	tflog.Info(ctx, "Successfully created MetricImpactSource", map[string]any{"diags": res.Diagnostics})

}

func (misr *metricImpactSourceResource) Delete(ctx context.Context, req resource.DeleteRequest, res *resource.DeleteResponse) {
	ctx = tflog.SetField(ctx, "resource", "metric_impact_source")
	ctx = tflog.SetField(ctx, "operation", "delete")

	var state metricImpactResourceModel
	diags := req.State.Get(ctx, &state)
	res.Diagnostics.Append(diags...)

	tflog.Info(ctx, "Deleting MetricImpactSource resource", map[string]any{"state": state})
	projectSlug := state.ProjectSlug.ValueStringPointer()
	slug := state.Slug.ValueStringPointer()

	err := misr.c.DeleteImpactSource(projectSlug, slug)
	if err != nil {
		tflog.Error(ctx, "Error deleting MetricImpactSource", map[string]any{"error": err.Error()})
		res.Diagnostics.AddError(
			"Error deleting MetricImpactSource",
			fmt.Sprintf("Could not delete metric impact source, unexpected error: %+v", err.Error()),
		)
		return
	}
}

func (misr *metricImpactSourceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, res *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, res)
}

func getNewStateFromMetricImpactSource(ccs *gqlclient.MetricImpactSource, projectSlug string) metricImpactResourceModel {
	return metricImpactResourceModel{
		ID:                         types.StringValue(ccs.Slug),
		Slug:                       types.StringValue(ccs.Slug),
		ProjectSlug:                types.StringValue(projectSlug),
		EnvSlug:                    types.StringValue(ccs.Environment.Slug),
		Name:                       types.StringValue(ccs.Name),
		ProviderType:               types.StringValue(strings.ToUpper(ccs.Provider)),
		IntegrationSlug:            types.StringValue(ccs.IntegrationAuthSlug),
		Query:                      types.StringValue(ccs.Query),
		LessIsBetter:               types.BoolValue(ccs.LessIsBetter),
		ManuallySetHealthThreshold: types.Float64PointerValue(ccs.ManuallySetHealthThreshold),
	}
}

func getMutableMetricImpactSourceStruct(plan metricImpactResourceModel) *gqlclient.MutableMetricImpactSource {
	return &gqlclient.MutableMetricImpactSource{
		EnvironmentSlug:            plan.EnvSlug.ValueString(),
		Name:                       plan.Name.ValueString(),
		Provider:                   strings.ToUpper(plan.ProviderType.ValueString()),
		Query:                      plan.Query.ValueString(),
		IntegrationSlug:            plan.IntegrationSlug.ValueString(),
		LessIsBetter:               plan.LessIsBetter.ValueBool(),
		ManuallySetHealthThreshold: plan.ManuallySetHealthThreshold.ValueFloat64(),
	}
}

package sleuth

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/sleuth-io/terraform-provider-sleuth/internal/gqlclient"
)

var (
	_ resource.Resource                = &incidentImpactSourceResource{}
	_ resource.ResourceWithConfigure   = &incidentImpactSourceResource{}
	_ resource.ResourceWithImportState = &incidentImpactSourceResource{}
)

type incidentImpactResourceModel struct {
	ID   types.String `tfsdk:"id"`
	Slug types.String `tfsdk:"slug"`

	ProjectSlug     types.String `tfsdk:"project_slug"`
	EnvironmentName types.String `tfsdk:"environment_name"`

	Name         types.String `tfsdk:"name"`
	ProviderName types.String `tfsdk:"provider_name"`

	PagerDutyInput   types.Object `tfsdk:"pagerduty_input"`
	DataDogInput     types.Object `tfsdk:"datadog_input"`
	JiraInput        types.Object `tfsdk:"jira_input"`
	BlamelessInput   types.Object `tfsdk:"blameless_input"`
	StatusPageInput  types.Object `tfsdk:"statuspage_input"`
	OpsGenieInput    types.Object `tfsdk:"opsgenie_input"`
	FireHydrantInput types.Object `tfsdk:"firehydrant_input"`
	ClubhouseInput   types.Object `tfsdk:"clubhouse_input"`
}

type pagerDutyInputResourceModel struct {
	RemoteServices  types.String `tfsdk:"remote_services"`
	RemoteUrgency   types.String `tfsdk:"remote_urgency"`
	IntegrationSlug types.String `tfsdk:"integration_slug"`
}

type dataDogInputResourceModel struct {
	Query                   types.String `tfsdk:"query"`
	RemotePriorityThreshold types.String `tfsdk:"remote_priority_threshold"`
	IntegrationSlug         types.String `tfsdk:"integration_slug"`
}

type jiraInputResourceModel struct {
	RemoteJQL       types.String `tfsdk:"remote_jql"`
	IntegrationSlug types.String `tfsdk:"integration_slug"`
}

type blamelessInputResourceModel struct {
	RemoteTypes             types.Set    `tfsdk:"remote_types"`
	RemoteSeverityThreshold types.String `tfsdk:"remote_severity_threshold"`
	IntegrationSlug         types.String `tfsdk:"integration_slug"`
}

type statuspageInputResourceModel struct {
	RemotePage                 types.String `tfsdk:"remote_page"`
	RemoteComponent            types.String `tfsdk:"remote_component"`
	RemoteImpact               types.String `tfsdk:"remote_impact"`
	IgnoreMaintenanceIncidents types.Bool   `tfsdk:"ignore_maintenance_incidents"`
	IntegrationSlug            types.String `tfsdk:"integration_slug"`
}

type opsgenieInputResourceModel struct {
	RemoteAlertTags         types.String `tfsdk:"remote_alert_tags"`
	RemoteIncidentTags      types.String `tfsdk:"remote_incident_tags"`
	RemotePriorityThreshold types.String `tfsdk:"remote_priority_threshold"`
	RemoteService           types.String `tfsdk:"remote_service"`
	RemoteUseAlerts         types.Bool   `tfsdk:"remote_use_alerts"`
	IntegrationSlug         types.String `tfsdk:"integration_slug"`
}

type firehydrantInputResourceModel struct {
	RemoteEnvironments       types.String `tfsdk:"remote_environments"`
	RemoteServices           types.String `tfsdk:"remote_services"`
	RemoteMitigatedIsHealthy types.Bool   `tfsdk:"remote_mitigated_is_healthy"`
}

type clubhouseInputResourceModel struct {
	RemoteQuery     types.String `tfsdk:"remote_query"`
	IntegrationSlug types.String `tfsdk:"integration_slug"`
}

// Types for translating between TF blocks into resource models
// removegst when we solely support TF protocol v6
var pagerDutyInputResourceTFTypes = map[string]attr.Type{
	"remote_services":  types.StringType,
	"remote_urgency":   types.StringType,
	"integration_slug": types.StringType,
}

var dataDogInputResourceTFTypes = map[string]attr.Type{
	"query":                     types.StringType,
	"remote_priority_threshold": types.StringType,
	"integration_slug":          types.StringType,
}

var jiraInputResourceTFTypes = map[string]attr.Type{
	"remote_jql":       types.StringType,
	"integration_slug": types.StringType,
}

var blamelessInputResourceTFTypes = map[string]attr.Type{
	"remote_types":              types.SetType{ElemType: types.StringType},
	"remote_severity_threshold": types.StringType,
	"integration_slug":          types.StringType,
}

var statuspageInputResourceTFTypes = map[string]attr.Type{
	"remote_page":                  types.StringType,
	"remote_component":             types.StringType,
	"remote_impact":                types.StringType,
	"ignore_maintenance_incidents": types.BoolType,
	"integration_slug":             types.StringType,
}

var opsgenieInputResourceTFTypes = map[string]attr.Type{
	"remote_alert_tags":         types.StringType,
	"remote_incident_tags":      types.StringType,
	"remote_priority_threshold": types.StringType,
	"remote_service":            types.StringType,
	"remote_use_alerts":         types.BoolType,
	"integration_slug":          types.StringType,
}

var firehydrantInputResourceTFTypes = map[string]attr.Type{
	"remote_environments":         types.StringType,
	"remote_services":             types.StringType,
	"remote_mitigated_is_healthy": types.BoolType,
}

var clubhouseInputResourceTFTypes = map[string]attr.Type{
	"remote_query":     types.StringType,
	"integration_slug": types.StringType,
}

type incidentImpactSourceResource struct {
	c *gqlclient.Client
}

func NewIncidentImpactSourceResource() resource.Resource {
	return &incidentImpactSourceResource{}
}

func (iisr *incidentImpactSourceResource) Schema(_ context.Context, _ resource.SchemaRequest, res *resource.SchemaResponse) {
	res.Schema = schema.Schema{
		MarkdownDescription: "Sleuth code change source.",
		Blocks: map[string]schema.Block{
			"pagerduty_input": schema.SingleNestedBlock{
				MarkdownDescription: "PagerDuty input",
				Attributes: map[string]schema.Attribute{
					"remote_services": schema.StringAttribute{
						MarkdownDescription: "List of remote services, empty string means all",
						Optional:            true,
						Computed:            true,
						Default:             stringdefault.StaticString(""),
					},
					"remote_urgency": schema.StringAttribute{
						MarkdownDescription: "PagerDuty remote urgency, options: HIGH, LOW, ANY",
						Optional:            true,
						Computed:            true,
					},
					"integration_slug": schema.StringAttribute{
						MarkdownDescription: "IntegrationAuthentication slug used",
						Optional:            true,
					},
				},
			},
			"datadog_input": schema.SingleNestedBlock{
				MarkdownDescription: "DataDog input",
				Attributes: map[string]schema.Attribute{
					"query": schema.StringAttribute{
						MarkdownDescription: `The query to scope the monitors to track. If you are using a custom facet you would need to add @ to the beginning of the facet name. If empty, all monitors in Datadog will be matched regardless of environment or service.
See [DataDog documentation](https://docs.datadoghq.com/monitors/manage/search/) for more information.`,
						Optional: true,
						Computed: true,
						Default:  stringdefault.StaticString(""),
					},
					"remote_priority_threshold": schema.StringAttribute{
						Optional: true,
						Computed: true,
						Default:  stringdefault.StaticString("ALL"),
						Description: `Monitor states with matching or higher priorities will be considered a failure in Sleuth.
Options: ALL, P1, P2, P3, P4, P5. Defaults to ALL`,
					},
					"integration_slug": schema.StringAttribute{
						Optional:    true,
						Description: "DataDog IntegrationAuthentication slug from app",
					},
				},
			},
			"jira_input": schema.SingleNestedBlock{
				MarkdownDescription: "JIRA input",
				Attributes: map[string]schema.Attribute{
					"remote_jql": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "JIRA active incidents issues JQL",
					},
					"integration_slug": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "JIRA IntegrationAuthentication slug from app",
					},
				},
			},
			"blameless_input": schema.SingleNestedBlock{
				MarkdownDescription: "Blameless input",
				Attributes: map[string]schema.Attribute{
					"remote_types": schema.SetAttribute{
						Optional:            true,
						ElementType:         basetypes.StringType{},
						MarkdownDescription: "The types of incidents to the monitors should track",
					},
					"remote_severity_threshold": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Incidents with matching or lower severities will be considered a failure in Sleuth",
					},
					"integration_slug": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Blameless IntegrationAuthentication slug from app",
					},
				},
			},
			"statuspage_input": schema.SingleNestedBlock{
				MarkdownDescription: "Statuspage input",
				Attributes: map[string]schema.Attribute{
					"remote_page": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Statuspage page the incident impact source should monitor",
					},
					"remote_component": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Statuspage component the incident impact source should monitor",
					},
					"remote_impact": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Incidents with matching or lower severities will be considered a failure in Sleuth",
					},
					"ignore_maintenance_incidents": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Option to ignore maintenance incidents",
					},
					"integration_slug": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Statuspage IntegrationAuthentication slug from app",
					},
				},
			},
			"opsgenie_input": schema.SingleNestedBlock{
				MarkdownDescription: "OpsGenie input",
				Attributes: map[string]schema.Attribute{
					"remote_alert_tags": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Optionally filter by alert tags",
					},
					"remote_incident_tags": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Optionally filter by incident tags",
					},
					"remote_priority_threshold": schema.StringAttribute{
						Optional:            true,
						Computed:            true,
						MarkdownDescription: "Monitor states with matching or higher priorities will be considered a failure in Sleuth",
						Default:             stringdefault.StaticString("ALL"),
					},
					"remote_service": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "Only taken into consideration when using OpsGenie Incidents. This value should be the Unique ID of the OpsGenie service.",
					},
					"remote_use_alerts": schema.BoolAttribute{
						Optional:            true,
						MarkdownDescription: "Use OpsGenie Alerts instead of Incidents",
						Default:             booldefault.StaticBool(false),
						Computed:            true,
					},
					"integration_slug": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "The slug for the integration",
					},
				},
			},
			"firehydrant_input": schema.SingleNestedBlock{
				MarkdownDescription: "FireHydrant input",
				Attributes: map[string]schema.Attribute{
					"remote_environments": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "The environment defined in FireHydrant to monitor",
					},
					"remote_services": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "The service defined in FireHydrant to monitor",
					},
					"remote_mitigated_is_healthy": schema.BoolAttribute{
						MarkdownDescription: "If true, incident considered to have ended once reaching mitigated Milestone or it is resolved",
						Default:             booldefault.StaticBool(false),
						Computed:            true,
					},
				},
			},
			"clubhouse_input": schema.SingleNestedBlock{
				MarkdownDescription: "Clubhouse input",
				Attributes: map[string]schema.Attribute{
					"remote_query": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: `Need help finding query expression? See the [documentation](https://help.shortcut.com/hc/en-us/articles/360000046646-Searching-in-Shortcut-Using-Search-Operators) for more information.`,
					},
					"integration_slug": schema.StringAttribute{
						Optional:            true,
						MarkdownDescription: "IntegrationAuthentication slug used",
					},
				},
			},
		},
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"slug": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"project_slug": schema.StringAttribute{
				MarkdownDescription: "The slug of the project that this incident impact source belongs to.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(), // ForceNew replacement
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Impact source name",
				Required:            true,
			},
			"provider_name": schema.StringAttribute{
				MarkdownDescription: "Impact source provider (options: PAGERDUTY)",
				Required:            true,
			},
			"environment_name": schema.StringAttribute{
				MarkdownDescription: "Impact source environment name",
				Required:            true,
			},
		},
	}
}

func (iisr *incidentImpactSourceResource) Configure(ctx context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	iisr.c = req.ProviderData.(*gqlclient.Client)
}

func (iisr *incidentImpactSourceResource) Metadata(ctx context.Context, req resource.MetadataRequest, res *resource.MetadataResponse) {
	res.TypeName = req.ProviderTypeName + "_incident_impact_source"
}

type providerData struct {
	pagerduty   *pagerDutyInputResourceModel
	datadog     *dataDogInputResourceModel
	jira        *jiraInputResourceModel
	blameless   *blamelessInputResourceModel
	statuspage  *statuspageInputResourceModel
	clubhouse   *clubhouseInputResourceModel
	firehydrant *firehydrantInputResourceModel
	opsgenie    *opsgenieInputResourceModel
}

// we have to manually parse the provider data because the TF protocol v5 doesn't support nested objects
func parseProviderData(ctx context.Context, plan incidentImpactResourceModel) (providerData, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	var pdInput *pagerDutyInputResourceModel
	asDiags := plan.PagerDutyInput.As(ctx, &pdInput, basetypes.ObjectAsOptions{})
	diags.Append(asDiags...)

	var ddInput *dataDogInputResourceModel
	asDiags = plan.DataDogInput.As(ctx, &ddInput, basetypes.ObjectAsOptions{})
	diags.Append(asDiags...)

	var jiraInput *jiraInputResourceModel
	asDiags = plan.JiraInput.As(ctx, &jiraInput, basetypes.ObjectAsOptions{})
	diags.Append(asDiags...)

	var blamelessInput *blamelessInputResourceModel
	asDiags = plan.BlamelessInput.As(ctx, &blamelessInput, basetypes.ObjectAsOptions{})
	diags.Append(asDiags...)

	var statuspageInput *statuspageInputResourceModel
	asDiags = plan.StatusPageInput.As(ctx, &statuspageInput, basetypes.ObjectAsOptions{})
	diags.Append(asDiags...)

	var opsgenieInput *opsgenieInputResourceModel
	asDiags = plan.OpsGenieInput.As(ctx, &opsgenieInput, basetypes.ObjectAsOptions{})
	diags.Append(asDiags...)

	var firehydrantInput *firehydrantInputResourceModel
	asDiags = plan.FireHydrantInput.As(ctx, &firehydrantInput, basetypes.ObjectAsOptions{})
	diags.Append(asDiags...)

	var clubhouseInput *clubhouseInputResourceModel
	asDiags = plan.ClubhouseInput.As(ctx, &clubhouseInput, basetypes.ObjectAsOptions{})
	diags.Append(asDiags...)

	return providerData{
		pagerduty:   pdInput,
		datadog:     ddInput,
		jira:        jiraInput,
		blameless:   blamelessInput,
		statuspage:  statuspageInput,
		opsgenie:    opsgenieInput,
		firehydrant: firehydrantInput,
		clubhouse:   clubhouseInput,
	}, diags

}

func (iisr *incidentImpactSourceResource) Create(ctx context.Context, req resource.CreateRequest, res *resource.CreateResponse) {
	ctx = tflog.SetField(ctx, "resource", "incident_impact_source")
	ctx = tflog.SetField(ctx, "operation", "create")

	var plan incidentImpactResourceModel
	diags := req.Plan.Get(ctx, &plan)
	res.Diagnostics.Append(diags...)

	tflog.Info(ctx, "Creating IncidentImpactSource resource", map[string]any{"name": plan.Name.ValueString(), "projectSlug": plan.ProjectSlug.ValueString()})
	var pdInput *pagerDutyInputResourceModel

	tflog.Error(ctx, fmt.Sprintf("REMOVE ME! %+v", plan.PagerDutyInput.IsNull()))

	diags = plan.PagerDutyInput.As(ctx, &pdInput, basetypes.ObjectAsOptions{})
	res.Diagnostics.Append(diags...)

	pd, diags := parseProviderData(ctx, plan)

	res.Diagnostics.Append(diags...)

	if res.Diagnostics.HasError() {
		tflog.Error(ctx, "Error getting IncidentImpactSource plan", map[string]any{"diagnostics": res.Diagnostics})
		return
	}

	projectSlug := plan.ProjectSlug.ValueString()
	input, diags := getMutableIncidentImpactSourceStruct(ctx, plan, pd)
	res.Diagnostics.Append(diags...)

	iis, err := iisr.c.CreateIncidentImpactSource(ctx, input)
	tflog.Info(ctx, fmt.Sprintf("Created IncidentImpactSource %+v", iis), map[string]any{"iis": iis, "err": err})
	if err != nil {
		tflog.Error(ctx, "Error creating IncidentImpactSource", map[string]any{"error": err.Error()})
		res.Diagnostics.AddError(
			"Error creating IncidentImpactSource",
			fmt.Sprintf("Could not create code change soure, unexpected error: %+v", err.Error()),
		)
		return
	}

	state, diags := getNewStateFromIncidentImpactSource(ctx, iis, projectSlug, pd)
	res.Diagnostics.Append(diags...)
	diags = res.State.Set(ctx, state)
	res.Diagnostics.Append(diags...)
	tflog.Info(ctx, fmt.Sprintf("Successfully created IncidentImpactSource"), map[string]any{"diags": res.Diagnostics})
}

func (iisr *incidentImpactSourceResource) Read(ctx context.Context, req resource.ReadRequest, res *resource.ReadResponse) {
	ctx = tflog.SetField(ctx, "resource", "incident_impact_source")
	ctx = tflog.SetField(ctx, "operation", "read")

	var state incidentImpactResourceModel
	diags := req.State.Get(ctx, &state)
	res.Diagnostics.Append(diags...)

	pd, diags := parseProviderData(ctx, state)

	tflog.Info(ctx, "Reading IncidentImpactSource resource", map[string]any{"state": state})
	projectSlug := state.ProjectSlug.ValueString()
	slug := state.Slug.ValueString()
	if state.ProjectSlug.ValueString() == "" {
		id := state.ID.ValueString()
		splits := strings.Split(id, "/")
		projectSlug = splits[0]
		slug = splits[1]
	}

	ccs, err := iisr.c.GetIncidentImpactSource(ctx, projectSlug, slug)
	if err != nil {
		tflog.Error(ctx, "Error reading IncidentImpactSource", map[string]any{"error": err.Error()})
		res.Diagnostics.AddError(
			"Error reading IncidentImpactSource",
			fmt.Sprintf("Could not read code change soure, unexpected error: %+v", err.Error()),
		)
		return
	}
	newState, diags := getNewStateFromIncidentImpactSource(ctx, ccs, projectSlug, pd)
	res.Diagnostics.Append(diags...)

	diags = res.State.Set(ctx, newState)
	res.Diagnostics.Append(diags...)
}

func (iisr *incidentImpactSourceResource) Update(ctx context.Context, req resource.UpdateRequest, res *resource.UpdateResponse) {
	ctx = tflog.SetField(ctx, "resource", "incident_impact_source")
	ctx = tflog.SetField(ctx, "operation", "update")

	var state incidentImpactResourceModel
	diags := req.State.Get(ctx, &state)
	res.Diagnostics.Append(diags...)

	var plan incidentImpactResourceModel
	diags = req.Plan.Get(ctx, &plan)
	res.Diagnostics.Append(diags...)

	pd, diags := parseProviderData(ctx, plan)

	tflog.Info(ctx, "Creating IncidentImpactSource resource", map[string]any{"name": plan.Name.ValueString(), "projectSlug": plan.ProjectSlug.ValueString()})

	var pdInput *pagerDutyInputResourceModel
	diags = plan.PagerDutyInput.As(ctx, &pdInput, basetypes.ObjectAsOptions{})
	res.Diagnostics.Append(diags...)

	tflog.Info(ctx, "Updating IncidentImpactSource resource", map[string]any{"plan": plan})

	if res.Diagnostics.HasError() {
		tflog.Error(ctx, "Error getting IncidentImpactSource plan", map[string]any{"diagnostics": res.Diagnostics})
		return
	}

	projectSlug := plan.ProjectSlug.ValueString()
	inputFields, diags := getMutableIncidentImpactSourceStruct(ctx, plan, pd)
	res.Diagnostics.Append(diags...)
	input := gqlclient.IncidentImpactSourceInputUpdateType{
		Slug:                          state.Slug.ValueString(),
		IncidentImpactSourceInputType: inputFields,
	}

	ccs, err := iisr.c.UpdateIncidentImpactSource(ctx, input)
	if err != nil {
		tflog.Error(ctx, "Error updating IncidentImpactSource", map[string]any{"error": err.Error()})
		res.Diagnostics.AddError(
			"Error updating IncidentImpactSource",
			fmt.Sprintf("Could not update code change soure, unexpected error: %+v", err.Error()),
		)
		return
	}

	newState, diags := getNewStateFromIncidentImpactSource(ctx, ccs, projectSlug, pd)
	res.Diagnostics.Append(diags...)

	diags = res.State.Set(ctx, newState)
	res.Diagnostics.Append(diags...)
	tflog.Info(ctx, fmt.Sprintf("Successfully created IncidentImpactSource"), map[string]any{"diags": res.Diagnostics})

}

func (iisr *incidentImpactSourceResource) Delete(ctx context.Context, req resource.DeleteRequest, res *resource.DeleteResponse) {
	ctx = tflog.SetField(ctx, "resource", "incident_impact_source")
	ctx = tflog.SetField(ctx, "operation", "delete")

	var state incidentImpactResourceModel
	diags := req.State.Get(ctx, &state)
	res.Diagnostics.Append(diags...)

	tflog.Info(ctx, "Deleting IncidentImpactSource resource", map[string]any{"state": state})
	projectSlug := state.ProjectSlug.ValueStringPointer()
	slug := state.Slug.ValueStringPointer()

	err := iisr.c.DeleteImpactSource(projectSlug, slug)
	if err != nil {
		tflog.Error(ctx, "Error deleting IncidentImpactSource", map[string]any{"error": err.Error()})
		res.Diagnostics.AddError(
			"Error deleting IncidentImpactSource",
			fmt.Sprintf("Could not delete code change soure, unexpected error: %+v", err.Error()),
		)
		return
	}
}

func (iisr *incidentImpactSourceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, res *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, res)
}

func getNewStateFromIncidentImpactSource(ctx context.Context, iis *gqlclient.IncidentImpactSource, projectSlug string, data providerData) (incidentImpactResourceModel, diag.Diagnostics) {
	iirm := incidentImpactResourceModel{
		ID:               types.StringValue(iis.Slug),
		Slug:             types.StringValue(iis.Slug),
		ProjectSlug:      types.StringValue(projectSlug),
		EnvironmentName:  types.StringValue(iis.Environment.Name),
		Name:             types.StringValue(iis.Name),
		ProviderName:     types.StringValue(strings.ToLower(iis.Provider)),
		PagerDutyInput:   types.ObjectNull(pagerDutyInputResourceTFTypes),
		DataDogInput:     types.ObjectNull(dataDogInputResourceTFTypes),
		JiraInput:        types.ObjectNull(jiraInputResourceTFTypes),
		BlamelessInput:   types.ObjectNull(blamelessInputResourceTFTypes),
		StatusPageInput:  types.ObjectNull(statuspageInputResourceTFTypes),
		OpsGenieInput:    types.ObjectNull(opsgenieInputResourceTFTypes),
		FireHydrantInput: types.ObjectNull(firehydrantInputResourceTFTypes),
		ClubhouseInput:   types.ObjectNull(clubhouseInputResourceTFTypes),
	}

	return getProviderSpecificStateValue(ctx, iis, iirm, data)
}

func getProviderSpecificStateValue(ctx context.Context, iis *gqlclient.IncidentImpactSource, stateObj incidentImpactResourceModel, data providerData) (incidentImpactResourceModel, diag.Diagnostics) {
	diags := diag.Diagnostics{}
	pd := pagerDutyInputResourceModel{
		RemoteUrgency:   types.StringValue(iis.ProviderData.PagerDutyProviderData.RemoteUrgency),
		RemoteServices:  types.StringValue(iis.ProviderData.PagerDutyProviderData.RemoteServices),
		IntegrationSlug: types.StringNull(),
	}

	if iis.IntegrationAuthSlug != "" {
		pd.IntegrationSlug = types.StringValue(iis.IntegrationAuthSlug)
	}

	dd := dataDogInputResourceModel{
		Query:                   types.StringValue(iis.ProviderData.DataDogProviderData.Query),
		RemotePriorityThreshold: types.StringValue(iis.ProviderData.DataDogProviderData.RemotePriorityThreshold),
		IntegrationSlug:         types.StringNull(),
	}

	if iis.IntegrationAuthSlug != "" {
		dd.IntegrationSlug = types.StringValue(iis.IntegrationAuthSlug)

	}

	jira := jiraInputResourceModel{
		RemoteJQL:       types.StringValue(iis.ProviderData.JiraProviderData.RemoteJql),
		IntegrationSlug: types.StringNull(),
	}
	if iis.IntegrationAuthSlug != "" {
		jira.IntegrationSlug = types.StringValue(iis.IntegrationAuthSlug)
	}

	var t []attr.Value
	for _, remoteType := range iis.ProviderData.BlamelessProviderData.RemoteTypes {
		typ := types.StringValue(remoteType)
		t = append(t, typ)
	}
	sv, errDiag := types.SetValue(types.StringType, t)

	diags.Append(errDiag...)

	blameless := blamelessInputResourceModel{
		RemoteSeverityThreshold: types.StringValue(iis.ProviderData.BlamelessProviderData.RemoteSeverityThreshold),
		RemoteTypes:             sv,
	}

	statuspage := statuspageInputResourceModel{
		RemotePage:                 types.StringValue(iis.ProviderData.StatuspageProviderData.RemotePage),
		RemoteComponent:            types.StringValue(iis.ProviderData.StatuspageProviderData.RemoteComponent),
		RemoteImpact:               types.StringValue(iis.ProviderData.StatuspageProviderData.RemoteImpact),
		IgnoreMaintenanceIncidents: types.BoolValue(iis.ProviderData.StatuspageProviderData.IgnoreMaintenanceIncidents),
		IntegrationSlug:            types.StringNull(),
	}
	if iis.IntegrationAuthSlug != "" {
		statuspage.IntegrationSlug = types.StringValue(iis.IntegrationAuthSlug)
	}

	opsgenie := opsgenieInputResourceModel{
		RemoteAlertTags:         types.StringValue(iis.ProviderData.OpsGenieProviderData.RemoteAlertTags),
		RemoteIncidentTags:      types.StringValue(iis.ProviderData.OpsGenieProviderData.RemoteIncidentTags),
		RemotePriorityThreshold: types.StringValue(iis.ProviderData.OpsGenieProviderData.RemotePriorityThreshold),
		RemoteService:           types.StringValue(iis.ProviderData.OpsGenieProviderData.RemoteService),
		RemoteUseAlerts:         types.BoolValue(iis.ProviderData.OpsGenieProviderData.RemoteUseAlerts),
		IntegrationSlug:         types.StringNull(),
	}
	if iis.IntegrationAuthSlug != "" {
		opsgenie.IntegrationSlug = types.StringValue(iis.IntegrationAuthSlug)
	}

	firehydrant := firehydrantInputResourceModel{
		RemoteEnvironments:       types.StringValue(iis.ProviderData.FireHydrantProviderData.RemoteEnvironments),
		RemoteServices:           types.StringValue(iis.ProviderData.FireHydrantProviderData.RemoteServices),
		RemoteMitigatedIsHealthy: types.BoolValue(iis.ProviderData.FireHydrantProviderData.RemoteMitigatedIsHealthy),
	}

	clubhouse := clubhouseInputResourceModel{
		RemoteQuery:     types.StringValue(iis.ProviderData.ClubhouseProviderData.RemoteQuery),
		IntegrationSlug: types.StringNull(),
	}
	if iis.IntegrationAuthSlug != "" {
		clubhouse.IntegrationSlug = types.StringValue(iis.IntegrationAuthSlug)
	}

	pdObj, errDiag := types.ObjectValueFrom(ctx, pagerDutyInputResourceTFTypes, pd)
	diags.Append(errDiag...)

	ddObj, errDiag := types.ObjectValueFrom(ctx, dataDogInputResourceTFTypes, dd)
	diags.Append(errDiag...)

	jiraObj, errDiag := types.ObjectValueFrom(ctx, jiraInputResourceTFTypes, jira)
	diags.Append(errDiag...)

	blamelessObj, errDiag := types.ObjectValueFrom(ctx, blamelessInputResourceTFTypes, blameless)
	diags.Append(errDiag...)

	statuspageObj, errDiag := types.ObjectValueFrom(ctx, statuspageInputResourceTFTypes, statuspage)
	diags.Append(errDiag...)

	opsgenieObj, errDiag := types.ObjectValueFrom(ctx, opsgenieInputResourceTFTypes, opsgenie)
	diags.Append(errDiag...)

	firehydrantObj, errDiag := types.ObjectValueFrom(ctx, firehydrantInputResourceTFTypes, firehydrant)
	diags.Append(errDiag...)

	clubhouseObj, errDiag := types.ObjectValueFrom(ctx, clubhouseInputResourceTFTypes, clubhouse)
	diags.Append(errDiag...)

	if data.pagerduty != nil {
		stateObj.PagerDutyInput = pdObj
	}
	if data.datadog != nil {
		stateObj.DataDogInput = ddObj
	}
	if data.jira != nil {
		stateObj.JiraInput = jiraObj
	}
	if data.blameless != nil {
		stateObj.BlamelessInput = blamelessObj
	}
	if data.statuspage != nil {
		stateObj.StatusPageInput = statuspageObj
	}
	if data.opsgenie != nil {
		stateObj.OpsGenieInput = opsgenieObj
	}
	if data.firehydrant != nil {
		stateObj.FireHydrantInput = firehydrantObj
	}
	if data.clubhouse != nil {
		stateObj.ClubhouseInput = clubhouseObj
	}

	return stateObj, diags

}

func getMutableIncidentImpactSourceStruct(ctx context.Context, plan incidentImpactResourceModel, data providerData) (gqlclient.IncidentImpactSourceInputType, diag.Diagnostics) {
	input := gqlclient.IncidentImpactSourceInputType{
		ProjectSlug:          plan.ProjectSlug.ValueString(),
		EnvironmentName:      strings.ToLower(plan.EnvironmentName.ValueString()),
		Name:                 plan.Name.ValueString(),
		Provider:             strings.ToUpper(plan.ProviderName.ValueString()),
		PagerDutyInputType:   nil,
		DataDogInputType:     nil,
		JiraInputType:        nil,
		BlamelessInputType:   nil,
		StatuspageInputType:  nil,
		OpsGenieInputType:    nil,
		FireHydrantInputType: nil,
		ClubhouseInputType:   nil,
	}

	return getProviderSpecificData(ctx, input, data)
}
func getProviderSpecificData(ctx context.Context, input gqlclient.IncidentImpactSourceInputType, data providerData) (gqlclient.IncidentImpactSourceInputType, diag.Diagnostics) {
	if data.pagerduty != nil {
		input.PagerDutyInputType = &gqlclient.PagerDutyInputType{
			RemoteServices: data.pagerduty.RemoteServices.ValueString(),
			RemoteUrgency:  data.pagerduty.RemoteUrgency.ValueString(),
		}
	}

	if data.datadog != nil {
		input.DataDogInputType = &gqlclient.DataDogInputType{
			DataDogProviderData: gqlclient.DataDogProviderData{
				Query:                   data.datadog.Query.ValueString(),
				RemotePriorityThreshold: data.datadog.RemotePriorityThreshold.ValueString(),
			},
			IntegrationSlug: data.datadog.IntegrationSlug.ValueString(),
		}
	}

	if data.jira != nil {
		input.JiraInputType = &gqlclient.JiraInputType{
			JiraProviderData: gqlclient.JiraProviderData{
				RemoteJql: data.jira.RemoteJQL.ValueString(),
			},
			IntegrationSlug: data.jira.IntegrationSlug.ValueString(),
		}
	}

	diags := diag.Diagnostics{}
	if data.blameless != nil {
		var remoteTypes []string
		diags = data.blameless.RemoteTypes.ElementsAs(context.Background(), &remoteTypes, false)
		if diags.HasError() {
			tflog.Error(ctx, "Error parsing remote types", map[string]any{"error": diags})
		}
		input.BlamelessInputType = &gqlclient.BlamelessInputType{
			BlamelessProviderData: gqlclient.BlamelessProviderData{
				RemoteTypes:             remoteTypes,
				RemoteSeverityThreshold: data.blameless.RemoteSeverityThreshold.ValueString(),
			},
			IntegrationSlug: data.blameless.IntegrationSlug.ValueString(),
		}
	}

	if data.statuspage != nil {
		input.StatuspageInputType = &gqlclient.StatuspageInputType{
			StatuspageProviderData: gqlclient.StatuspageProviderData{
				RemotePage:                 data.statuspage.RemotePage.ValueString(),
				RemoteComponent:            data.statuspage.RemoteComponent.ValueString(),
				RemoteImpact:               data.statuspage.RemoteImpact.ValueString(),
				IgnoreMaintenanceIncidents: data.statuspage.IgnoreMaintenanceIncidents.ValueBool(),
			},
			IntegrationSlug: data.statuspage.IntegrationSlug.ValueString(),
		}
	}

	if data.opsgenie != nil {
		input.OpsGenieInputType = &gqlclient.OpsGenieInputType{
			OpsGenieProviderData: gqlclient.OpsGenieProviderData{
				RemoteAlertTags:         data.opsgenie.RemoteAlertTags.ValueString(),
				RemoteIncidentTags:      data.opsgenie.RemoteIncidentTags.ValueString(),
				RemotePriorityThreshold: data.opsgenie.RemotePriorityThreshold.ValueString(),
				RemoteService:           data.opsgenie.RemoteService.ValueString(),
				RemoteUseAlerts:         data.opsgenie.RemoteUseAlerts.ValueBool(),
			},
			IntegrationSlug: data.opsgenie.IntegrationSlug.ValueString(),
		}
	}

	if data.firehydrant != nil {
		input.FireHydrantInputType = &gqlclient.FireHydrantInputType{
			FireHydrantProviderData: gqlclient.FireHydrantProviderData{
				RemoteEnvironments:       data.firehydrant.RemoteEnvironments.ValueString(),
				RemoteServices:           data.firehydrant.RemoteServices.ValueString(),
				RemoteMitigatedIsHealthy: data.firehydrant.RemoteMitigatedIsHealthy.ValueBool(),
			},
		}
	}

	if data.clubhouse != nil {
		input.ClubhouseInputType = &gqlclient.ClubhouseInputType{
			ClubhouseProviderData: gqlclient.ClubhouseProviderData{
				RemoteQuery: data.clubhouse.RemoteQuery.ValueString(),
			},
			IntegrationSlug: data.clubhouse.IntegrationSlug.ValueString(),
		}
	}

	return input, diags
}

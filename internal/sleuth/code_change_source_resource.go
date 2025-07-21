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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/sleuth-io/terraform-provider-sleuth/internal/gqlclient"
)

var (
	_ resource.Resource                = &codeChangeSourceResource{}
	_ resource.ResourceWithConfigure   = &codeChangeSourceResource{}
	_ resource.ResourceWithImportState = &codeChangeSourceResource{}
)

type codeChangeResourceModel struct {
	ProjectSlug types.String `tfsdk:"project_slug"`
	Name        types.String `tfsdk:"name"`
	Slug        types.String `tfsdk:"slug"`
	ID          types.String `tfsdk:"id"`

	Repository          *repositoryResourceModel           `tfsdk:"repository"`
	EnvironmentMappings []environmentMappingsResourceModel `tfsdk:"environment_mappings"`
	BuildMappings       []buildMappingsResourceModel       `tfsdk:"build_mappings"`

	DeployTrackingType types.String `tfsdk:"deploy_tracking_type"`
	CollectImpact      types.Bool   `tfsdk:"collect_impact"`
	PathPrefix         types.String `tfsdk:"path_prefix"`
	NotifyInSlack      types.Bool   `tfsdk:"notify_in_slack"`
	IncludeInDashboard types.Bool   `tfsdk:"include_in_dashboard"`
	AutoTrackingDelay  types.Int64  `tfsdk:"auto_tracking_delay"`
}

type repositoryResourceModel struct {
	Owner           types.String `tfsdk:"owner"`
	Name            types.String `tfsdk:"name"`
	URL             types.String `tfsdk:"url"`
	Provider        types.String `tfsdk:"provider"`
	IntegrationSlug types.String `tfsdk:"integration_slug"`
	RepoUID         types.String `tfsdk:"repo_uid"`
	ProjectUID      types.String `tfsdk:"project_uid"`
	Webhook         types.Object `tfsdk:"webhook"`
}

type webhookResourceModel struct {
	URL    types.String `tfsdk:"url"`
	Secret types.String `tfsdk:"secret"`
}

func (w webhookResourceModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"url":    types.StringType,
		"secret": types.StringType,
	}
}

type environmentMappingsResourceModel struct {
	EnvironmentSlug types.String `tfsdk:"environment_slug"`
	Branch          types.String `tfsdk:"branch"`
	ID              types.String `tfsdk:"id"`
}

type buildMappingsResourceModel struct {
	EnvironmentSlug          types.String `tfsdk:"environment_slug"`
	Provider                 types.String `tfsdk:"provider"`
	IntegrationSlug          types.String `tfsdk:"integration_slug"`
	BuildName                types.String `tfsdk:"build_name"`
	JobName                  types.String `tfsdk:"job_name"`
	ProjectKey               types.String `tfsdk:"project_key"`
	ProjectName              types.String `tfsdk:"project_name"`
	MatchBranchToEnvironment types.Bool   `tfsdk:"match_branch_to_environment"`
	IsCustom                 types.Bool   `tfsdk:"is_custom"`
}

const azureProvider = "azure"

type codeChangeSourceResource struct {
	c *gqlclient.Client
}

func NewCodeChangeSourceResource() resource.Resource {
	return &codeChangeSourceResource{}
}

func (ccsr *codeChangeSourceResource) Schema(_ context.Context, _ resource.SchemaRequest, res *resource.SchemaResponse) {
	res.Schema = schema.Schema{
		MarkdownDescription: "Sleuth code change source.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"project_slug": schema.StringAttribute{
				MarkdownDescription: "The slug of the project that this code change source belongs to.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(), // ForceNew replacement
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Code change source name",
				Required:            true,
			},
			"slug": schema.StringAttribute{
				Computed: true,
			},
			"deploy_tracking_type": schema.StringAttribute{
				MarkdownDescription: "How to track deploys. Valid choices are build, manual, auto_pr, auto_tag, auto_push",
				Required:            true,
			},
			"collect_impact": schema.BoolAttribute{
				MarkdownDescription: "Whether to collect impact for its deploys",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"path_prefix": schema.StringAttribute{
				MarkdownDescription: "What code source path to limit this deployment to. Useful for monorepos. Must be used with the [jsonencode()](https://developer.hashicorp.com/terraform/language/functions/jsonencode) function to specify the paths to include and/or exclude respectively. (see example above)",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"notify_in_slack": schema.BoolAttribute{
				MarkdownDescription: "Whether to send Slack notifications for deploys or not",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"include_in_dashboard": schema.BoolAttribute{
				MarkdownDescription: "Whether to include deploys from this change source in the metrics dashboard",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"auto_tracking_delay": schema.Int64Attribute{
				MarkdownDescription: "The delay to add to a deployment event",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(0),
			},
			"repository": schema.SingleNestedAttribute{
				Description: "Repository details",
				Required:    true,
				Attributes: map[string]schema.Attribute{
					"owner": schema.StringAttribute{
						MarkdownDescription: "The repository owner, usually the organization or user name",
						Required:            true,
					},
					"name": schema.StringAttribute{
						MarkdownDescription: "The repository name",
						Required:            true,
					},
					"url": schema.StringAttribute{
						MarkdownDescription: "The repository URL, used for links",
						Required:            true,
					},
					"provider": schema.StringAttribute{
						MarkdownDescription: "The repository provider, options: AZURE, BITBUCKET, CUSTOM_GIT, GITHUB, GITHUB_ENTERPRISE, GITLAB",
						Required:            true,
					},
					"integration_slug": schema.StringAttribute{
						MarkdownDescription: "IntegrationAuthentication slug used",
						Optional:            true,
						Computed:            true,
					},
					"repo_uid": schema.StringAttribute{
						MarkdownDescription: "Repository UID, required only for AZURE provider. You can obtain data from [API](https://learn.microsoft.com/en-us/rest/api/azure/devops/git/repositories/list?view=azure-devops-rest-6.0&tabs=HTTP)",
						Optional:            true,
					},
					"project_uid": schema.StringAttribute{
						MarkdownDescription: "Project UID, required only for AZURE provider. You can obtain data from [API](https://learn.microsoft.com/en-us/rest/api/azure/devops/git/repositories/list?view=azure-devops-rest-6.0&tabs=HTTP)",
						Optional:            true,
					},
					"webhook": schema.SingleNestedAttribute{
						MarkdownDescription: "Webhook configuration for registering deploys from code integrations in read-only mode",
						Computed:            true,
						Attributes: map[string]schema.Attribute{
							"url": schema.StringAttribute{
								MarkdownDescription: "Webhook URL",
								Computed:            true,
							},
							"secret": schema.StringAttribute{
								MarkdownDescription: "Webhook secret to present in payloads sent to the webhook URL",
								Computed:            true,
								Sensitive:           true,
							},
						},
					},
				},
			},
			"environment_mappings": schema.ListNestedAttribute{
				MarkdownDescription: "Environment mappings of the code change source. They must be ordered by environment_slug ascending to avoid Terraform plan changes.",
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"environment_slug": schema.StringAttribute{
							MarkdownDescription: "The environment slug for mapping",
							Required:            true,
						},
						"branch": schema.StringAttribute{
							MarkdownDescription: "The repository branch name for the environment",
							Required:            true,
						},
						"id": schema.StringAttribute{
							MarkdownDescription: "Computed ID",
							Computed:            true,
						},
					},
				},
			},
			"build_mappings": schema.ListNestedAttribute{
				MarkdownDescription: "Build mappings of the code change source. They must be ordered by environment_slug ascending to avoid Terraform plan changes.",
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"environment_slug": schema.StringAttribute{
							MarkdownDescription: "The environment slug",
							Required:            true,
						},
						"provider": schema.StringAttribute{
							MarkdownDescription: "The build provider. Options: AZURE, BITBUCKET_PIPELINES, BUILDKITE, CIRCLECI, GITHUB, GITLAB, JENKINS",
							Required:            true,
						},
						"integration_slug": schema.StringAttribute{
							MarkdownDescription: "IntegrationAuthentication slug used",
							Optional:            true,
							Computed:            true,
						},
						"build_name": schema.StringAttribute{
							MarkdownDescription: "The remote build or pipeline name",
							Required:            true,
						},
						"job_name": schema.StringAttribute{
							MarkdownDescription: "The job or stage within the build or pipeline, if supported",
							Optional:            true,
						},
						"project_key": schema.StringAttribute{
							MarkdownDescription: "The build project key. If both project_key and project_name are provided, project_key takes precedence.",
							Optional:            true,
						},
						"project_name": schema.StringAttribute{
							MarkdownDescription: "The build project name. If both project_key and project_name are provided, project_key takes precedence.",
							Optional:            true,
						},
						"match_branch_to_environment": schema.BoolAttribute{
							MarkdownDescription: "Whether only builds performed on the branch mapped from the environment are " +
								"tracked or not. Basically if you only want Sleuth to find builds that were triggered" +
								"by a change on the branch that is configured for the environment, set this to false. " +
								"Defaults to true",
							Optional: true,
							Computed: true,
							Default:  booldefault.StaticBool(true),
						},
						"is_custom": schema.BoolAttribute{
							MarkdownDescription: "Whether this is a custom build mapping or not. This needs to be set to true " +
								"if a build name or job name isn't visible in Sleuth. Defaults to false",
							Optional: true,
							Computed: true,
							Default:  booldefault.StaticBool(false),
						},
					},
				},
			},
		},
	}
}

func (ccsr *codeChangeSourceResource) Configure(ctx context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	ccsr.c = req.ProviderData.(*gqlclient.Client)
}

func (ccsr *codeChangeSourceResource) Metadata(ctx context.Context, req resource.MetadataRequest, res *resource.MetadataResponse) {
	res.TypeName = req.ProviderTypeName + "_code_change_source"
}

func (ccsr *codeChangeSourceResource) Create(ctx context.Context, req resource.CreateRequest, res *resource.CreateResponse) {
	ctx = tflog.SetField(ctx, "resource", "CodeChangeSource")
	ctx = tflog.SetField(ctx, "operation", "create")

	var plan codeChangeResourceModel
	diags := req.Plan.Get(ctx, &plan)
	res.Diagnostics.Append(diags...)

	if res.Diagnostics.HasError() {
		tflog.Error(ctx, "Error getting CodeChangeSource plan", map[string]any{"diagnostics": res.Diagnostics})
		return
	}

	tflog.Info(ctx, "Creating CodeChangeSource resource", map[string]any{"name": plan.Name.ValueString(), "projectSlug": plan.ProjectSlug.ValueString()})

	projectSlug := plan.ProjectSlug.ValueString()
	inputFields, err := getMutableCodeChangeSourceStruct(plan)
	if err != nil {
		res.Diagnostics.AddError("Could not create input object", fmt.Sprintf("Could not create input object: %+v", err.Error()))
		return
	}

	input := gqlclient.CreateCodeChangeSourceMutationInput{
		ProjectSlug:             projectSlug,
		InitializeChanges:       true,
		MutableCodeChangeSource: inputFields,
	}

	if err := validateCodeChangeInput(*inputFields); err != nil {
		tflog.Error(ctx, "Error validating CodeChangeSource input", map[string]any{"err": err.Error()})
		res.Diagnostics.AddError("Error validating input", fmt.Sprintf("Input object validation failed, err: %+v", err.Error()))

		return
	}

	ccs, err := ccsr.c.CreateCodeChangeSource(ctx, input)
	if err != nil {
		tflog.Error(ctx, "Error creating CodeChangeSource", map[string]any{"error": err.Error()})
		res.Diagnostics.AddError(
			"Error creating CodeChangeSource",
			fmt.Sprintf("Could not create code change source, unexpected error: %+v", err.Error()),
		)
		return
	}

	state, diags := getNewStateFromCodeChangeSource(ctx, ccs, projectSlug, plan)

	// if they are both empty, make sure they match (could be [] or nil)
	// if len(plan.BuildMappings) < 1 && len(state.BuildMappings) < 1 {
	// 	state.BuildMappings = plan.BuildMappings
	// }

	res.Diagnostics.Append(diags...)

	diags = res.State.Set(ctx, state)
	res.Diagnostics.Append(diags...)
	tflog.Info(ctx, "Successfully created CodeChangeSource", map[string]any{"diags": res.Diagnostics})
}

func (ccsr *codeChangeSourceResource) Read(ctx context.Context, req resource.ReadRequest, res *resource.ReadResponse) {
	ctx = tflog.SetField(ctx, "resource", "code_change_source")
	ctx = tflog.SetField(ctx, "operation", "read")

	var state codeChangeResourceModel
	diags := req.State.Get(ctx, &state)
	res.Diagnostics.Append(diags...)

	if res.Diagnostics.HasError() {
		tflog.Error(ctx, "Error reading CodeChangeSource", map[string]any{"diagnostics": res.Diagnostics})
		return
	}

	tflog.Info(ctx, "Reading CodeChangeSource resource", map[string]any{"state": state})
	projectSlug := state.ProjectSlug.ValueString()
	slug := state.Slug.ValueString()
	// when importing a resource, only ID will be set with project slug & slug
	if projectSlug == "" {
		id := state.ID.ValueString()
		splits := strings.Split(id, "/")
		if len(splits) != 2 {
			res.Diagnostics.AddError("Error importing CodeChangeSource", "Imported code change source must have an ID of the form 'project_slug/change_source_slug'")
			return
		}

		projectSlug = splits[0]
		slug = splits[1]
	}

	ccs, err := ccsr.c.GetCodeChangeSource(ctx, &projectSlug, &slug)
	if err != nil {
		tflog.Error(ctx, "Error reading CodeChangeSource", map[string]any{"error": err.Error()})
		res.Diagnostics.AddError(
			"Error reading CodeChangeSource",
			fmt.Sprintf("Could not read code change source, unexpected error: %+v", err.Error()),
		)
		return
	}
	newState, diags := getNewStateFromCodeChangeSource(ctx, ccs, projectSlug, state)
	res.Diagnostics.Append(diags...)

	diags = res.State.Set(ctx, newState)
	res.Diagnostics.Append(diags...)

}

func (ccsr *codeChangeSourceResource) Update(ctx context.Context, req resource.UpdateRequest, res *resource.UpdateResponse) {
	ctx = tflog.SetField(ctx, "resource", "code_change_source")
	ctx = tflog.SetField(ctx, "operation", "update")

	var state codeChangeResourceModel
	diags := req.State.Get(ctx, &state)

	res.Diagnostics.Append(diags...)

	var plan codeChangeResourceModel
	diags = req.Plan.Get(ctx, &plan)
	res.Diagnostics.Append(diags...)

	tflog.Info(ctx, "Creating CodeChangeSource resource", map[string]any{"name": plan.Name.ValueString(), "projectSlug": plan.ProjectSlug.ValueString()})
	tflog.Info(ctx, "Updating CodeChangeSource resource", map[string]any{"plan": plan})

	if res.Diagnostics.HasError() {
		tflog.Error(ctx, "Error getting CodeChangeSource plan", map[string]any{"diagnostics": res.Diagnostics})
		return
	}

	projectSlug := plan.ProjectSlug.ValueString()
	inputFields, err := getMutableCodeChangeSourceStruct(plan)
	if err != nil {
		res.Diagnostics.AddError("Could not create input object", fmt.Sprintf("Could not create input object: %+v", err.Error()))
		return
	}

	input := gqlclient.UpdateCodeChangeSourceMutationInput{
		ProjectSlug:             projectSlug,
		Slug:                    state.Slug.ValueString(),
		MutableCodeChangeSource: inputFields,
	}

	if err := validateCodeChangeInput(*inputFields); err != nil {
		tflog.Error(ctx, "Error validating CodeChangeSource input", map[string]any{"err": err.Error()})
		res.Diagnostics.AddError("Error validating input", fmt.Sprintf("Input object validation failed, err: %+v", err.Error()))

		return
	}

	ccs, err := ccsr.c.UpdateCodeChangeSource(ctx, input)
	if err != nil {
		tflog.Error(ctx, "Error updating CodeChangeSource", map[string]any{"error": err.Error()})
		res.Diagnostics.AddError(
			"Error updating CodeChangeSource",
			fmt.Sprintf("Could not update code change soure, unexpected error: %+v", err.Error()),
		)
		return
	}

	newState, diags := getNewStateFromCodeChangeSource(ctx, ccs, projectSlug, plan)

	res.Diagnostics.Append(diags...)

	diags = res.State.Set(ctx, newState)
	res.Diagnostics.Append(diags...)
	tflog.Info(ctx, "Successfully created CodeChangeSource", map[string]any{"diags": res.Diagnostics})

}

func (ccsr *codeChangeSourceResource) Delete(ctx context.Context, req resource.DeleteRequest, res *resource.DeleteResponse) {
	ctx = tflog.SetField(ctx, "resource", "code_change_source")
	ctx = tflog.SetField(ctx, "operation", "delete")

	var state codeChangeResourceModel
	diags := req.State.Get(ctx, &state)
	res.Diagnostics.Append(diags...)

	tflog.Info(ctx, "Deleting CodeChangeSource resource", map[string]any{"state": state})
	projectSlug := state.ProjectSlug.ValueStringPointer()
	slug := state.Slug.ValueStringPointer()

	err := ccsr.c.DeleteChangeSource(ctx, projectSlug, slug)
	if err != nil {
		tflog.Error(ctx, "Error deleting CodeChangeSource", map[string]any{"error": err.Error()})
		res.Diagnostics.AddError(
			"Error deleting CodeChangeSource",
			fmt.Sprintf("Could not delete code change soure, unexpected error: %+v", err.Error()),
		)
		return
	}
}

func (ccsr *codeChangeSourceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, res *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, res)
}

func validateCodeChangeInput(ccs gqlclient.MutableCodeChangeSource) error {
	if strings.ToLower(ccs.Repository.Provider) != azureProvider {
		return nil
	}

	if ccs.Repository.ProjectUID == "" || ccs.Repository.RepoUID == "" || ccs.Repository.IntegrationSlug == "" {
		return fmt.Errorf("project_uid, repo_uid and integration_slug are required for AZURE provider")
	}
	return nil
}

func getNewStateFromCodeChangeSource(ctx context.Context, ccs *gqlclient.CodeChangeSource, projectSlug string, plan codeChangeResourceModel) (codeChangeResourceModel, diag.Diagnostics) {
	var environmentMappings []environmentMappingsResourceModel

	diags := diag.Diagnostics{}

	for _, em := range ccs.EnvironmentMappings {
		emm := environmentMappingsResourceModel{
			EnvironmentSlug: types.StringValue(em.EnvironmentSlug),
			Branch:          types.StringValue(em.Branch),
			ID:              types.StringValue(fmt.Sprintf("%s/%s", projectSlug, em.EnvironmentSlug)),
		}
		environmentMappings = append(environmentMappings, emm)
	}

	var buildMappings []buildMappingsResourceModel = []buildMappingsResourceModel{}
	// Build a lookup from plan for project_name preservation
	planBuildMappingLookup := map[string]buildMappingsResourceModel{}
	for _, pbm := range plan.BuildMappings {
		key := pbm.EnvironmentSlug.ValueString() + ":" + pbm.BuildName.ValueString()
		planBuildMappingLookup[key] = pbm
	}
	for _, bm := range ccs.DeployTrackingBuildMappings {
		key := bm.Environment.Slug + ":" + bm.BuildName
		planBM, hasPlan := planBuildMappingLookup[key]
		projectName := ""
		if hasPlan {
			projectName = planBM.ProjectName.ValueString()
		}
		// Find the original provider from the plan for this build mapping
		originalProvider := strings.ToUpper(bm.Provider) // fallback to uppercase API response
		for _, planBM := range plan.BuildMappings {
			if planBM.EnvironmentSlug.ValueString() == bm.Environment.Slug && planBM.BuildName.ValueString() == bm.BuildName {
				originalProvider = planBM.Provider.ValueString()
				break
			}
		}
		buildMappingObj := buildMappingsResourceModel{
			EnvironmentSlug:          types.StringValue(bm.Environment.Slug),
			Provider:                 types.StringValue(originalProvider),
			IntegrationSlug:          types.StringValue(bm.IntegrationSlug),
			BuildName:                types.StringValue(bm.BuildName),
			JobName:                  types.StringValue(bm.JobName),
			ProjectKey:               types.StringNull(),
			ProjectName:              types.StringNull(),
			MatchBranchToEnvironment: types.BoolValue(bm.MatchBranchToEnvironment),
			IsCustom:                 types.BoolValue(bm.IsCustom),
		}

		if projectName != "" {
			buildMappingObj.ProjectName = types.StringValue(projectName)
		}

		if bm.IntegrationSlug == "" {
			buildMappingObj.IntegrationSlug = types.StringNull()
		}

		if bm.JobName == "" {
			buildMappingObj.JobName = types.StringNull()
		}

		if hasPlan && !planBM.ProjectKey.IsNull() && bm.BuildProjectKey != "" {
			buildMappingObj.ProjectKey = types.StringValue(bm.BuildProjectKey)
		}

		buildMappings = append(buildMappings, buildMappingObj)
	}

	if len(buildMappings) < 1 && len(plan.BuildMappings) < 1 {
		buildMappings = plan.BuildMappings
	}

	r := repositoryResourceModel{
		Owner:           types.StringValue(ccs.Repository.Owner),
		Name:            types.StringValue(ccs.Repository.Name),
		URL:             types.StringValue(ccs.Repository.Url),
		Provider:        types.StringValue(plan.Repository.Provider.ValueString()), // preserve original case
		IntegrationSlug: types.StringNull(),
		RepoUID:         types.StringNull(),
		ProjectUID:      types.StringNull(),
		Webhook:         types.ObjectNull(webhookResourceModel{}.AttributeTypes()),
	}

	if ccs.Repository.IntegrationAuth != nil {
		r.IntegrationSlug = types.StringValue(ccs.Repository.IntegrationAuth.Slug)
	}

	if ccs.Repository.Webhook != nil {
		webhook := webhookResourceModel{
			URL:    types.StringValue(ccs.Repository.Webhook.Url),
			Secret: types.StringValue(ccs.Repository.Webhook.Secret),
		}

		var webhookDiags diag.Diagnostics
		r.Webhook, webhookDiags = types.ObjectValueFrom(ctx, webhook.AttributeTypes(), webhook)
		diags.Append(webhookDiags...)
	}

	if strings.ToLower(ccs.Repository.Provider) == azureProvider {
		r.RepoUID = types.StringValue(ccs.Repository.RepoUID)
		r.ProjectUID = types.StringValue(ccs.Repository.ProjectUID)
	}

	return codeChangeResourceModel{
		ProjectSlug:         types.StringValue(projectSlug),
		Name:                types.StringValue(ccs.Name),
		Slug:                types.StringValue(ccs.Slug),
		ID:                  types.StringValue(ccs.Slug),
		Repository:          &r,
		EnvironmentMappings: environmentMappings,
		BuildMappings:       buildMappings,
		DeployTrackingType:  types.StringValue(ccs.DeployTrackingType),
		CollectImpact:       types.BoolValue(ccs.CollectImpact),
		PathPrefix:          types.StringValue(ccs.PathPrefix),
		NotifyInSlack:       types.BoolValue(ccs.NotifyInSlack),
		IncludeInDashboard:  types.BoolValue(ccs.IncludeInDashboard),
		AutoTrackingDelay:   types.Int64Value(int64(ccs.AutoTrackingDelay)),
	}, diags
}

func getMutableCodeChangeSourceStruct(plan codeChangeResourceModel) (*gqlclient.MutableCodeChangeSource, error) {
	var environmentMappings []gqlclient.BranchMapping
	environmentMappingsLookup := map[string]string{}
	for _, em := range plan.EnvironmentMappings {
		environmentSlug := em.EnvironmentSlug.ValueString()
		branch := em.Branch.ValueString()
		environmentMappings = append(environmentMappings, gqlclient.BranchMapping{
			EnvironmentSlug: environmentSlug,
			Branch:          em.Branch.ValueString(),
		})
		environmentMappingsLookup[environmentSlug] = branch
	}

	var buildMappingsT []gqlclient.BuildMapping
	for _, bm := range plan.BuildMappings {
		environmentSlug := bm.EnvironmentSlug.ValueString()
		buildBranch, ok := environmentMappingsLookup[environmentSlug]
		if !ok {
			return nil, fmt.Errorf("could not find branch for build mapping for environment slug: %s. Did you forget to include this or all environments in the `environment_mappings` field?", environmentSlug)
		}

		projectKey := bm.ProjectKey.ValueString()
		projectName := bm.ProjectName.ValueString()

		buildMapping := gqlclient.BuildMapping{
			EnvironmentSlug:          environmentSlug,
			Provider:                 strings.ToUpper(bm.Provider.ValueString()),
			BuildName:                bm.BuildName.ValueString(),
			JobName:                  bm.JobName.ValueString(),
			IntegrationSlug:          bm.IntegrationSlug.ValueString(),
			BuildBranch:              buildBranch,
			MatchBranchToEnvironment: bm.MatchBranchToEnvironment.ValueBool(),
			IsCustom:                 bm.IsCustom.ValueBool(),
		}

		if projectKey != "" {
			buildMapping.BuildProjectKey = projectKey
		}
		if projectName != "" {
			buildMapping.BuildProjectName = projectName
		}

		buildMappingsT = append(buildMappingsT, buildMapping)
	}

	return &gqlclient.MutableCodeChangeSource{
		Name: plan.Name.ValueString(),
		Repository: gqlclient.MutableRepository{
			RepositoryBase: gqlclient.RepositoryBase{
				Owner:      plan.Repository.Owner.ValueString(),
				Name:       plan.Repository.Name.ValueString(),
				Provider:   strings.ToUpper(plan.Repository.Provider.ValueString()),
				Url:        plan.Repository.URL.ValueString(),
				ProjectUID: plan.Repository.ProjectUID.ValueString(),
				RepoUID:    plan.Repository.RepoUID.ValueString(),
			},
			IntegrationSlug: plan.Repository.IntegrationSlug.ValueString(),
		},
		DeployTrackingType:  plan.DeployTrackingType.ValueString(),
		CollectImpact:       plan.CollectImpact.ValueBool(),
		PathPrefix:          plan.PathPrefix.ValueString(),
		NotifyInSlack:       plan.NotifyInSlack.ValueBool(),
		IncludeInDashboard:  plan.IncludeInDashboard.ValueBool(),
		AutoTrackingDelay:   int(plan.AutoTrackingDelay.ValueInt64()),
		EnvironmentMappings: environmentMappings,
		BuildMappings:       buildMappingsT,
	}, nil
}

package sleuth

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/sleuth-io/terraform-provider-sleuth/internal/gqlclient"
)

var (
	_ resource.Resource                = &environmentResource{}
	_ resource.ResourceWithConfigure   = &environmentResource{}
	_ resource.ResourceWithImportState = &environmentResource{}
)

type envResourceModel struct {
	ID          types.String `tfsdk:"id"`
	ProjectSlug types.String `tfsdk:"project_slug"`
	Name        types.String `tfsdk:"name"`
	Slug        types.String `tfsdk:"slug"`
	Description types.String `tfsdk:"description"`
	Color       types.String `tfsdk:"color"`
}

type environmentResource struct {
	c *gqlclient.Client
}

func NewEnvironmentResource() resource.Resource {
	return &environmentResource{}
}

func (p *environmentResource) Schema(_ context.Context, _ resource.SchemaRequest, res *resource.SchemaResponse) {
	res.Schema = schema.Schema{
		Description: "Environment resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"project_slug": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The project for this environment",
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Environment name",
				Required:            true,
			},
			"slug": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Environment slug",
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Environment description",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"color": schema.StringAttribute{
				MarkdownDescription: "The color for the UI",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("#cecece"),
			},
		},
	}
}

func (p *environmentResource) Configure(ctx context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	p.c = req.ProviderData.(*gqlclient.Client)
}

func (p *environmentResource) Metadata(ctx context.Context, req resource.MetadataRequest, res *resource.MetadataResponse) {
	res.TypeName = req.ProviderTypeName + "_environment"
}

func (p *environmentResource) Create(ctx context.Context, req resource.CreateRequest, res *resource.CreateResponse) {
	ctx = tflog.SetField(ctx, "resource", "environment")
	ctx = tflog.SetField(ctx, "operation", "create")

	var plan envResourceModel
	diags := req.Plan.Get(ctx, &plan)
	res.Diagnostics.Append(diags...)

	tflog.Info(ctx, "Creating Environment resource", map[string]any{"plan": plan})

	inputFields := getMutableEnvStruct(plan)
	projectSlug := plan.ProjectSlug.ValueString()
	envName := plan.Name.ValueString()
	input := gqlclient.CreateEnvironmentMutationInput{ProjectSlug: projectSlug, MutableEnvironment: &inputFields}

	// We create the environment automatically when Project is created, so we need to check if it already exists
	existingEnv, err := p.c.GetEnvironmentByName(&projectSlug, &envName)

	if err != nil && err != gqlclient.ErrNotFound {
		res.Diagnostics.AddError("Error obtaining environment", fmt.Sprintf("Could not obtain environment, unexpected error: %+v", err.Error()))
		return
	}

	var env *gqlclient.Environment
	if existingEnv != nil {
		input := gqlclient.UpdateEnvironmentMutationInput{ProjectSlug: projectSlug, Slug: existingEnv.Slug, MutableEnvironment: &inputFields}
		env, err = p.c.UpdateEnvironment(input)
		if err != nil {
			tflog.Error(ctx, "Error updating Environment", map[string]any{"error": err.Error()})
			res.Diagnostics.AddError(
				"Error creating Environment",
				fmt.Sprintf("Could not create environment, unexpected error: %+v", err.Error()),
			)
			return
		}
	} else {
		env, err = p.c.CreateEnvironment(input)
		if err != nil {
			tflog.Error(ctx, "Error creating Environment", map[string]any{"error": err.Error()})
			res.Diagnostics.AddError(
				"Error creating Environment",
				fmt.Sprintf("Could not create environment, unexpected error: %+v", err.Error()),
			)
			return
		}
	}

	tflog.Info(ctx, "Created Environment", map[string]any{"environment": env})

	state := getNewStateFromEnv(env, projectSlug)

	diags = res.State.Set(ctx, state)
	res.Diagnostics.Append(diags...)
}

func (p *environmentResource) Read(ctx context.Context, req resource.ReadRequest, res *resource.ReadResponse) {
	ctx = tflog.SetField(ctx, "resource", "environment")
	ctx = tflog.SetField(ctx, "operation", "read")

	var state envResourceModel
	diags := req.State.Get(ctx, &state)

	projectSlug := state.ProjectSlug.ValueString()
	slug := state.Slug.ValueString()
	// when importing a resource, only ID will be set with project slug & slug
	if projectSlug == "" {
		id := state.ID.ValueString()
		splits := strings.Split(id, "/")
		projectSlug = splits[0]
		slug = splits[1]
	}

	tflog.Info(ctx, "Refreshing Environment resource", map[string]any{"state": fmt.Sprintf("%+v", state)})
	res.Diagnostics.Append(diags...)

	if res.Diagnostics.HasError() {
		return
	}

	env, err := p.c.GetEnvironment(&projectSlug, &slug)
	if err != nil {
		tflog.Error(ctx, fmt.Sprintf("Error obtaining environment: %+v", err))
		res.Diagnostics.AddError(
			"Error Reading Environment",
			fmt.Sprintf("Could not read Environment Slug %s, %+v", state.Slug.ValueString(), err.Error()),
		)
		return
	}
	if env == nil {
		return
	}

	newState := getNewStateFromEnv(env, projectSlug)

	diags = res.State.Set(ctx, &newState)
	res.Diagnostics.Append(diags...)

}

func (p *environmentResource) Update(ctx context.Context, req resource.UpdateRequest, res *resource.UpdateResponse) {
	ctx = tflog.SetField(ctx, "resource", "environment")
	ctx = tflog.SetField(ctx, "operation", "update")

	var plan envResourceModel
	diags := req.Plan.Get(ctx, &plan)
	res.Diagnostics.Append(diags...)

	var state envResourceModel
	diags = req.State.Get(ctx, &state)
	res.Diagnostics.Append(diags...)

	tflog.Info(ctx, "Updating Environment resource", map[string]any{"plan": plan, "state": state})

	inputFields := getMutableEnvStruct(plan)
	projectSlug := state.ProjectSlug.ValueString()
	slug := state.Slug.ValueString()
	input := gqlclient.UpdateEnvironmentMutationInput{
		ProjectSlug:        projectSlug,
		Slug:               slug,
		MutableEnvironment: &inputFields,
	}

	env, err := p.c.UpdateEnvironment(input)
	if err != nil {
		diags.AddError("Error updating environment", err.Error())
		return
	}

	tflog.Info(ctx, "Updated Environment", map[string]any{"environment": env})

	newState := getNewStateFromEnv(env, projectSlug)

	diags = res.State.Set(ctx, newState)
	res.Diagnostics.Append(diags...)
}

func (p *environmentResource) Delete(ctx context.Context, req resource.DeleteRequest, res *resource.DeleteResponse) {
	tflog.SetField(ctx, "resource", "environment")
	tflog.SetField(ctx, "operation", "delete")

	var state envResourceModel
	req.State.Get(ctx, &state)

	tflog.Info(ctx, "Deleting Environment resource", map[string]any{"state": fmt.Sprintf("%+v", state)})
	projectSlug := state.ProjectSlug.ValueStringPointer()
	err := p.c.DeleteEnvironment(ctx, projectSlug, state.Slug.ValueStringPointer())
	if err != nil {
		tflog.Error(ctx, "Unexpected error deleting environment", map[string]any{"error": err.Error()})
		res.Diagnostics.AddError("Unexpected error deleting environment", err.Error())
		return
	}

	res.State.RemoveResource(ctx)
}

func (p *environmentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, res *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, res)
}

func getNewStateFromEnv(env *gqlclient.Environment, projectSlug string) envResourceModel {
	return envResourceModel{
		ID:          types.StringValue(env.Slug),
		ProjectSlug: types.StringValue(projectSlug),
		Name:        types.StringValue(env.Name),
		Slug:        types.StringValue(env.Slug),
		Description: types.StringValue(env.Description),
		Color:       types.StringValue(env.Color),
	}
}

func getMutableEnvStruct(plan envResourceModel) gqlclient.MutableEnvironment {
	return gqlclient.MutableEnvironment{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
		Color:       plan.Color.ValueString(),
	}
}

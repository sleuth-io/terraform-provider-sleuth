package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sleuth-io/terraform-provider-sleuth/internal/gqlclient"
	"time"
)

func resourceProject() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "Sleuth project.",

		CreateContext: resourceProjectCreate,
		ReadContext:   resourceProjectRead,
		UpdateContext: resourceProjectUpdate,
		DeleteContext: resourceProjectDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "Project name",
				Type:        schema.TypeString,
				Required:    true,
			},
			"slug": {
				Description: "Project slug",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"description": {
				Description: "Project description",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
			},
			"issue_tracker_provider_type": {
				Description: "Where to find issues linked to by changes",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "SOURCE_PROVIDER",
			},
			"build_provider": {
				Description: "Where to find builds related to changes",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "NONE",
			},
			"change_failure_rate_boundary": {
				Description: "The health rating at which point it will be considered a failure",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "UNHEALTHY",
			},
			"impact_sensitivity": {
				Description: "How many impact measures Sleuth takes into account when auto-determining a deploys health.",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "NORMAL",
			},
			"failure_sensitivity": {
				Description: "The amount of time (in seconds) a deploy must spend in a failure status (Unhealthy, Incident, etc.) before it is determined a failure. Setting this value to a longer time means that less deploys will be classified.",
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     420,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceProjectCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*gqlclient.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	inputFields := gqlclient.MutableProject{}
	input := gqlclient.CreateProjectMutationInput{MutableProject: &inputFields}

	populateProjectInput(d, &inputFields)

	proj, err := c.CreateProject(input)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(proj.Slug)
	setProjectFields(d, proj)

	return diags
}

func resourceProjectUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*gqlclient.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	projectSlug := d.Id()

	slug := d.Get("slug").(string)
	inputFields := gqlclient.MutableProject{}
	input := gqlclient.UpdateProjectMutationInput{Slug: slug, MutableProject: &inputFields}

	populateProjectInput(d, &inputFields)

	proj, err := c.UpdateProject(&projectSlug, input)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("last_updated", time.Now().Format(time.RFC850))
	setProjectFields(d, proj)

	return diags
}

func resourceProjectRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*gqlclient.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	projectSlug := d.Id()

	proj, err := c.GetProject(&projectSlug)
	if err != nil {
		return diag.FromErr(err)
	}

	setProjectFields(d, proj)

	return diags

}

func setProjectFields(d *schema.ResourceData, proj *gqlclient.Project) {

	d.Set("name", proj.Name)
	d.Set("slug", proj.Slug)
	d.Set("description", proj.Description)
	d.Set("issue_tracker_provider_type", proj.IssueTrackerProvider)
	d.Set("build_provider", proj.BuildProvider)
	d.Set("change_failure_rate_boundary", proj.ChangeFailureRateBoundary)
	d.Set("impact_sensitivity", proj.ImpactSensitivity)
	d.Set("failure_sensitivity", proj.FailureSensitivity)
}

func populateProjectInput(d *schema.ResourceData, input *gqlclient.MutableProject) bool {
	input.Name = d.Get("name").(string)
	input.Description = d.Get("description").(string)
	input.IssueTrackerProvider = d.Get("issue_tracker_provider_type").(string)
	input.BuildProvider = d.Get("build_provider").(string)
	input.ChangeFailureRateBoundary = d.Get("change_failure_rate_boundary").(string)
	input.ImpactSensitivity = d.Get("impact_sensitivity").(string)
	input.FailureSensitivity = d.Get("failure_sensitivity").(int)
	return true
}

func resourceProjectDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*gqlclient.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	projectSlug := d.Id()

	err := c.DeleteProject(&projectSlug)
	if err != nil {
		// Ignore missing as the project gets deleted when the last env gets deleted
		if err.Error() != "Missing" {
			return diag.FromErr(err)
		}
	}

	d.SetId("")

	return diags
}

package provider

import (
	"context"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/sleuth-io/terraform-provider-sleuth/internal/gqlclient"
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
				Deprecated:  "Project description will be removed",
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
			"change_lead_time_start_definition": {
				Description: "The event that will be taken as a start definition (first commit, issue transition or whichever comes first) - options: COMMIT (default), ISSUE, FIRST_EVENT.",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "COMMIT",
			},
			"change_lead_time_issue_states": {
				Description: "Issue state IDs used for start definition (only used if change_lead_time_start_definition is ISSUE or FIRST_EVENT.",
				Type:        schema.TypeSet,
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			"change_lead_time_strict_matching": {
				Description: "When enabled Sleuth will only look for issue references in PR titles and PR branch names. If strict issue matching is disabled, Sleuth will expand the search for issue references to PR descriptions and commit messages.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
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
	setProjectFields(ctx, d, proj)

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
	setProjectFields(ctx, d, proj)

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

	if proj == nil {
		d.SetId("")
	} else {
		setProjectFields(ctx, d, proj)
	}

	return diags

}

func setProjectFields(ctx context.Context, d *schema.ResourceData, proj *gqlclient.Project) {
	d.Set("name", proj.Name)
	d.Set("slug", proj.Slug)
	d.Set("description", proj.Description)
	d.Set("issue_tracker_provider_type", proj.IssueTrackerProvider)
	d.Set("build_provider", proj.BuildProvider)
	d.Set("change_failure_rate_boundary", proj.ChangeFailureRateBoundary)
	d.Set("impact_sensitivity", proj.ImpactSensitivity)
	d.Set("failure_sensitivity", proj.FailureSensitivity)
	d.Set("change_lead_time_start_definition", proj.CltStartDefinition)
	d.Set("change_lead_time_strict_matching", proj.StrictIssueMatching)

	cltStateInts := []int{}
	for _, val := range proj.CltStartStates {
		x, err := strconv.Atoi(val.ID)
		if err != nil {
			tflog.Error(ctx, "Error converting ID to int")
			continue
		}
		cltStateInts = append(cltStateInts, x)
	}
	d.Set("change_lead_time_issue_states", cltStateInts)
}

func populateProjectInput(d *schema.ResourceData, input *gqlclient.MutableProject) bool {
	input.Name = d.Get("name").(string)
	input.Description = d.Get("description").(string)
	input.IssueTrackerProvider = d.Get("issue_tracker_provider_type").(string)
	input.BuildProvider = d.Get("build_provider").(string)
	input.ChangeFailureRateBoundary = d.Get("change_failure_rate_boundary").(string)
	input.ImpactSensitivity = d.Get("impact_sensitivity").(string)
	input.FailureSensitivity = d.Get("failure_sensitivity").(int)
	input.CltStartDefinition = d.Get("change_lead_time_start_definition").(string)

	input.StrictIssueMatching = d.Get("change_lead_time_strict_matching").(bool)

	cltStates := d.Get("change_lead_time_issue_states").(*schema.Set)
	cltStatesList := cltStates.List()
	cltStatesInts := make([]int, len(cltStatesList))
	for i := range cltStatesList {
		cltStatesInts[i] = cltStatesList[i].(int)
	}
	input.CltStartStates = cltStatesInts

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

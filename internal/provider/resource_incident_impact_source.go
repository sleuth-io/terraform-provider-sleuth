package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sleuth-io/terraform-provider-sleuth/internal/gqlclient"
	"strings"
	"time"
)

func resourceIncidentImpactSource() *schema.Resource {
	return &schema.Resource{
		Description: "Sleuth incident impact source",

		CreateContext: resourceIncidentImpactSourceCreate,
		ReadContext:   resourceIncidentImpactSourceRead,
		UpdateContext: resourceIncidentImpactSourceUpdate,
		DeleteContext: resourceIncidentImpactSourceDelete,

		Schema: map[string]*schema.Schema{
			"project_slug": {
				Description: "Project slug",
				Type:        schema.TypeString,
				Required:    true,
			},
			"name": {
				Description: "Impact source name",
				Type:        schema.TypeString,
				Required:    true,
			},
			// can't use `provider` because terraform tries to import provider
			"provider_name": {
				Description: "Impact source provider (options: PAGERDUTY)",
				Type:        schema.TypeString,
				Required:    true,
			},
			"environment_name": {
				Description: "Impact source environment name",
				Type:        schema.TypeString,
				Required:    true,
			},
			"pagerduty_input": {
				Description: "PagerDuty input",
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"remote_services": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "List of remote services, empty string means all",
							Default:     "",
						},
						"remote_urgency": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "PagerDuty remote urgency, options: HIGH, LOW, ANY",
						},
						"historic_init": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Populate with data from the last 30 days",
						},
						"integration_slug": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "IntegrationAuthentication slug used",
						},
					},
				},
			},
			"slug": {
				Description: "Impact source slug",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceIncidentImpactSourceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*gqlclient.Client)

	projectSlug, slug, err := getSlugsFromID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	tflog.Debug(ctx, fmt.Sprintf("Project slug, incident slug %s, %s", projectSlug, slug))
	iis, err := c.GetIncidentImpactSource(ctx, projectSlug, slug)
	if err != nil {
		return diag.FromErr(err)
	}

	setFields(ctx, d, iis, projectSlug)

	return nil
}

func resourceIncidentImpactSourceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*gqlclient.Client)

	projectSlug := d.Get("project_slug").(string)
	createIncidentImpactSourceMutationInput := gqlclient.IncidentImpactSourceInputType{
		ProjectSlug:     projectSlug,
		Name:            d.Get("name").(string),
		Provider:        d.Get("provider_name").(string),
		EnvironmentName: strings.ToLower(d.Get("environment_name").(string)),
		PagerDutyInputType: gqlclient.PagerDutyInputType{
			RemoteServices: d.Get("pagerduty_input.0.remote_services").(string),
			RemoteUrgency:  d.Get("pagerduty_input.0.remote_urgency").(string),
		},
	}

	incidentImpact, err := c.CreateIncidentImpactSource(ctx, createIncidentImpactSourceMutationInput)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(fmt.Sprintf("%s/%s", projectSlug, incidentImpact.Slug))

	setFields(ctx, d, incidentImpact, projectSlug)

	return nil
}

func resourceIncidentImpactSourceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*gqlclient.Client)

	projectSlug, slug, err := getSlugsFromID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	input := gqlclient.IncidentImpactSourceInputUpdateType{
		Slug: slug,
		IncidentImpactSourceInputType: gqlclient.IncidentImpactSourceInputType{
			ProjectSlug:     projectSlug,
			Name:            d.Get("name").(string),
			Provider:        d.Get("provider_name").(string),
			EnvironmentName: strings.ToLower(d.Get("environment_name").(string)),
			PagerDutyInputType: gqlclient.PagerDutyInputType{
				RemoteServices: d.Get("pagerduty_input.0.remote_services").(string),
				RemoteUrgency:  d.Get("pagerduty_input.0.remote_urgency").(string),
			},
		},
	}

	proj, err := c.UpdateIncidentImpactSource(ctx, input)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("last_updated", time.Now().Format(time.RFC850))

	setFields(ctx, d, proj, projectSlug)

	return nil
}

func resourceIncidentImpactSourceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*gqlclient.Client)

	projectSlug, slug, err := getSlugsFromID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	input := gqlclient.IncidentImpactSourceDeleteInputType{Slug: slug, ProjectSlug: projectSlug}

	succ, err := c.DeleteIncidentImpactSource(ctx, input)
	if err != nil {
		return diag.FromErr(err)
	}
	if !succ {
		return diag.FromErr(fmt.Errorf("unsuccessful deletion of incident impact source"))
	}

	return nil
}

func getSlugsFromID(id string) (string, string, error) {
	splits := strings.Split(id, "/")
	if len(splits) != 2 {
		return "", "", fmt.Errorf("invalid resource ID: %s", id)
	}

	return splits[0], splits[1], nil
}

func setFields(ctx context.Context, d *schema.ResourceData, is *gqlclient.IncidentImpactSource, projectSlug string) {
	d.Set("name", is.Name)
	d.Set("slug", is.Slug)
	d.Set("provider_name", strings.ToUpper(is.Provider))
	d.Set("environment_name", is.Environment.Name)
	d.Set("project_slug", projectSlug)

	pager_duty_input := make(map[string]interface{})
	pager_duty_input["remote_services"] = is.ProviderData.PagerDutyProviderData.RemoteServices
	pager_duty_input["remote_urgency"] = is.ProviderData.PagerDutyProviderData.RemoteUrgency

	d.Set("pagerduty_input", []map[string]interface{}{pager_duty_input})

}

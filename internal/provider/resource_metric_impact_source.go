package provider

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/sleuth-io/terraform-provider-sleuth/internal/gqlclient"
)

func resourceMetricImpactSource() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "Sleuth error impact source.",

		CreateContext: resourceMetricImpactSourceCreate,
		ReadContext:   resourceMetricImpactSourceRead,
		UpdateContext: resourceMetricImpactSourceUpdate,
		DeleteContext: resourceMetricImpactSourceDelete,

		Schema: map[string]*schema.Schema{
			"project_slug": {
				Description: "The project for this environment",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"environment_slug": {
				Description: "The environment slug",
				Type:        schema.TypeString,
				Required:    true,
			},
			"name": {
				Description: "Impact source name",
				Type:        schema.TypeString,
				Required:    true,
			},
			"provider_type": {
				Description: "Integration provider type",
				Type:        schema.TypeString,
				Required:    true,
			},
			"integration_slug": {
				Description: "The integration slug",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"query": {
				Description: "The metric query",
				Type:        schema.TypeString,
				Required:    true,
			},
			"less_is_better": {
				Description: "Whether smaller values are better or not",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"manually_set_health_threshold": {
				Description: "The manually set threshold to start marking failed values",
				Type:        schema.TypeFloat,
				Optional:    true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceMetricImpactSourceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*gqlclient.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	projectSlug := d.Get("project_slug").(string)
	inputFields := gqlclient.MutableMetricImpactSource{}
	input := gqlclient.CreateMetricImpactSourceMutationInput{ProjectSlug: projectSlug, MutableMetricImpactSource: &inputFields}

	populateMetricImpactSource(d, &inputFields)

	src, err := c.CreateMetricImpactSource(input)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(fmt.Sprintf("%s/%s", projectSlug, src.Slug))
	setMetricImpactSourceFields(d, projectSlug, src)

	return diags
}

func resourceMetricImpactSourceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*gqlclient.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	parsed := strings.Split(d.Id(), "/")
	projectSlug := parsed[0]
	impactSourceSlug := parsed[1]

	inputFields := gqlclient.MutableMetricImpactSource{}
	input := gqlclient.UpdateMetricImpactSourceMutationInput{ProjectSlug: projectSlug, Slug: impactSourceSlug, MutableMetricImpactSource: &inputFields}
	populateMetricImpactSource(d, &inputFields)

	proj, err := c.UpdateMetricImpactSource(input)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("last_updated", time.Now().Format(time.RFC850))
	setMetricImpactSourceFields(d, projectSlug, proj)

	return diags
}

func resourceMetricImpactSourceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*gqlclient.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	parsed := strings.Split(d.Id(), "/")
	projectSlug := parsed[0]
	environmentSlug := parsed[1]

	source, err := c.GetMetricImpactSource(&projectSlug, &environmentSlug)
	if err != nil {
		return diag.FromErr(err)
	} else if source == nil {
		d.SetId("")
	} else {
		setMetricImpactSourceFields(d, projectSlug, source)
	}

	return diags

}

func setMetricImpactSourceFields(d *schema.ResourceData, projectSlug string, env *gqlclient.MetricImpactSource) {

	d.Set("project_slug", projectSlug)
	d.Set("name", env.Name)
	d.Set("environment_slug", env.Environment.Slug)
	d.Set("provider_type", env.Provider)
	d.Set("query", env.Query)
	d.Set("less_is_better", env.LessIsBetter)
	d.Set("manually_set_health_threshold", env.ManuallySetHealthThreshold)
}

func populateMetricImpactSource(d *schema.ResourceData, input *gqlclient.MutableMetricImpactSource) bool {
	input.Name = d.Get("name").(string)
	var envRaw = d.Get("environment_slug").(string)

	var envSlug string
	if strings.Contains(envRaw, "/") {
		envSlug = strings.Split(envRaw, "/")[1]
	} else {
		envSlug = envRaw
	}

	var providerType = d.Get("provider_type").(string)
	var integrationSlug = d.Get("integration_slug").(string)
	if integrationSlug == "" {
		integrationSlug = providerType
	}
	input.EnvironmentSlug = envSlug
	input.Provider = strings.ToUpper(providerType)
	input.Query = d.Get("query").(string)
	input.IntegrationSlug = integrationSlug
	input.LessIsBetter = d.Get("less_is_better").(bool)
	input.ManuallySetHealthThreshold = d.Get("manually_set_health_threshold").(float64)

	return true
}

func resourceMetricImpactSourceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*gqlclient.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	parsed := strings.Split(d.Id(), "/")
	projectSlug := parsed[0]
	environmentSlug := parsed[1]

	err := c.DeleteImpactSource(&projectSlug, &environmentSlug)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}

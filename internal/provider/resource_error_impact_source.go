package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sleuth-io/terraform-provider-sleuth/internal/gqlclient"
	"strings"
	"time"
)

func resourceErrorImpactSource() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "Sleuth error impact source.",

		CreateContext: resourceErrorImpactSourceCreate,
		ReadContext:   resourceErrorImpactSourceRead,
		UpdateContext: resourceErrorImpactSourceUpdate,
		DeleteContext: resourceErrorImpactSourceDelete,

		Schema: map[string]*schema.Schema{
			"project_slug": {
				Description: "The project for this impact source",
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
			"error_org_key": {
				Description: "The organization key of the integration provider",
				Type:        schema.TypeString,
				Required:    true,
			},
			"error_project_key": {
				Description: "The project key of the integration provider",
				Type:        schema.TypeString,
				Required:    true,
			},
			"error_environment": {
				Description: "The environment of the integration provider",
				Type:        schema.TypeString,
				Required:    true,
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

func resourceErrorImpactSourceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*gqlclient.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	projectSlug := d.Get("project_slug").(string)
	inputFields := gqlclient.MutableErrorImpactSource{}
	input := gqlclient.CreateErrorImpactSourceMutationInput{ProjectSlug: projectSlug, MutableErrorImpactSource: &inputFields}

	populateErrorImpactSource(d, &inputFields)

	src, err := c.CreateErrorImpactSource(input)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(fmt.Sprintf("%s/%s", projectSlug, src.Slug))
	setErrorImpactSourceFields(d, projectSlug, src)

	return diags
}

func resourceErrorImpactSourceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*gqlclient.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	parsed := strings.Split(d.Id(), "/")
	projectSlug := parsed[0]
	impactSourceSlug := parsed[1]

	inputFields := gqlclient.MutableErrorImpactSource{}
	input := gqlclient.UpdateErrorImpactSourceMutationInput{ProjectSlug: projectSlug, Slug: impactSourceSlug, MutableErrorImpactSource: &inputFields}
	populateErrorImpactSource(d, &inputFields)

	proj, err := c.UpdateErrorImpactSource(input)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("last_updated", time.Now().Format(time.RFC850))
	setErrorImpactSourceFields(d, projectSlug, proj)

	return diags
}

func resourceErrorImpactSourceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*gqlclient.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	parsed := strings.Split(d.Id(), "/")
	projectSlug := parsed[0]
	environmentSlug := parsed[1]

	source, err := c.GetErrorImpactSource(&projectSlug, &environmentSlug)
	if err != nil {
		return diag.FromErr(err)
	} else if source == nil {
		d.SetId("")
	} else {
		setErrorImpactSourceFields(d, projectSlug, source)
	}

	return diags

}

func setErrorImpactSourceFields(d *schema.ResourceData, projectSlug string, env *gqlclient.ErrorImpactSource) {

	d.Set("project_slug", projectSlug)
	d.Set("name", env.Name)
	d.Set("environment_slug", fmt.Sprintf("%s/%s", projectSlug, env.Environment.Slug))
	d.Set("provider_type", env.Provider)
	d.Set("error_org_key", env.ErrorOrgKey)
	d.Set("error_project_key", env.ErrorProjectKey)
	d.Set("error_environment", env.ErrorEnvironment)
	d.Set("manually_set_health_threshold", env.ManuallySetHealthThreshold)
}

func populateErrorImpactSource(d *schema.ResourceData, input *gqlclient.MutableErrorImpactSource) bool {
	input.Name = d.Get("name").(string)
	var envRaw = d.Get("environment_slug").(string)
	var envSlug = strings.Split(envRaw, "/")[1]
	input.EnvironmentSlug = envSlug
	input.Provider = strings.ToUpper(d.Get("provider_type").(string))
	input.ErrorOrgKey = d.Get("error_org_key").(string)
	input.ErrorProjectKey = d.Get("error_project_key").(string)
	input.ErrorEnvironment = d.Get("error_environment").(string)
	input.ManuallySetHealthThreshold = d.Get("manually_set_health_threshold").(float64)

	return true
}

func resourceErrorImpactSourceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

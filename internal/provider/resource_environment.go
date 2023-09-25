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

func resourceEnvironment() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "Sleuth environment.",

		CreateContext: resourceEnvironmentCreate,
		ReadContext:   resourceEnvironmentRead,
		UpdateContext: resourceEnvironmentUpdate,
		DeleteContext: resourceEnvironmentDelete,

		Schema: map[string]*schema.Schema{
			"project_slug": {
				Description: "The project for this environment",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"name": {
				Description: "Environment name",
				Type:        schema.TypeString,
				Required:    true,
			},
			"slug": {
				Description: "Environment slug",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"description": {
				Description: "Environment description",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
			},
			"color": {
				Description: "The color for the UI",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "#cecece",
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceEnvironmentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*gqlclient.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	projectSlug := d.Get("project_slug").(string)
	name := d.Get("name").(string)

	existingEnv, _ := c.GetEnvironmentByName(&projectSlug, &name)
	if existingEnv != nil {
		d.SetId(fmt.Sprintf("%s/%s", projectSlug, existingEnv.Slug))
		resourceEnvironmentUpdate(ctx, d, meta)
	} else {
		inputFields := gqlclient.MutableEnvironment{}
		input := gqlclient.CreateEnvironmentMutationInput{ProjectSlug: projectSlug, MutableEnvironment: &inputFields}

		populateInput(d, &inputFields)

		env, err := c.CreateEnvironment(input)
		if err != nil {
			return diag.FromErr(err)
		}
		d.SetId(fmt.Sprintf("%s/%s", projectSlug, env.Slug))
		setEnvironmentFields(d, projectSlug, env)
	}

	return diags
}

func resourceEnvironmentUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*gqlclient.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	parsed := strings.Split(d.Id(), "/")
	projectSlug := parsed[0]
	environmentSlug := parsed[1]

	inputFields := gqlclient.MutableEnvironment{}
	input := gqlclient.UpdateEnvironmentMutationInput{ProjectSlug: projectSlug, Slug: environmentSlug, MutableEnvironment: &inputFields}
	populateInput(d, &inputFields)

	proj, err := c.UpdateEnvironment(input)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("last_updated", time.Now().Format(time.RFC850))
	setEnvironmentFields(d, projectSlug, proj)

	return diags
}

func resourceEnvironmentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*gqlclient.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	parsed := strings.Split(d.Id(), "/")
	projectSlug := parsed[0]
	environmentSlug := parsed[1]

	env, err := c.GetEnvironment(&projectSlug, &environmentSlug)
	if err != nil {
		return diag.FromErr(err)
	} else if env == nil {
		d.SetId("")
	} else {
		setEnvironmentFields(d, projectSlug, env)
	}

	return diags

}

func setEnvironmentFields(d *schema.ResourceData, projectSlug string, env *gqlclient.Environment) {

	d.Set("project_slug", projectSlug)
	d.Set("name", env.Name)
	d.Set("slug", env.Slug)
	d.Set("description", env.Description)
	d.Set("color", env.Color)
}

func populateInput(d *schema.ResourceData, input *gqlclient.MutableEnvironment) bool {
	input.Name = d.Get("name").(string)
	input.Description = d.Get("description").(string)
	input.Color = d.Get("color").(string)
	return true
}

func resourceEnvironmentDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*gqlclient.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	parsed := strings.Split(d.Id(), "/")
	projectSlug := parsed[0]
	environmentSlug := parsed[1]

	err := c.DeleteEnvironment(ctx, &projectSlug, &environmentSlug)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}

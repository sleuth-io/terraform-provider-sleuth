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
		Description: "Sample resource in the Terraform provider scaffolding.",

		CreateContext: resourceProjectCreate,
		ReadContext:   resourceProjectRead,
		UpdateContext: resourceProjectUpdate,
		DeleteContext: resourceProjectDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				// This description is used by the documentation generator and the language server.
				Description: "Project name",
				Type:        schema.TypeString,
				Required:    true,
			},
			"failure_sensitivity": {
				// This description is used by the documentation generator and the language server.
				Description: "Failure sensitivity",
				Type:        schema.TypeInt,
				Optional:    true,
			},
		},
	}
}

func resourceProjectCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*gqlclient.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	name := d.Get("name").(string)

	input := gqlclient.ProjectCreationMutationInput{Name: name}

	val, ok := d.GetOk("failureSensitivity")
	if ok {
		input.FailureSensitivity = val.(string)
	}

	proj, err := c.CreateProject( input)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(proj.Slug)

	resourceProjectRead(ctx, d, meta)

	return diags
}

func resourceProjectRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*gqlclient.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	projectSlug := d.Id()

	_, err := c.GetProject(&projectSlug)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags

}

func resourceProjectUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*gqlclient.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	projectSlug := d.Id()

	input := gqlclient.ProjectUpdateMutationInput{}
	changed := false

	val, ok := d.GetOk("failure_sensitivity")
	if ok {
		input.FailureSensitivity = val.(int)
		changed = true

	}

	if d.HasChange("name") {
		name := d.Get("name").(string)
		input.Name = name
		changed = true
	}

	if changed {
		_, err := c.UpdateProject(&projectSlug, input)
		if err != nil {
			return diag.FromErr(err)
		}
		d.Set("last_updated", time.Now().Format(time.RFC850))
	}

	return diags
}

func resourceProjectDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*gqlclient.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	projectSlug := d.Id()

	err := c.DeleteProject( &projectSlug)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/sleuth-io/terraform-provider-sleuth/internal/gqlclient"
)

func init() {
	// Set descriptions to support markdown syntax, this will be used in document generation
	// and the language server.
	schema.DescriptionKind = schema.StringMarkdown

	// Customize the content of descriptions when output. For example you can add defaults on
	// to the exported descriptions if present.
	// schema.SchemaDescriptionBuilder = func(s *schema.Schema) string {
	// 	desc := s.Description
	// 	if s.Default != nil {
	// 		desc += fmt.Sprintf(" Defaults to `%v`.", s.Default)
	// 	}
	// 	return strings.TrimSpace(desc)
	// }
}

func New(version string) func() *schema.Provider {
	return func() *schema.Provider {
		p := &schema.Provider{
			Schema: map[string]*schema.Schema{
				"baseurl": &schema.Schema{
					Type:        schema.TypeString,
					Description: "Ignore this, as it is only used by Sleuth developers",
					Optional:    true,

					DefaultFunc: schema.EnvDefaultFunc("SLEUTH_BASEURL", "https://app.sleuth.io"),
				},
				"api_key": &schema.Schema{
					Type:        schema.TypeString,
					Description: "The Sleuth organization's Api key",
					Required:    true,
					DefaultFunc: schema.EnvDefaultFunc("SLEUTH_API_KEY", nil),
				},
			},
			//DataSourcesMap: map[string]*schema.Resource{
			//	"scaffolding_data_source": dataSourceScaffolding(),
			//},
			ResourcesMap: map[string]*schema.Resource{
				"sleuth_project":                resourceProject(),
				"sleuth_environment":            resourceEnvironment(),
				"sleuth_error_impact_source":    resourceErrorImpactSource(),
				"sleuth_metric_impact_source":   resourceMetricImpactSource(),
				"sleuth_code_change_source":     resourceCodeChangeSource(),
				"sleuth_incident_impact_source": resourceIncidentImpactSource(),
			},
		}

		p.ConfigureContextFunc = configure(version, p)

		return p
	}
}

func configure(version string, p *schema.Provider) func(context.Context, *schema.ResourceData) (interface{}, diag.Diagnostics) {
	return func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {

		apiKey := d.Get("api_key").(string)

		var baseurl *string

		hVal, ok := d.GetOk("baseurl")
		if ok {
			tempBaseurl := hVal.(string)
			baseurl = &tempBaseurl
		}

		// Warning or errors can be collected in a slice type
		var diags diag.Diagnostics

		ua := p.UserAgent("sleuth_terraform_provider", version)
		c, err := gqlclient.NewClient(baseurl, &apiKey, ua)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to create Sleuth client",
				Detail:   "Unable to authenticate api key for authenticated Sleuth client",
			})

			return nil, diags
		}

		return c, diags
	}
}

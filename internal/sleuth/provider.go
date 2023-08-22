package sleuth

import (
	"context"
	"fmt"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/meta"

	"github.com/sleuth-io/terraform-provider-sleuth/internal/gqlclient"
)

// Ensure the implementation satisfies the expected interfaces
var (
	_ provider.Provider = &sleuthProvider{}
)

// New is a helper function to simplify provider server and testing implementation.
func New(v string) provider.Provider {
	return &sleuthProvider{v: v}
}

// sleuthProvider is the provider implementation.
type sleuthProvider struct {
	v string
}

// sleuthProviderModel maps provider schema data to a Go type.
type sleuthProviderModel struct {
	APIKey  types.String `tfsdk:"api_key"`
	BaseURL types.String `tfsdk:"baseurl"`
}

// Metadata returns the provider type name.
func (p *sleuthProvider) Metadata(_ context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "sleuth"
	resp.Version = p.v
}

// Schema defines the provider-level schema for configuration data.
func (p *sleuthProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "",
		Attributes: map[string]schema.Attribute{
			"api_key": schema.StringAttribute{
				MarkdownDescription: "The Sleuth organization's Api key",
				Optional:            true,
			},
			"baseurl": schema.StringAttribute{
				MarkdownDescription: "Ignore this, as it is only used by Sleuth developers",
				Optional:            true,
			},
		},
	}
}

func (p *sleuthProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring Sleuth client")

	apiKeyFallback := os.Getenv("SLEUTH_API_KEY")
	baseURLENVFallback := os.Getenv("SLEUTH_BASEURL")

	// Retrieve provider data from configuration
	var config sleuthProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.APIKey.IsNull() && apiKeyFallback == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key"),
			"API key must be set",
			"API key must be set",
		)
		return
	}

	apiKey := config.APIKey
	if config.APIKey.IsNull() {
		apiKey = types.StringValue(apiKeyFallback)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	baseURL := config.BaseURL
	if baseURL.IsNull() {
		if baseURLENVFallback != "" {
			baseURL = types.StringValue(baseURLENVFallback)
		} else {
			baseURL = types.StringValue("https://app.sleuth.io")
		}
	}

	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "sleuth_base_url", baseURL)
	ua := userAgent(req.TerraformVersion, "terraform-provider-sleuth", p.v)
	c, err := gqlclient.NewClient(baseURL.ValueStringPointer(), apiKey.ValueStringPointer(), ua)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating new client",
			fmt.Sprintf("%+v", err),
		)
		return
	}

	// Make the client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = c
	resp.ResourceData = c

	tflog.Info(ctx, "Configured Sleuth client", map[string]any{"success": true})
}

// DataSources defines the data sources implemented in the provider.
func (p *sleuthProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

// Resources defines the resources implemented in the provider.
func (p *sleuthProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewProjectResource,
		NewCodeChangeSourceResource,
		NewEnvironmentResource,
		NewMetricImpactSourceResource,
	}
}

// modified from Plugin SDK (https://github.com/hashicorp/terraform-plugin-sdk/blob/ee14c4b6cb40fe4c6dc8ad2e50eda4c7f29cd291/helper/schema/provider.go#L489)
func userAgent(terraformVersion, name, version string) string {
	ua := fmt.Sprintf("Terraform/%s (+https://www.terraform.io) Terraform-Plugin-SDK/%s", terraformVersion, meta.SDKVersionString())
	if name != "" {
		ua += " " + name
		if version != "" {
			ua += "/" + version
		}
	}

	return ua
}

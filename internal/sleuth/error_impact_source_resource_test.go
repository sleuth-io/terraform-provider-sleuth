package sleuth

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAccErrorImpactSourceResource_v6(t *testing.T) {
	// tests are run in parallel both locally & on CI, so we need to generate a random name so slugs don't collide
	randomStr := acctest.RandStringFromCharSet(5, acctest.CharSetAlpha)
	projName := fmt.Sprintf("Terraform test project %s", randomStr)
	slug := fmt.Sprintf("terraform-test-project-%s", randomStr)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			//v6 will fail for v0.14.0
			tfversion.SkipBetween(tfversion.Version0_14_0, tfversion.Version0_15_0),
		},
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: createErrorImpactSourceConfig(projName),

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sleuth_project.terraform_acc_test", "name", projName),
					resource.TestCheckResourceAttr("sleuth_project.terraform_acc_test", "slug", slug),

					resource.TestCheckResourceAttr("sleuth_environment.terraform_acc_test", "name", "Terraform test environment"),

					resource.TestCheckResourceAttr("sleuth_error_impact_source.sentry_terraform_acc_test", "name", "Sentry errors"),
					resource.TestCheckResourceAttr("sleuth_error_impact_source.sentry_terraform_acc_test", "provider_type", "SENTRY"),
					resource.TestCheckResourceAttr("sleuth_error_impact_source.sentry_terraform_acc_test", "error_org_key", "sleuthio"),
					resource.TestCheckResourceAttr("sleuth_error_impact_source.sentry_terraform_acc_test", "error_project_key", "sleuth-dev"),
					resource.TestCheckResourceAttr("sleuth_error_impact_source.sentry_terraform_acc_test", "error_environment", "prod"),
					resource.TestCheckNoResourceAttr("sleuth_error_impact_source.sentry_terraform_acc_test", "manually_set_health_threshold"),
				),
			},
			// Update testing
			{
				Config: updateErrorImpactSourceConfig(projName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sleuth_project.terraform_acc_test", "name", projName),
					resource.TestCheckResourceAttr("sleuth_project.terraform_acc_test", "slug", slug),

					resource.TestCheckResourceAttr("sleuth_environment.terraform_acc_test", "name", "Terraform test environment"),

					resource.TestCheckResourceAttr("sleuth_error_impact_source.sentry_terraform_acc_test", "name", "Sentry errors update"),
					resource.TestCheckResourceAttr("sleuth_error_impact_source.sentry_terraform_acc_test", "provider_type", "SENTRY"),
					resource.TestCheckResourceAttr("sleuth_error_impact_source.sentry_terraform_acc_test", "error_org_key", "sleuthio"),
					resource.TestCheckResourceAttr("sleuth_error_impact_source.sentry_terraform_acc_test", "error_project_key", "sleuth-dev"),
					resource.TestCheckResourceAttr("sleuth_error_impact_source.sentry_terraform_acc_test", "error_environment", "staging"),
					resource.TestCheckResourceAttr("sleuth_error_impact_source.sentry_terraform_acc_test", "manually_set_health_threshold", "5"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func createErrorImpactSourceConfig(name string) string {
	return fmt.Sprintf(`
resource "sleuth_project" "terraform_acc_test" {
	name = "%s"
}

resource "sleuth_environment" "terraform_acc_test" {
	project_slug = sleuth_project.terraform_acc_test.slug
	name = "Terraform test environment"
}

resource "sleuth_error_impact_source" "sentry_terraform_acc_test" {
	project_slug = sleuth_project.terraform_acc_test.slug
	environment_slug = sleuth_environment.terraform_acc_test.slug
	name = "Sentry errors"
	provider_type = "SENTRY"
	error_org_key = "sleuthio"
	error_project_key = "sleuth-dev"
	error_environment = "prod"
	integration_slug = "sentry"
}
`, name)
}

func updateErrorImpactSourceConfig(name string) string {
	return fmt.Sprintf(`
resource "sleuth_project" "terraform_acc_test" {
	name = "%s"
}

resource "sleuth_environment" "terraform_acc_test" {
	project_slug = sleuth_project.terraform_acc_test.slug
	name = "Terraform test environment"
}

resource "sleuth_error_impact_source" "sentry_terraform_acc_test" {
  	project_slug = sleuth_project.terraform_acc_test.slug
  	environment_slug = sleuth_environment.terraform_acc_test.slug
  	name = "Sentry errors update"
  	provider_type = "SENTRY"
  	error_org_key = "sleuthio"
  	error_project_key = "sleuth-dev"
  	error_environment = "staging"
	manually_set_health_threshold = 5.0
	integration_slug = "sentry"
}`, name)
}

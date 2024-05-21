package sleuth

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAccProjectResource_v6(t *testing.T) {
	// tests are run in parallel both locally & on CI, so we need to generate a random name so slugs don't collide
	randomStr := acctest.RandStringFromCharSet(5, acctest.CharSetAlpha)
	name := fmt.Sprintf("Terraform test project %s", randomStr)
	updatedName := fmt.Sprintf("Terraform test project updated %s", randomStr)
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
				Config: createConfig(name),

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sleuth_project.terraform_acc_test", "name", name),
					resource.TestCheckResourceAttr("sleuth_project.terraform_acc_test", "slug", slug),
					resource.TestCheckResourceAttr("sleuth_project.terraform_acc_test", "issue_tracker_provider_type", "SOURCE_PROVIDER"),
					resource.TestCheckResourceAttr("sleuth_project.terraform_acc_test", "build_provider", "NONE"),
					resource.TestCheckResourceAttr("sleuth_project.terraform_acc_test", "change_failure_rate_boundary", "UNHEALTHY"),
					resource.TestCheckResourceAttr("sleuth_project.terraform_acc_test", "impact_sensitivity", "NORMAL"),
					resource.TestCheckResourceAttr("sleuth_project.terraform_acc_test", "failure_sensitivity", "420"),
				),
			},
			// Update testing
			{
				Config: updateConfig(updatedName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sleuth_project.terraform_acc_test", "name", updatedName),
					resource.TestCheckResourceAttr("sleuth_project.terraform_acc_test", "slug", slug),
					resource.TestCheckResourceAttr("sleuth_project.terraform_acc_test", "issue_tracker_provider_type", "SOURCE_PROVIDER"),
					resource.TestCheckResourceAttr("sleuth_project.terraform_acc_test", "build_provider", "GITHUB"),
					resource.TestCheckResourceAttr("sleuth_project.terraform_acc_test", "change_failure_rate_boundary", "INCIDENT"),
					resource.TestCheckResourceAttr("sleuth_project.terraform_acc_test", "impact_sensitivity", "FINE"),
					resource.TestCheckResourceAttr("sleuth_project.terraform_acc_test", "failure_sensitivity", "200"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func createConfig(name string) string {
	return fmt.Sprintf(`
resource "sleuth_project" "terraform_acc_test" {
  name = "%s"
}
`, name)
}

func updateConfig(name string) string {
	return fmt.Sprintf(`
resource "sleuth_project" "terraform_acc_test" {
  name = "%s"
  build_provider = "GITHUB"
  change_failure_rate_boundary = "INCIDENT"
  impact_sensitivity = "FINE"
  failure_sensitivity = 200
}`, name)
}

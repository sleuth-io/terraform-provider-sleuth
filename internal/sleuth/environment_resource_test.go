package sleuth

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAccEnvironmentResource_v6(t *testing.T) {
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
				Config: createEnvConfig(name),

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sleuth_project.terraform_acc_test", "name", name),
					resource.TestCheckResourceAttr("sleuth_project.terraform_acc_test", "slug", slug),

					resource.TestCheckResourceAttr("sleuth_environment.terraform_acc_test", "name", "staging"),
					resource.TestCheckResourceAttr("sleuth_environment.terraform_acc_test", "project_slug", slug),
					resource.TestCheckResourceAttr("sleuth_environment.terraform_acc_test", "description", "description abc"),
					resource.TestCheckResourceAttr("sleuth_environment.terraform_acc_test", "color", "#cecece"),
				),
			},
			// Update testing
			{
				Config: updateEnvConfig(updatedName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sleuth_project.terraform_acc_test", "name", updatedName),
					resource.TestCheckResourceAttr("sleuth_project.terraform_acc_test", "slug", slug),

					resource.TestCheckResourceAttr("sleuth_environment.terraform_acc_test", "name", "staging updated"),
					resource.TestCheckResourceAttr("sleuth_environment.terraform_acc_test", "project_slug", slug),
					resource.TestCheckResourceAttr("sleuth_environment.terraform_acc_test", "description", "description updated"),
					resource.TestCheckResourceAttr("sleuth_environment.terraform_acc_test", "color", "#ffffff"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

// Production name is omitted on purpose because of a bug with removing the env before the project
func createEnvConfig(name string) string {
	return fmt.Sprintf(`
resource "sleuth_project" "terraform_acc_test" {
	name = "%s"
}

resource "sleuth_environment" "terraform_acc_test" {
	project_slug = sleuth_project.terraform_acc_test.slug
	name = "staging"
	description = "description abc"
}
`, name)
}

func updateEnvConfig(name string) string {
	return fmt.Sprintf(`
resource "sleuth_project" "terraform_acc_test" {
	name = "%s"
}

resource "sleuth_environment" "terraform_acc_test" {
	project_slug = sleuth_project.terraform_acc_test.slug
	name = "staging updated"
	description = "description updated"
	color = "#ffffff"
}
`, name)
}

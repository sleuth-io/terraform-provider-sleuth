package sleuth

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAccChangeSourceResource_v6(t *testing.T) {
	// tests are run in parallel both locally & on CI, so we need to generate a random name so slugs don't collide
	randomStr := acctest.RandStringFromCharSet(5, acctest.CharSetAlpha)
	projectString := fmt.Sprintf("Terraform test project %s", randomStr)
	projectSlug := fmt.Sprintf("terraform-test-project-%s", randomStr)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			/* Skipping v0.14 & 0.15 to avoid the following error:
			   Provider "registry.terraform.io/hashicorp/sleuth" planned an invalid value
			   for
			   sleuth_metric_impact_source.terraform_acc_test_dd.manually_set_health_threshold:
			   planned value cty.NumberFloatVal(99.99) does not match config value
			   cty.MustParseNumberVal("99.99").

				No issues here, just a SDK bug/blurp
			*/
			tfversion.SkipBetween(tfversion.Version0_14_0, tfversion.Version0_15_0),
		},
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: createCodeChangeConfig(projectString),

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sleuth_project.terraform_acc_test", "name", projectString),
					resource.TestCheckResourceAttr("sleuth_project.terraform_acc_test", "slug", projectSlug),

					resource.TestCheckResourceAttr("sleuth_environment.terraform_acc_test", "name", "staging"),
					resource.TestCheckResourceAttr("sleuth_environment.terraform_acc_test", "slug", "staging"),

					resource.TestCheckResourceAttr("sleuth_code_change_source.terraform_acc_test", "name", "Terraform code change source"),
					resource.TestCheckResourceAttr("sleuth_code_change_source.terraform_acc_test", "deploy_tracking_type", "build"),
					resource.TestCheckResourceAttr("sleuth_code_change_source.terraform_acc_test", "collect_impact", "true"),

					resource.TestCheckResourceAttr("sleuth_code_change_source.terraform_acc_test", "repository.name", "Terraform provider sleuth"),
					resource.TestCheckResourceAttr("sleuth_code_change_source.terraform_acc_test", "repository.provider", "GITHUB"),
					resource.TestCheckResourceAttr("sleuth_code_change_source.terraform_acc_test", "repository.url", "https://github.com/sleuth-io/terraform-provider-sleuth"),

					resource.TestCheckResourceAttr("sleuth_code_change_source.terraform_acc_test", "environment_mappings.#", "2"),
					resource.TestCheckResourceAttr("sleuth_code_change_source.terraform_acc_test", "environment_mappings.0.environment_slug", "production"),
					resource.TestCheckResourceAttr("sleuth_code_change_source.terraform_acc_test", "environment_mappings.0.branch", "main"),
					resource.TestCheckResourceAttr("sleuth_code_change_source.terraform_acc_test", "environment_mappings.1.environment_slug", "staging"),
					resource.TestCheckResourceAttr("sleuth_code_change_source.terraform_acc_test", "environment_mappings.1.branch", "main"),

					resource.TestCheckResourceAttr("sleuth_code_change_source.terraform_acc_test", "build_mappings.#", "2"),
					resource.TestCheckResourceAttr("sleuth_code_change_source.terraform_acc_test", "build_mappings.0.environment_slug", "production"),
					resource.TestCheckResourceAttr("sleuth_code_change_source.terraform_acc_test", "build_mappings.0.build_name", "release"),
					resource.TestCheckResourceAttr("sleuth_code_change_source.terraform_acc_test", "build_mappings.0.project_key", "sleuth-io/terraform-provider-sleuth"),
					resource.TestCheckResourceAttr("sleuth_code_change_source.terraform_acc_test", "build_mappings.0.match_branch_to_environment", "false"),
					resource.TestCheckResourceAttr("sleuth_code_change_source.terraform_acc_test", "build_mappings.1.environment_slug", "staging"),
					resource.TestCheckResourceAttr("sleuth_code_change_source.terraform_acc_test", "build_mappings.1.build_name", "Tests"),
					resource.TestCheckResourceAttr("sleuth_code_change_source.terraform_acc_test", "build_mappings.1.project_key", "sleuth-io/terraform-provider-sleuth"),
					resource.TestCheckResourceAttr("sleuth_code_change_source.terraform_acc_test", "build_mappings.1.match_branch_to_environment", "true"),
				),
			},
			{
				Config:             createCodeChangeConfig(projectString),
				ExpectNonEmptyPlan: false,
				Check: func(state *terraform.State) error {
					// Manually add delay because if we try to delete Code change it will fail because it's still in initailizing state
					// If you are getting errors with this test, try increasing the sleep time
					time.Sleep(60 * time.Second)
					return nil
				},
			},
			// Update testing
			{
				Config: updateCodeChangeConfig(projectString),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sleuth_project.terraform_acc_test", "name", projectString),
					resource.TestCheckResourceAttr("sleuth_project.terraform_acc_test", "slug", projectSlug),

					resource.TestCheckResourceAttr("sleuth_environment.terraform_acc_test", "name", "staging"),
					resource.TestCheckResourceAttr("sleuth_environment.terraform_acc_test", "slug", "staging"),

					resource.TestCheckResourceAttr("sleuth_code_change_source.terraform_acc_test", "name", "Terraform code change source updated"),
					resource.TestCheckResourceAttr("sleuth_code_change_source.terraform_acc_test", "deploy_tracking_type", "manual"),
					resource.TestCheckResourceAttr("sleuth_code_change_source.terraform_acc_test", "collect_impact", "false"),

					resource.TestCheckResourceAttr("sleuth_code_change_source.terraform_acc_test", "repository.name", "terraform-provider-sleuth"),
					resource.TestCheckResourceAttr("sleuth_code_change_source.terraform_acc_test", "repository.provider", "GITHUB"),
					resource.TestCheckResourceAttr("sleuth_code_change_source.terraform_acc_test", "repository.url", "https://github.com/sleuth-io/terraform-provider-sleuth"),

					resource.TestCheckResourceAttr("sleuth_code_change_source.terraform_acc_test", "environment_mappings.#", "2"),
					resource.TestCheckResourceAttr("sleuth_code_change_source.terraform_acc_test", "environment_mappings.0.environment_slug", "production"),
					resource.TestCheckResourceAttr("sleuth_code_change_source.terraform_acc_test", "environment_mappings.0.branch", "main"),
					resource.TestCheckResourceAttr("sleuth_code_change_source.terraform_acc_test", "environment_mappings.1.environment_slug", "staging"),
					resource.TestCheckResourceAttr("sleuth_code_change_source.terraform_acc_test", "environment_mappings.1.branch", "main"),

					resource.TestCheckNoResourceAttr("sleuth_code_change_source.terraform_acc_test", "build_mappings.#"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccCodeChangeSourceResource_v5(t *testing.T) {
	// tests are run in parallel both locally & on CI, so we need to generate a random name so slugs don't collide
	randomStr := acctest.RandStringFromCharSet(5, acctest.CharSetAlpha)
	projectString := fmt.Sprintf("Terraform test project %s", randomStr)
	projectSlug := fmt.Sprintf("terraform-test-project-%s", randomStr)

	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			/* Skipping v0.14 & 0.15 to avoid the following error:
			   Provider "registry.terraform.io/hashicorp/sleuth" planned an invalid value
			   for
			   sleuth_metric_impact_source.terraform_acc_test_dd.manually_set_health_threshold:
			   planned value cty.NumberFloatVal(99.99) does not match config value
			   cty.MustParseNumberVal("99.99").

				No issues here, just a SDK bug/blurp
			*/
			tfversion.SkipBetween(tfversion.Version0_14_0, tfversion.Version0_15_0),
		},
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: createCodeChangeConfig(projectString),

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sleuth_project.terraform_acc_test", "name", projectString),
					resource.TestCheckResourceAttr("sleuth_project.terraform_acc_test", "slug", projectSlug),

					resource.TestCheckResourceAttr("sleuth_environment.terraform_acc_test", "name", "staging"),
					resource.TestCheckResourceAttr("sleuth_environment.terraform_acc_test", "slug", "staging"),

					resource.TestCheckResourceAttr("sleuth_code_change_source.terraform_acc_test", "name", "Terraform code change source"),
					resource.TestCheckResourceAttr("sleuth_code_change_source.terraform_acc_test", "deploy_tracking_type", "build"),
					resource.TestCheckResourceAttr("sleuth_code_change_source.terraform_acc_test", "collect_impact", "true"),

					resource.TestCheckResourceAttr("sleuth_code_change_source.terraform_acc_test", "repository.name", "Terraform provider sleuth"),
					resource.TestCheckResourceAttr("sleuth_code_change_source.terraform_acc_test", "repository.provider", "GITHUB"),
					resource.TestCheckResourceAttr("sleuth_code_change_source.terraform_acc_test", "repository.url", "https://github.com/sleuth-io/terraform-provider-sleuth"),

					resource.TestCheckResourceAttr("sleuth_code_change_source.terraform_acc_test", "environment_mappings.#", "2"),
					resource.TestCheckResourceAttr("sleuth_code_change_source.terraform_acc_test", "environment_mappings.0.environment_slug", "production"),
					resource.TestCheckResourceAttr("sleuth_code_change_source.terraform_acc_test", "environment_mappings.0.branch", "main"),
					resource.TestCheckResourceAttr("sleuth_code_change_source.terraform_acc_test", "environment_mappings.1.environment_slug", "staging"),
					resource.TestCheckResourceAttr("sleuth_code_change_source.terraform_acc_test", "environment_mappings.1.branch", "main"),

					resource.TestCheckResourceAttr("sleuth_code_change_source.terraform_acc_test", "build_mappings.#", "2"),
					resource.TestCheckResourceAttr("sleuth_code_change_source.terraform_acc_test", "build_mappings.0.environment_slug", "production"),
					resource.TestCheckResourceAttr("sleuth_code_change_source.terraform_acc_test", "build_mappings.0.build_name", "release"),
					resource.TestCheckResourceAttr("sleuth_code_change_source.terraform_acc_test", "build_mappings.0.project_key", "sleuth-io/terraform-provider-sleuth"),
					resource.TestCheckResourceAttr("sleuth_code_change_source.terraform_acc_test", "build_mappings.0.match_branch_to_environment", "false"),
					resource.TestCheckResourceAttr("sleuth_code_change_source.terraform_acc_test", "build_mappings.1.environment_slug", "staging"),
					resource.TestCheckResourceAttr("sleuth_code_change_source.terraform_acc_test", "build_mappings.1.build_name", "Tests"),
					resource.TestCheckResourceAttr("sleuth_code_change_source.terraform_acc_test", "build_mappings.1.project_key", "sleuth-io/terraform-provider-sleuth"),
					resource.TestCheckResourceAttr("sleuth_code_change_source.terraform_acc_test", "build_mappings.1.match_branch_to_environment", "true"),
				),
			},
			{
				Config:             createCodeChangeConfig(projectString),
				ExpectNonEmptyPlan: false,
				Check: func(state *terraform.State) error {
					// Manually add delay because if we try to delete Code change it will fail because it's still in initailizing state
					// If you are getting errors with this test, try increasing the sleep time
					time.Sleep(60 * time.Second)
					return nil
				},
			},
			// Update testing
			{
				Config: updateCodeChangeConfig(projectString),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sleuth_project.terraform_acc_test", "name", projectString),
					resource.TestCheckResourceAttr("sleuth_project.terraform_acc_test", "slug", projectSlug),

					resource.TestCheckResourceAttr("sleuth_environment.terraform_acc_test", "name", "staging"),
					resource.TestCheckResourceAttr("sleuth_environment.terraform_acc_test", "slug", "staging"),

					resource.TestCheckResourceAttr("sleuth_code_change_source.terraform_acc_test", "name", "Terraform code change source updated"),
					resource.TestCheckResourceAttr("sleuth_code_change_source.terraform_acc_test", "deploy_tracking_type", "manual"),
					resource.TestCheckResourceAttr("sleuth_code_change_source.terraform_acc_test", "collect_impact", "false"),

					resource.TestCheckResourceAttr("sleuth_code_change_source.terraform_acc_test", "repository.name", "terraform-provider-sleuth"),
					resource.TestCheckResourceAttr("sleuth_code_change_source.terraform_acc_test", "repository.provider", "GITHUB"),
					resource.TestCheckResourceAttr("sleuth_code_change_source.terraform_acc_test", "repository.url", "https://github.com/sleuth-io/terraform-provider-sleuth"),

					resource.TestCheckResourceAttr("sleuth_code_change_source.terraform_acc_test", "environment_mappings.#", "2"),
					resource.TestCheckResourceAttr("sleuth_code_change_source.terraform_acc_test", "environment_mappings.0.environment_slug", "production"),
					resource.TestCheckResourceAttr("sleuth_code_change_source.terraform_acc_test", "environment_mappings.0.branch", "main"),
					resource.TestCheckResourceAttr("sleuth_code_change_source.terraform_acc_test", "environment_mappings.1.environment_slug", "staging"),
					resource.TestCheckResourceAttr("sleuth_code_change_source.terraform_acc_test", "environment_mappings.1.branch", "main"),

					resource.TestCheckNoResourceAttr("sleuth_code_change_source.terraform_acc_test", "build_mappings.#"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

// Production env is created automatically, and we don't want to track it inside TF otherwise we can't delete project because of the env blocking deletion
func createCodeChangeConfig(name string) string {
	return fmt.Sprintf(`
resource "sleuth_project" "terraform_acc_test" {
	name = "%s"
}

resource "sleuth_environment" "terraform_acc_test" {
	project_slug = sleuth_project.terraform_acc_test.slug
	name = "staging"
}

resource "sleuth_code_change_source" "terraform_acc_test" {
    project_slug = sleuth_project.terraform_acc_test.slug
    name = "Terraform code change source"
    repository {
        name = "Terraform provider sleuth"
        owner = "sleuth-io"
        provider = "GITHUB"
        url = "https://github.com/sleuth-io/terraform-provider-sleuth"
    }
    environment_mappings {
        environment_slug = "production"
        branch = "main"
    }
    environment_mappings {
        environment_slug = sleuth_environment.terraform_acc_test.slug
        branch = "main"
    }

	build_mappings {
		environment_slug = "production"
		build_name = "release"
		project_key = "sleuth-io/terraform-provider-sleuth"
		provider = "GITHUB"
		match_branch_to_environment = false
	}
    build_mappings {
        environment_slug = sleuth_environment.terraform_acc_test.slug
        build_name = "Tests"
        project_key = "sleuth-io/terraform-provider-sleuth"
        provider = "GITHUB"
        match_branch_to_environment = true
    }

  	deploy_tracking_type = "build"
  	collect_impact = true
}
`, name)
}

func updateCodeChangeConfig(name string) string {
	return fmt.Sprintf(`
resource "sleuth_project" "terraform_acc_test" {
	name = "%s"
}

resource "sleuth_environment" "terraform_acc_test" {
	project_slug = sleuth_project.terraform_acc_test.slug
	name = "staging"
}

resource "sleuth_code_change_source" "terraform_acc_test" {
    project_slug = sleuth_project.terraform_acc_test.slug
    name = "Terraform code change source updated"
    repository {
        name = "terraform-provider-sleuth"
        owner = "sleuth-io"
        provider = "GITHUB"
        url = "https://github.com/sleuth-io/terraform-provider-sleuth"
    }
    environment_mappings {
        environment_slug = "production"
        branch = "main"
    }
    environment_mappings {
        environment_slug = sleuth_environment.terraform_acc_test.slug
        branch = "main"
    }

  	deploy_tracking_type = "manual"
  	collect_impact = false
}
`, name)
}

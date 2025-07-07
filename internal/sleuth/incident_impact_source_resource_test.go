package sleuth

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAccIncidentImpactSourceResource_v6(t *testing.T) {
	// tests are run in parallel both locally & on CI, so we need to generate a random name so slugs don't collide
	randomStr := acctest.RandStringFromCharSet(5, acctest.CharSetAlpha)
	projectString := fmt.Sprintf("Terraform test project %s", randomStr)
	projectSlug := fmt.Sprintf("terraform-test-project-%s", randomStr)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			/* Skipping v0.14 up till v1.2 to avoid the following error:
			 Error: Provider produced invalid plan

			Provider "registry.terraform.io/hashicorp/sleuth" planned an invalid value
			for sleuth_incident_impact_source.terraform_acc_test_pd.pagerduty_input:
			planned for existence but config wants absense.

			See:
			- https://github.com/hashicorp/terraform/pull/32463
			- https://github.com/hashicorp/terraform-plugin-framework/issues/603
			*/
			tfversion.SkipBetween(tfversion.Version0_14_0, tfversion.Version1_3_0),
		},
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: createIncidentImpactConfig(projectString),

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sleuth_project.terraform_acc_test", "name", projectString),
					resource.TestCheckResourceAttr("sleuth_project.terraform_acc_test", "slug", projectSlug),

					resource.TestCheckResourceAttr("sleuth_environment.terraform_acc_test", "name", "staging"),
					resource.TestCheckResourceAttr("sleuth_environment.terraform_acc_test", "slug", "staging"),
					// PagerDuty
					//resource.TestCheckResourceAttr("sleuth_incident_impact_source.terraform_acc_test_pd", "name", "PagerDuty TF incident impact"),
					//resource.TestCheckResourceAttr("sleuth_incident_impact_source.terraform_acc_test_pd", "environment_name", "staging"),
					//resource.TestCheckResourceAttr("sleuth_incident_impact_source.terraform_acc_test_pd", "provider_name", "pagerduty"),
					//resource.TestCheckResourceAttr("sleuth_incident_impact_source.terraform_acc_test_pd", "pagerduty_input.remote_services", "PIMPOA4"),
					//resource.TestCheckResourceAttr("sleuth_incident_impact_source.terraform_acc_test_pd", "pagerduty_input.remote_urgency", "ANY"),
					// DataDog
					resource.TestCheckResourceAttr("sleuth_incident_impact_source.terraform_acc_test_dd", "name", "DataDog TF incident impact"),
					resource.TestCheckResourceAttr("sleuth_incident_impact_source.terraform_acc_test_dd", "environment_name", "staging"),
					resource.TestCheckResourceAttr("sleuth_incident_impact_source.terraform_acc_test_dd", "provider_name", "datadog"),
					resource.TestCheckResourceAttr("sleuth_incident_impact_source.terraform_acc_test_dd", "datadog_input.query", "@query=1234"),
					resource.TestCheckResourceAttr("sleuth_incident_impact_source.terraform_acc_test_dd", "datadog_input.remote_priority_threshold", "P4"),
					//resource.TestCheckResourceAttr("sleuth_incident_impact_source.terraform_acc_test_dd", "pagerduty_input.integration_slug", "datadog"),
					// JIRA
					resource.TestCheckResourceAttr("sleuth_incident_impact_source.terraform_acc_test_jira", "name", "Jira TF incident impact"),
					resource.TestCheckResourceAttr("sleuth_incident_impact_source.terraform_acc_test_jira", "environment_name", "staging"),
					resource.TestCheckResourceAttr("sleuth_incident_impact_source.terraform_acc_test_jira", "provider_name", "jira"),
					resource.TestCheckResourceAttr("sleuth_incident_impact_source.terraform_acc_test_jira", "jira_input.remote_jql", "created >= -30d order by created DESC"),
				),
			},
			// Update testing
			{
				Config: updateIncidentImpactConfig(projectString),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sleuth_project.terraform_acc_test", "name", projectString),
					resource.TestCheckResourceAttr("sleuth_project.terraform_acc_test", "slug", projectSlug),

					resource.TestCheckResourceAttr("sleuth_environment.terraform_acc_test", "name", "staging"),
					resource.TestCheckResourceAttr("sleuth_environment.terraform_acc_test", "slug", "staging"),
					// PagerDuty
					//resource.TestCheckResourceAttr("sleuth_incident_impact_source.terraform_acc_test_pd", "name", "PagerDuty TF incident impact updated"),
					//resource.TestCheckResourceAttr("sleuth_incident_impact_source.terraform_acc_test_pd", "environment_name", "staging"),
					//resource.TestCheckResourceAttr("sleuth_incident_impact_source.terraform_acc_test_pd", "provider_name", "pagerduty"),
					//resource.TestCheckNoResourceAttr("sleuth_incident_impact_source.terraform_acc_test_pd", "pagerduty_input.remote_services"),
					//resource.TestCheckNoResourceAttr("sleuth_incident_impact_source.terraform_acc_test_pd", "pagerduty_input.remote_urgency"),
					// DataDog
					resource.TestCheckResourceAttr("sleuth_incident_impact_source.terraform_acc_test_dd", "name", "DataDog TF incident impact updated"),
					resource.TestCheckResourceAttr("sleuth_incident_impact_source.terraform_acc_test_dd", "environment_name", "staging"),
					resource.TestCheckResourceAttr("sleuth_incident_impact_source.terraform_acc_test_dd", "provider_name", "datadog"),
					resource.TestCheckResourceAttr("sleuth_incident_impact_source.terraform_acc_test_dd", "datadog_input.query", ""),
					resource.TestCheckResourceAttr("sleuth_incident_impact_source.terraform_acc_test_dd", "datadog_input.remote_priority_threshold", "ALL"),
					// JIRA
					resource.TestCheckResourceAttr("sleuth_incident_impact_source.terraform_acc_test_jira", "name", "Jira TF incident impact updated"),
					resource.TestCheckResourceAttr("sleuth_incident_impact_source.terraform_acc_test_jira", "environment_name", "staging"),
					resource.TestCheckResourceAttr("sleuth_incident_impact_source.terraform_acc_test_jira", "provider_name", "jira"),
					resource.TestCheckResourceAttr("sleuth_incident_impact_source.terraform_acc_test_jira", "jira_input.remote_jql", "created >= -30d order by created ASC"),
				),
			},
			{
				Config: updateIncidentImpactConfig2(projectString),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sleuth_project.terraform_acc_test", "name", projectString),
					resource.TestCheckResourceAttr("sleuth_project.terraform_acc_test", "slug", projectSlug),

					resource.TestCheckResourceAttr("sleuth_environment.terraform_acc_test", "name", "staging"),
					resource.TestCheckResourceAttr("sleuth_environment.terraform_acc_test", "slug", "staging"),
					// PagerDuty
					//resource.TestCheckResourceAttr("sleuth_incident_impact_source.terraform_acc_test_pd", "name", "PagerDuty TF incident impact"),
					//resource.TestCheckResourceAttr("sleuth_incident_impact_source.terraform_acc_test_pd", "environment_name", "Production"),
					//resource.TestCheckResourceAttr("sleuth_incident_impact_source.terraform_acc_test_pd", "provider_name", "pagerduty"),
					//resource.TestCheckResourceAttr("sleuth_incident_impact_source.terraform_acc_test_pd", "pagerduty_input.remote_services", ""),
					//resource.TestCheckResourceAttr("sleuth_incident_impact_source.terraform_acc_test_pd", "pagerduty_input.remote_urgency", "HIGH"),
					// DataDog
					resource.TestCheckResourceAttr("sleuth_incident_impact_source.terraform_acc_test_dd", "name", "DataDog TF incident impact"),
					resource.TestCheckResourceAttr("sleuth_incident_impact_source.terraform_acc_test_dd", "environment_name", "Production"),
					resource.TestCheckResourceAttr("sleuth_incident_impact_source.terraform_acc_test_dd", "provider_name", "datadog"),
					resource.TestCheckResourceAttr("sleuth_incident_impact_source.terraform_acc_test_dd", "datadog_input.query", "@query=12345"),
					resource.TestCheckResourceAttr("sleuth_incident_impact_source.terraform_acc_test_dd", "datadog_input.remote_priority_threshold", "P3"),
					// JIRA
					resource.TestCheckResourceAttr("sleuth_incident_impact_source.terraform_acc_test_jira", "name", "Jira TF incident impact"),
					resource.TestCheckResourceAttr("sleuth_incident_impact_source.terraform_acc_test_jira", "environment_name", "Production"),
					resource.TestCheckResourceAttr("sleuth_incident_impact_source.terraform_acc_test_jira", "provider_name", "jira"),
					resource.TestCheckResourceAttr("sleuth_incident_impact_source.terraform_acc_test_jira", "jira_input.remote_jql", "created >= -10d order by created DESC"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

/* For local testing update DD & JIRA auth slugs */
func createIncidentImpactConfig(name string) string {
	return fmt.Sprintf(`
resource "sleuth_project" "terraform_acc_test" {
	name = "%s"
}

resource "sleuth_environment" "terraform_acc_test" {
	project_slug = sleuth_project.terraform_acc_test.slug
	name = "staging"
}

//resource "sleuth_incident_impact_source" "terraform_acc_test_pd" {
//  	project_slug = sleuth_project.terraform_acc_test.slug
//  	name = "PagerDuty TF incident impact"
//  	environment_name = sleuth_environment.terraform_acc_test.name
//  	provider_name = "pagerduty"
//	pagerduty_input = {
//		remote_services = "PIMPOA4"
//		remote_urgency = "ANY"
//	}
//}

resource "sleuth_incident_impact_source" "terraform_acc_test_dd" {
 	project_slug = sleuth_project.terraform_acc_test.slug
 	name = "DataDog TF incident impact"
 	environment_name = sleuth_environment.terraform_acc_test.name
 	provider_name = "datadog"
	datadog_input = {
        query = "@query=1234"
        remote_priority_threshold = "P4"
        #integration_slug = "datadog-local"
    }
}

resource "sleuth_incident_impact_source" "terraform_acc_test_jira" {
 	project_slug = sleuth_project.terraform_acc_test.slug
 	name = "Jira TF incident impact"
 	environment_name = sleuth_environment.terraform_acc_test.name
 	provider_name = "jira"
	jira_input = {
		remote_jql = "created >= -30d order by created DESC"
		integration_slug = "jira-cloud-jira-hot"
	}
}
`, name)
}

func updateIncidentImpactConfig(name string) string {
	return fmt.Sprintf(`
resource "sleuth_project" "terraform_acc_test" {
	name = "%s"
}

resource "sleuth_environment" "terraform_acc_test" {
	project_slug = sleuth_project.terraform_acc_test.slug
	name = "staging"
}

//resource "sleuth_incident_impact_source" "terraform_acc_test_pd" {
//  	project_slug = sleuth_project.terraform_acc_test.slug
//  	name = "PagerDuty TF incident impact updated"
//  	environment_name = sleuth_environment.terraform_acc_test.name
//  	provider_name = "pagerduty"
//}

resource "sleuth_incident_impact_source" "terraform_acc_test_dd" {
 	project_slug = sleuth_project.terraform_acc_test.slug
 	name = "DataDog TF incident impact updated"
 	environment_name = sleuth_environment.terraform_acc_test.name
 	provider_name = "datadog"
	datadog_input = {
        remote_priority_threshold = "ALL"
    }
}

# remote_jql is always required
resource "sleuth_incident_impact_source" "terraform_acc_test_jira" {
 	project_slug = sleuth_project.terraform_acc_test.slug
 	name = "Jira TF incident impact updated"
 	environment_name = sleuth_environment.terraform_acc_test.name
 	provider_name = "jira"
	jira_input = {
		remote_jql = "created >= -30d order by created ASC"
		integration_slug = "jira-cloud-jira-hot"
	}
}
`, name)
}

func updateIncidentImpactConfig2(name string) string {
	return fmt.Sprintf(`
resource "sleuth_project" "terraform_acc_test" {
	name = "%s"
}

resource "sleuth_environment" "terraform_acc_test" {
	project_slug = sleuth_project.terraform_acc_test.slug
	name = "staging"
}

//resource "sleuth_incident_impact_source" "terraform_acc_test_pd" {
//  	project_slug = sleuth_project.terraform_acc_test.slug
//  	name = "PagerDuty TF incident impact"
//  	environment_name = "Production"
//  	provider_name = "pagerduty"
//	pagerduty_input = {
//		remote_urgency = "HIGH"
//	}
//}

resource "sleuth_incident_impact_source" "terraform_acc_test_dd" {
 	project_slug = sleuth_project.terraform_acc_test.slug
 	name = "DataDog TF incident impact"
 	environment_name = "Production"
 	provider_name = "datadog"
	datadog_input = {
		query = "@query=12345"
        remote_priority_threshold = "P3"
    }
}

resource "sleuth_incident_impact_source" "terraform_acc_test_jira" {
 	project_slug = sleuth_project.terraform_acc_test.slug
 	name = "Jira TF incident impact"
 	environment_name = "Production"
 	provider_name = "jira"
	jira_input = {
		remote_jql = "created >= -10d order by created DESC"
		integration_slug = "jira-cloud-jira-hot"
	}
}
`, name)
}

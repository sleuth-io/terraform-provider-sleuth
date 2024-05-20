package sleuth

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAccMetricImpactSourceResource_v6(t *testing.T) {
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
			tfversion.SkipBetween(tfversion.Version0_14_0, tfversion.Version1_0_0),
		},
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: createMetricImpactConfig(projectString),

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sleuth_project.terraform_acc_test", "name", projectString),
					resource.TestCheckResourceAttr("sleuth_project.terraform_acc_test", "slug", projectSlug),

					resource.TestCheckResourceAttr("sleuth_environment.terraform_acc_test", "name", "staging"),
					resource.TestCheckResourceAttr("sleuth_environment.terraform_acc_test", "slug", "staging"),

					resource.TestCheckResourceAttr("sleuth_metric_impact_source.terraform_acc_test_dd", "name", "Datadog acceptance test"),
					resource.TestCheckResourceAttr("sleuth_metric_impact_source.terraform_acc_test_dd", "slug", "datadog-acceptance-test"),
					resource.TestCheckResourceAttr("sleuth_metric_impact_source.terraform_acc_test_dd", "provider_type", "DATADOG"),
					resource.TestCheckResourceAttr("sleuth_metric_impact_source.terraform_acc_test_dd", "query", "avg(last_5m):avg:datadog.agent.check_run.duration{check:datadog.agent.up}.as_count() > 0"),
					resource.TestCheckResourceAttr("sleuth_metric_impact_source.terraform_acc_test_dd", "less_is_better", "false"),

					resource.TestCheckResourceAttrSet("sleuth_metric_impact_source.terraform_acc_test_dd", "id"),
					resource.TestCheckResourceAttrSet("sleuth_metric_impact_source.terraform_acc_test_dd", "integration_slug"),

					resource.TestCheckResourceAttr("sleuth_metric_impact_source.terraform_acc_test_cw", "name", "RDS CPU"),
					resource.TestCheckResourceAttr("sleuth_metric_impact_source.terraform_acc_test_cw", "slug", "rds-cpu"),
					resource.TestCheckResourceAttr("sleuth_metric_impact_source.terraform_acc_test_cw", "provider_type", "CLOUDWATCH"),
					resource.TestCheckResourceAttr("sleuth_metric_impact_source.terraform_acc_test_cw", "query", "{\"metrics\":[[\"AWS/RDS\",\"CPUUtilization\",\"DBInstanceIdentifier\",\"my-db-identifier\",{\"id\":\"m1\"}]],\"period\":300,\"region\":\"us-east-1\",\"stacked\":false,\"stat\":\"Average\",\"view\":\"timeSeries\"}"),
					resource.TestCheckResourceAttr("sleuth_metric_impact_source.terraform_acc_test_cw", "less_is_better", "true"),

					resource.TestCheckResourceAttrSet("sleuth_metric_impact_source.terraform_acc_test_cw", "id"),
					resource.TestCheckResourceAttrSet("sleuth_metric_impact_source.terraform_acc_test_cw", "integration_slug"),
				),
			},
			// Update testing
			{
				Config: updateMetricImpactConfig(projectString),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sleuth_project.terraform_acc_test", "name", projectString),
					resource.TestCheckResourceAttr("sleuth_project.terraform_acc_test", "slug", projectSlug),

					resource.TestCheckResourceAttr("sleuth_environment.terraform_acc_test", "name", "staging"),
					resource.TestCheckResourceAttr("sleuth_environment.terraform_acc_test", "slug", "staging"),

					resource.TestCheckResourceAttr("sleuth_metric_impact_source.terraform_acc_test_dd", "name", "Datadog acceptance test update"),
					resource.TestCheckResourceAttr("sleuth_metric_impact_source.terraform_acc_test_dd", "slug", "datadog-acceptance-test"),
					resource.TestCheckResourceAttr("sleuth_metric_impact_source.terraform_acc_test_dd", "provider_type", "DATADOG"),
					resource.TestCheckResourceAttr("sleuth_metric_impact_source.terraform_acc_test_dd", "query", "avg(last_5m):avg:datadog.agent.check_run.duration{check:datadog.agent.up}.as_count() > 100"),
					resource.TestCheckResourceAttr("sleuth_metric_impact_source.terraform_acc_test_dd", "less_is_better", "true"),
					resource.TestCheckResourceAttr("sleuth_metric_impact_source.terraform_acc_test_dd", "manually_set_health_threshold", "99.99"),

					resource.TestCheckResourceAttr("sleuth_metric_impact_source.terraform_acc_test_cw", "name", "RDS CPU updated"),
					resource.TestCheckResourceAttr("sleuth_metric_impact_source.terraform_acc_test_cw", "slug", "rds-cpu"),
					resource.TestCheckResourceAttr("sleuth_metric_impact_source.terraform_acc_test_cw", "provider_type", "CLOUDWATCH"),
					resource.TestCheckResourceAttr("sleuth_metric_impact_source.terraform_acc_test_cw", "query", "{\"metrics\":[[\"AWS/RDS\",\"CPUUtilization\",\"DBInstanceIdentifier\",\"my-db-identifier\",{\"id\":\"m1\"}]],\"period\":600,\"region\":\"us-east-1\",\"stacked\":false,\"stat\":\"Average\",\"view\":\"timeSeries\"}"),
					resource.TestCheckResourceAttr("sleuth_metric_impact_source.terraform_acc_test_cw", "less_is_better", "false"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccMetricImpactSourceResource_v5(t *testing.T) {
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
			tfversion.SkipBetween(tfversion.Version0_14_0, tfversion.Version1_0_0),
		},
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: createMetricImpactConfig(projectString),

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sleuth_project.terraform_acc_test", "name", projectString),
					resource.TestCheckResourceAttr("sleuth_project.terraform_acc_test", "slug", projectSlug),

					resource.TestCheckResourceAttr("sleuth_environment.terraform_acc_test", "name", "staging"),
					resource.TestCheckResourceAttr("sleuth_environment.terraform_acc_test", "slug", "staging"),

					resource.TestCheckResourceAttr("sleuth_metric_impact_source.terraform_acc_test_dd", "name", "Datadog acceptance test"),
					resource.TestCheckResourceAttr("sleuth_metric_impact_source.terraform_acc_test_dd", "slug", "datadog-acceptance-test"),
					resource.TestCheckResourceAttr("sleuth_metric_impact_source.terraform_acc_test_dd", "provider_type", "DATADOG"),
					resource.TestCheckResourceAttr("sleuth_metric_impact_source.terraform_acc_test_dd", "query", "avg(last_5m):avg:datadog.agent.check_run.duration{check:datadog.agent.up}.as_count() > 0"),
					resource.TestCheckResourceAttr("sleuth_metric_impact_source.terraform_acc_test_dd", "less_is_better", "false"),

					resource.TestCheckResourceAttrSet("sleuth_metric_impact_source.terraform_acc_test_dd", "id"),
					resource.TestCheckResourceAttrSet("sleuth_metric_impact_source.terraform_acc_test_dd", "integration_slug"),

					resource.TestCheckResourceAttr("sleuth_metric_impact_source.terraform_acc_test_cw", "name", "RDS CPU"),
					resource.TestCheckResourceAttr("sleuth_metric_impact_source.terraform_acc_test_cw", "slug", "rds-cpu"),
					resource.TestCheckResourceAttr("sleuth_metric_impact_source.terraform_acc_test_cw", "provider_type", "CLOUDWATCH"),
					resource.TestCheckResourceAttr("sleuth_metric_impact_source.terraform_acc_test_cw", "query", "{\"metrics\":[[\"AWS/RDS\",\"CPUUtilization\",\"DBInstanceIdentifier\",\"my-db-identifier\",{\"id\":\"m1\"}]],\"period\":300,\"region\":\"us-east-1\",\"stacked\":false,\"stat\":\"Average\",\"view\":\"timeSeries\"}"),
					resource.TestCheckResourceAttr("sleuth_metric_impact_source.terraform_acc_test_cw", "less_is_better", "true"),

					resource.TestCheckResourceAttrSet("sleuth_metric_impact_source.terraform_acc_test_cw", "id"),
					resource.TestCheckResourceAttrSet("sleuth_metric_impact_source.terraform_acc_test_cw", "integration_slug"),
				),
			},
			// Update testing
			{
				Config: updateMetricImpactConfig(projectString),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sleuth_project.terraform_acc_test", "name", projectString),
					resource.TestCheckResourceAttr("sleuth_project.terraform_acc_test", "slug", projectSlug),

					resource.TestCheckResourceAttr("sleuth_environment.terraform_acc_test", "name", "staging"),
					resource.TestCheckResourceAttr("sleuth_environment.terraform_acc_test", "slug", "staging"),

					resource.TestCheckResourceAttr("sleuth_metric_impact_source.terraform_acc_test_dd", "name", "Datadog acceptance test update"),
					resource.TestCheckResourceAttr("sleuth_metric_impact_source.terraform_acc_test_dd", "slug", "datadog-acceptance-test"),
					resource.TestCheckResourceAttr("sleuth_metric_impact_source.terraform_acc_test_dd", "provider_type", "DATADOG"),
					resource.TestCheckResourceAttr("sleuth_metric_impact_source.terraform_acc_test_dd", "query", "avg(last_5m):avg:datadog.agent.check_run.duration{check:datadog.agent.up}.as_count() > 100"),
					resource.TestCheckResourceAttr("sleuth_metric_impact_source.terraform_acc_test_dd", "less_is_better", "true"),
					resource.TestCheckResourceAttr("sleuth_metric_impact_source.terraform_acc_test_dd", "manually_set_health_threshold", "99.99"),

					resource.TestCheckResourceAttr("sleuth_metric_impact_source.terraform_acc_test_cw", "name", "RDS CPU updated"),
					resource.TestCheckResourceAttr("sleuth_metric_impact_source.terraform_acc_test_cw", "slug", "rds-cpu"),
					resource.TestCheckResourceAttr("sleuth_metric_impact_source.terraform_acc_test_cw", "provider_type", "CLOUDWATCH"),
					resource.TestCheckResourceAttr("sleuth_metric_impact_source.terraform_acc_test_cw", "query", "{\"metrics\":[[\"AWS/RDS\",\"CPUUtilization\",\"DBInstanceIdentifier\",\"my-db-identifier\",{\"id\":\"m1\"}]],\"period\":600,\"region\":\"us-east-1\",\"stacked\":false,\"stat\":\"Average\",\"view\":\"timeSeries\"}"),
					resource.TestCheckResourceAttr("sleuth_metric_impact_source.terraform_acc_test_cw", "less_is_better", "false"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func createMetricImpactConfig(name string) string {
	return fmt.Sprintf(`
resource "sleuth_project" "terraform_acc_test" {
	name = "%s"
}

resource "sleuth_environment" "terraform_acc_test" {
	project_slug = sleuth_project.terraform_acc_test.slug
	name = "staging"
}

resource "sleuth_metric_impact_source" "terraform_acc_test_dd" {
	project_slug = sleuth_project.terraform_acc_test.slug
	environment_slug = sleuth_environment.terraform_acc_test.slug
	name = "Datadog acceptance test"
	provider_type = "DATADOG"
	query = "avg(last_5m):avg:datadog.agent.check_run.duration{check:datadog.agent.up}.as_count() > 0"
	integration_slug="datadog"
}

resource "sleuth_metric_impact_source" "terraform_acc_test_cw" {
    project_slug = sleuth_project.terraform_acc_test.slug
    environment_slug = sleuth_environment.terraform_acc_test.slug
    name = "RDS CPU"
    provider_type = "CLOUDWATCH"
    query = jsonencode({
        "metrics": [
            [ "AWS/RDS", "CPUUtilization", "DBInstanceIdentifier", "my-db-identifier", { "id": "m1" } ]
        ],
        "view": "timeSeries",
        "stacked": false,
        "region": "us-east-1",
        "stat": "Average",
        "period": 300
    })
    less_is_better = true
    integration_slug = "aws-cloudwatch-staging-key-staging-key"
}

`, name)
}

func updateMetricImpactConfig(name string) string {
	return fmt.Sprintf(`
resource "sleuth_project" "terraform_acc_test" {
	name = "%s"
}

resource "sleuth_environment" "terraform_acc_test" {
	project_slug = sleuth_project.terraform_acc_test.slug
	name = "staging"
}

resource "sleuth_metric_impact_source" "terraform_acc_test_dd" {
	project_slug= sleuth_project.terraform_acc_test.slug
	environment_slug = sleuth_environment.terraform_acc_test.slug
	name = "Datadog acceptance test update"
	provider_type = "DATADOG"
	query = "avg(last_5m):avg:datadog.agent.check_run.duration{check:datadog.agent.up}.as_count() > 100"
	less_is_better = true
	manually_set_health_threshold = 99.99
	integration_slug = "datadog"
}

resource "sleuth_metric_impact_source" "terraform_acc_test_cw" {
	project_slug = sleuth_project.terraform_acc_test.slug
	environment_slug = sleuth_environment.terraform_acc_test.slug
	name = "RDS CPU updated"
	provider_type = "CLOUDWATCH"
    query = jsonencode({
        "metrics": [
            [ "AWS/RDS", "CPUUtilization", "DBInstanceIdentifier", "my-db-identifier", { "id": "m1" } ]
        ],
        "view": "timeSeries",
        "stacked": false,
        "region": "us-east-1",
        "stat": "Average",
        "period": 600
    })
	less_is_better = false
  	integration_slug="aws-cloudwatch-staging-key-staging-key"
}
`, name)
}

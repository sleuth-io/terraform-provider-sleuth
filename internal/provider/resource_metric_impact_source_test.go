package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceMetricImpactSource(t *testing.T) {
	if err := testAccCheckOrganization(); err != nil {
		t.Skipf("Skipping because %s.", err.Error())
	}
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceMetricImpactSource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"sleuth_metric_impact_source.mysource", "name", regexp.MustCompile("^My dd source blah")),
				),
			},
			{
				ResourceName:      "sleuth_metric_impact_source.mysource",
				ImportState:       true,
				ImportStateVerify: false,
			},
		},
	})
}

const testAccResourceMetricImpactSource = `
resource "sleuth_project" "myproject" {
	name = "My project blah"
}

resource "sleuth_environment" "myenv" {
	project_slug = "${sleuth_project.myproject.id}"
	name = "Production"
}

resource "sleuth_metric_impact_source" "mysource" {
	project_slug = "${sleuth_project.myproject.id}"
	environment_slug = "${sleuth_environment.myenv.id}"
	name = "My dd source blah"
	provider_type = "datadog"
	query = "aws.ecs.memory_utilization"
	manually_set_health_threshold = 50
}
`

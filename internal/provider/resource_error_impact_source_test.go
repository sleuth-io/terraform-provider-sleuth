package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceErrorImpactSource(t *testing.T) {
	if err := testAccCheckOrganization(); err != nil {
		t.Skipf("Skipping because %s.", err.Error())
	}
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceErrorImpactSource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"sleuth_error_impact_source.mysource", "name", regexp.MustCompile("^My source blah")),
				),
			},
			{
				ResourceName:      "sleuth_error_impact_source.mysource",
				ImportState:       true,
				ImportStateVerify: false,
			},
		},
	})
}

const testAccResourceErrorImpactSource = `
resource "sleuth_project" "myproject" {
	name = "My project blah"
}

resource "sleuth_environment" "myenv" {
	project_slug = "${sleuth_project.myproject.id}"
	name = "Production"
}

resource "sleuth_error_impact_source" "mysource" {
	project_slug = "${sleuth_project.myproject.id}"
	environment_slug = "${sleuth_environment.myenv.id}"
	name = "My source blah"
	provider_type = "sentry"
	error_org_key = "sleuth-demo"
	error_project_key = "sleuth-dev"
	error_environment = "Production"
	manually_set_health_threshold = 50
}
`

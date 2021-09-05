package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceSleuth(t *testing.T) {
	if err := testAccCheckOrganization(); err != nil {
		t.Skipf("Skipping because %s.", err.Error())
	}
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceSleuth,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"sleuth_project.myproject", "name", regexp.MustCompile("^My project blah")),
				),
			},
			{
				ResourceName:      "sleuth_project.myproject",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccResourceSleuth = `
resource "sleuth_project" "myproject" {
	name = "My project blah"
	description = "blah"
	impact_sensitivity = "FINE"
	failure_sensitivity = 500
}
`

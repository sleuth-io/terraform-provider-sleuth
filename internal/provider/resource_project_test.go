package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceProject(t *testing.T) {
	if err := testAccCheckOrganization(); err != nil {
		t.Skipf("Skipping because %s.", err.Error())
	}
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceProject,
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

const testAccResourceProject = `
resource "sleuth_project" "myproject" {
	name = "My project blah"
	impact_sensitivity = "FINE"
	failure_sensitivity = 500
}
`

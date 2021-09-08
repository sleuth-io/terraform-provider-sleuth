package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceEnvironment(t *testing.T) {
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
						"sleuth_environment.myenvironment", "name", regexp.MustCompile("^My environment blah")),
				),
			},
			{
				ResourceName:      "sleuth_environment.myenvironment",
				ImportState:       true,
				ImportStateVerify: false,
			},
		},
	})
}

const testAccResourceSleuth = `
resource "sleuth_project" "myproject" {
	name = "My project blah"
}

resource "sleuth_environment" "myenvironment" {
	project_slug = "${sleuth_project.myproject.id}"
	name = "My environment blah"
	description = "blah"
	color = "#ffffff"
}
`

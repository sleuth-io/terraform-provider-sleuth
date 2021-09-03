package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceSleuth(t *testing.T) {
	t.Skip("requires env vars against a running sleuth")

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceSleuth,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"sleuth_project.myproject", "name", regexp.MustCompile("^myproj")),
				),
			},
		},
	})
}

const testAccResourceSleuth = `
resource "sleuth_project" "myproject" {
	name = "myproject"
}
`

package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceCodeSourceSource(t *testing.T) {
	if err := testAccCheckOrganization(); err != nil {
		t.Skipf("Skipping because %s.", err.Error())
	}
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceCodeSourceSource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"sleuth_code_change_source.mysource", "name", regexp.MustCompile("^Sleuth Test")),
				),
			},
			{
				ResourceName:      "sleuth_code_change_source.mysource",
				ImportState:       true,
				ImportStateVerify: false,
			},
		},
	})
}

const testAccResourceCodeSourceSource = `
resource "sleuth_project" "myproject" {
	name = "My project blah"
}

resource "sleuth_environment" "myenv" {
	project_slug = "${sleuth_project.myproject.id}"
	name = "Production"
}

resource "sleuth_environment" "stg" {
	project_slug = "${sleuth_project.myproject.id}"
	name = "Staging"
}

resource "sleuth_code_change_source" "mysource" {
	project_slug = "${sleuth_project.myproject.id}"
	name = "Sleuth Test"
	repository {
		name = "sleuth-test"
		owner = "mrdon"
		provider = "GITHUB"
		url = "https://github.com/mrdon/sleuth-test"
	}
	environment_mappings {
		environment_slug = "${sleuth_environment.myenv.id}"
		branch = "master"
	}
	environment_mappings {
		environment_slug = "${sleuth_environment.stg.id}"
		branch = "master"
	}
	deploy_tracking_type = "auto_pr"
}
`

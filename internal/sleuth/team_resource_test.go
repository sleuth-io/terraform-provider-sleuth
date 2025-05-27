package sleuth

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAccTeamResource_v6(t *testing.T) {
	randomStr := acctest.RandStringFromCharSet(5, acctest.CharSetAlpha)
	teamName := fmt.Sprintf("Terraform test team %s", randomStr)
	updatedTeamName := fmt.Sprintf("Terraform test team updated %s", randomStr)
	parentName := fmt.Sprintf("Terraform parent team %s", randomStr)

	member1 := "dbrown@sleuth.io"
	member2 := "detkin@sleuth.io"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBetween(tfversion.Version0_14_0, tfversion.Version0_15_0),
		},
		Steps: []resource.TestStep{
			// Create team
			{
				Config: testAccTeamConfig(teamName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sleuth_team.terraform_acc_test", "name", teamName),
				),
			},
			// Create parent and subteam
			{
				Config: testAccTeamWithParentConfig(parentName, teamName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sleuth_team.terraform_acc_parent", "name", parentName),
					resource.TestCheckResourceAttrSet("sleuth_team.terraform_acc_test", "parent_slug"),
				),
			},
			// Add members (with parent_slug)
			{
				Config: testAccTeamWithMembersAndParentConfig(parentName, teamName, []string{member1, member2}),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sleuth_team.terraform_acc_test", "name", teamName),
					resource.TestCheckTypeSetElemAttr("sleuth_team.terraform_acc_test", "members.*", member1),
					resource.TestCheckTypeSetElemAttr("sleuth_team.terraform_acc_test", "members.*", member2),
					resource.TestCheckResourceAttrSet("sleuth_team.terraform_acc_test", "parent_slug"),
				),
			},
			// Remove a member and update name (with parent_slug)
			{
				Config: testAccTeamWithMembersAndParentConfig(parentName, updatedTeamName, []string{member2}),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("sleuth_team.terraform_acc_test", "name", updatedTeamName),
					resource.TestCheckTypeSetElemAttr("sleuth_team.terraform_acc_test", "members.*", member2),
					resource.TestCheckResourceAttrSet("sleuth_team.terraform_acc_test", "parent_slug"),
				),
			},
		},
	})
}

func testAccTeamConfig(name string) string {
	return fmt.Sprintf(`
resource "sleuth_team" "terraform_acc_test" {
  name = "%s"
}
`, name)
}

func testAccTeamWithParentConfig(parentName, subteamName string) string {
	return fmt.Sprintf(`
resource "sleuth_team" "terraform_acc_parent" {
  name = "%s"
}
resource "sleuth_team" "terraform_acc_test" {
  name = "%s"
  parent_slug = sleuth_team.terraform_acc_parent.slug
}
`, parentName, subteamName)
}

func testAccTeamWithMembersAndParentConfig(parentName, name string, members []string) string {
	membersStr := ""
	for _, m := range members {
		membersStr += fmt.Sprintf("  \"%s\",\n", m)
	}
	return fmt.Sprintf(`
resource "sleuth_team" "terraform_acc_parent" {
  name = "%s"
}
resource "sleuth_team" "terraform_acc_test" {
  name = "%s"
  parent_slug = sleuth_team.terraform_acc_parent.slug
  members = [
%s  ]
}
`, parentName, name, membersStr)
}

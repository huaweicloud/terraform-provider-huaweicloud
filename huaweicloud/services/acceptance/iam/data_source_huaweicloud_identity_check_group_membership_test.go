package iam

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataCheckGroupMembership_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceName()

		check = "data.huaweicloud_identity_check_group_membership.test"
		dc    = acceptance.InitDataSourceCheck(check)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source: "hashicorp/random",
				// The version of the random provider must be greater than 3.3.0 to support the 'numeric' parameter.
				VersionConstraint: "3.3.0",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccDataCheckGroupMembership_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(check, "result", "false"),
				),
			},
			{
				Config: testAccDataCheckGroupMembership_basic_step2(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(check, "result", "true"),
				),
			},
		},
	})
}

func testAccDataCheckGroupMembership_base(name string) string {
	return fmt.Sprintf(`
resource "random_string" "test" {
  length           = 10
  min_numeric      = 1
  min_special      = 1
  min_lower        = 1
  override_special = "@!"
}

resource "huaweicloud_identity_group" "test" {
  name = "%[1]s"
}

resource "huaweicloud_identity_user" "test" {
  name     = "%[1]s"
  password = random_string.test.result
  enabled  = true
}
`, name)
}

func testAccDataCheckGroupMembership_basic_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_identity_check_group_membership" "test" {
  group_id = huaweicloud_identity_group.test.id
  user_id  = huaweicloud_identity_user.test.id
}
`, testAccDataCheckGroupMembership_base(name))
}

func testAccDataCheckGroupMembership_basic_step2(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_identity_group_membership" "test" {
  group = huaweicloud_identity_group.test.id
  users = [huaweicloud_identity_user.test.id]
}

data "huaweicloud_identity_check_group_membership" "test" {
  group_id = huaweicloud_identity_group.test.id
  user_id  = huaweicloud_identity_user.test.id

  depends_on = [huaweicloud_identity_group_membership.test]
}
`, testAccDataCheckGroupMembership_base(name))
}

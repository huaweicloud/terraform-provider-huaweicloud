package iam

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataCheckGroupMembership_basic(t *testing.T) {
	var (
		groupName     = acceptance.RandomAccResourceName()
		userName      = acceptance.RandomAccResourceName()
		otherUserName = acceptance.RandomAccResourceName()
		password      = acceptance.RandomPassword()

		check = "data.huaweicloud_identity_check_group_membership.test"
		dc    = acceptance.InitDataSourceCheck(check)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataCheckGroupMembership_basic_step1(groupName, userName, password),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(check, "result", "true"),
				),
			},
			{
				Config: testAccDataCheckGroupMembership_basic_step2(groupName, otherUserName, password),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(check, "result", "false"),
				),
			},
		},
	})
}

func testAccDataCheckGroupMembership_base(groupName, userName, password string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_group" "test" {
  name = "%[1]s"
}

resource "huaweicloud_identity_user" "test" {
  name     = "%[2]s"
  password = "%[3]s"
  enabled  = true
}
`, groupName, userName, password)
}

func testAccDataCheckGroupMembership_basic_step1(groupName, userName, password string) string {
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
`, testAccDataCheckGroupMembership_base(groupName, userName, password))
}

func testAccDataCheckGroupMembership_basic_step2(groupName, userName, password string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_identity_check_group_membership" "test" {
  group_id = huaweicloud_identity_group.test.id
  user_id  = huaweicloud_identity_user.test.id
}
`, testAccDataCheckGroupMembership_base(groupName, userName, password))
}

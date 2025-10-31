package iam

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccIdentityCheckGroupMembership_basic(t *testing.T) {
	var (
		groupName      = acceptance.RandomAccResourceName()
		userName1      = acceptance.RandomAccResourceName()
		userName2      = acceptance.RandomAccResourceName()
		password       = acceptance.RandomPassword()
		dataSourceName = "data.huaweicloud_identity_check_group_membership.test"
	)
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityCheckGroupMembership(groupName, userName1, password),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "result", "true"),
				),
			},
			{
				Config: testAccIdentityCheckGroupMembershipNot(groupName, userName2, password),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "result", "false"),
				),
			},
		},
	})
}

func testAccIdentityCheckGroupMembership(groupName, userName, password string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_group" "group_1" {
  name = "%s"
}

resource "huaweicloud_identity_user" "user_1" {
  name     = "%s"
  password = "%s"
  enabled  = true
}
   
resource "huaweicloud_identity_group_membership" "membership_1" {
  group = huaweicloud_identity_group.group_1.id
  users = [huaweicloud_identity_user.user_1.id]
}

data "huaweicloud_identity_check_group_membership" "test" {
  group_id = huaweicloud_identity_group.group_1.id
  user_id  = huaweicloud_identity_user.user_1.id

  depends_on = [huaweicloud_identity_group_membership.membership_1]
}
`, groupName, userName, password)
}

func testAccIdentityCheckGroupMembershipNot(groupName, userName, password string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_group" "group_1" {
  name = "%s"
}

resource "huaweicloud_identity_user" "user_1" {
  name     = "%s"
  password = "%s"
  enabled  = true
}

data "huaweicloud_identity_check_group_membership" "test" {
  group_id = huaweicloud_identity_group.group_1.id
  user_id  = huaweicloud_identity_user.user_1.id
}
`, groupName, userName, password)
}

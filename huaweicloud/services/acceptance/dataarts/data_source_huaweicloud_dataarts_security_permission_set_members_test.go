package dataarts

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSecurityPermissionSetMembers_basic(t *testing.T) {
	var (
		rName = acceptance.RandomAccResourceName()

		all = "data.huaweicloud_dataarts_security_permission_set_members.all"
		dc  = acceptance.InitDataSourceCheck(all)

		byMemberName   = "data.huaweicloud_dataarts_security_permission_set_members.filter_by_member_name"
		dcByMemberName = acceptance.InitDataSourceCheck(byMemberName)

		byMemberType   = "data.huaweicloud_dataarts_security_permission_set_members.filter_by_member_type"
		dcByMemberType = acceptance.InitDataSourceCheck(byMemberType)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDataArtsWorkSpaceID(t)
			acceptance.TestAccPreCheckDataArtsManagerID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSecurityPermissionSetMembers_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(all, "workspace_id", acceptance.HW_DATAARTS_WORKSPACE_ID),
					resource.TestCheckResourceAttrSet(all, "permission_set_id"),
					resource.TestCheckResourceAttrSet(all, "region"),
					resource.TestMatchResourceAttr(all, "members.#", regexp.MustCompile(`[1-9]([0-9]*)?`)),
					resource.TestCheckResourceAttrSet(all, "members.0.id"),
					resource.TestCheckResourceAttrSet(all, "members.0.name"),
					resource.TestCheckResourceAttrSet(all, "members.0.type"),
					// Filter by member_name
					dcByMemberName.CheckResourceExists(),
					resource.TestCheckOutput("is_member_name_filter_useful", "true"),
					// Filter by member_type
					dcByMemberType.CheckResourceExists(),
					resource.TestCheckOutput("is_member_type_filter_useful", "true"),
					resource.TestCheckResourceAttr(byMemberType, "members.0.type", "USER"),
				),
			},
		},
	})
}

func testAccDataSecurityPermissionSetMembers_basic_base(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_dataarts_studio_workspace_user_roles" "test" {
  workspace_id = "%[1]s"
}

data "huaweicloud_identity_users" "test" {
  name = "%[2]s"
}

resource "huaweicloud_dataarts_studio_workspace_user" "test" {
  workspace_id = "%[1]s"
  user_id      = try(data.huaweicloud_identity_users.test.users[0].id, "NOT_FOUND")

  dynamic "roles" {
    for_each = slice(data.huaweicloud_dataarts_studio_workspace_user_roles.test.roles, 0, 1)

	content {
      id = roles.value.id
    }
  }
}

resource "huaweicloud_dataarts_security_permission_set" "test" {
  workspace_id = "%[1]s"
  name         = "%[3]s"
  parent_id    = "0"
  manager_id   = "%[4]s"
}

resource "huaweicloud_dataarts_security_permission_set_member" "test" {
  workspace_id      = "%[1]s"
  permission_set_id = huaweicloud_dataarts_security_permission_set.test.id
  object_id         = huaweicloud_dataarts_studio_workspace_user.test.id
  name              = "%[3]s"
  type              = "USER"
}
`, acceptance.HW_DATAARTS_WORKSPACE_ID, acceptance.HW_USER_NAME, name, acceptance.HW_DATAARTS_MANAGER_ID)
}

func testAccDataSecurityPermissionSetMembers_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

# Query all members and without any filter
data "huaweicloud_dataarts_security_permission_set_members" "all" {
  depends_on = [
    huaweicloud_dataarts_security_permission_set_member.test,
  ]

  workspace_id      = "%[2]s"
  permission_set_id = huaweicloud_dataarts_security_permission_set.test.id
}

# Filter by member_name
locals {
  member_name = huaweicloud_dataarts_security_permission_set_member.test.name
}

data "huaweicloud_dataarts_security_permission_set_members" "filter_by_member_name" {
  # The behavior of parameter 'member_name' of the data source is 'Required', means this parameter does not
  # have 'Know After Apply' behavior.
  depends_on = [
    huaweicloud_dataarts_security_permission_set_member.test,
  ]

  workspace_id      = "%[2]s"
  permission_set_id = huaweicloud_dataarts_security_permission_set.test.id
  member_name       = local.member_name
}

locals {
  member_name_filter_result = [
    for v in data.huaweicloud_dataarts_security_permission_set_members.filter_by_member_name.members[*].name :
      v == local.member_name
  ]
}

output "is_member_name_filter_useful" {
  value = length(local.member_name_filter_result) > 0 && alltrue(local.member_name_filter_result)
}

# Filter by member_type
locals {
  member_type = "USER"
}

data "huaweicloud_dataarts_security_permission_set_members" "filter_by_member_type" {
  depends_on = [
    huaweicloud_dataarts_security_permission_set_member.test,
  ]

  workspace_id      = "%[2]s"
  permission_set_id = huaweicloud_dataarts_security_permission_set.test.id
  member_type       = local.member_type
}

locals {
  member_type_filter_result = [
    for v in data.huaweicloud_dataarts_security_permission_set_members.filter_by_member_type.members[*].type :
      v == local.member_type
  ]
}

output "is_member_type_filter_useful" {
  value = length(local.member_type_filter_result) > 0 && alltrue(local.member_type_filter_result)
}
`, testAccDataSecurityPermissionSetMembers_basic_base(name), acceptance.HW_DATAARTS_WORKSPACE_ID)
}

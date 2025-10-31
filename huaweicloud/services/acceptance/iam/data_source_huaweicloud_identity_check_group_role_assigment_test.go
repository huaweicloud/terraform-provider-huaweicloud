package iam

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccIdentityCheckGroupRoleAssignment_basic(t *testing.T) {
	var (
		groupName1           = acceptance.RandomAccResourceName()
		dataSourceDomain     = "data.huaweicloud_identity_check_group_role_assignment.test_domain"
		dataSourceProject    = "data.huaweicloud_identity_check_group_role_assignment.test_project"
		dataSourceProjectAll = "data.huaweicloud_identity_check_group_role_assignment.test_project_all"
		dataSourceDomainNot  = "data.huaweicloud_identity_check_group_role_assignment.test_domain_not"
	)
	dcDomain := acceptance.InitDataSourceCheck(dataSourceDomain)
	dcProject := acceptance.InitDataSourceCheck(dataSourceProject)
	dcProjectAll := acceptance.InitDataSourceCheck(dataSourceProjectAll)
	dcDomainNot := acceptance.InitDataSourceCheck(dataSourceDomainNot)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityCheckGroupRoleAssignment(groupName1),
				Check: resource.ComposeTestCheckFunc(
					dcDomain.CheckResourceExists(),
					dcProject.CheckResourceExists(),
					dcProjectAll.CheckResourceExists(),
					dcDomainNot.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceDomain, "result", "true"),
					resource.TestCheckResourceAttr(dataSourceProject, "result", "true"),
					resource.TestCheckResourceAttr(dataSourceProjectAll, "result", "true"),
					resource.TestCheckResourceAttr(dataSourceDomainNot, "result", "false"),
				),
			},
		},
	})
}

func testAccIdentityCheckGroupRoleAssignment(groupName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_group" "test" {
  name = "%s"
}

data "huaweicloud_identity_role" "role_1" {
  display_name = "AAD FullAccess"
}

data "huaweicloud_identity_role" "role_2" {
  display_name = "AAD ReadOnlyAccess"
}

data "huaweicloud_identity_role" "role_3" {
  display_name = "AND ReadOnlyAccess"
}

data "huaweicloud_identity_role" "role_4" {
  display_name = "APIG Administrator"
}

resource "huaweicloud_identity_group_role_assignment" "test_domain" {
  group_id  = huaweicloud_identity_group.test.id
  role_id   = data.huaweicloud_identity_role.role_1.id
  domain_id = "%s"
}

resource "huaweicloud_identity_group_role_assignment" "test_project" {
  group_id   = huaweicloud_identity_group.test.id
  role_id    = data.huaweicloud_identity_role.role_2.id
  project_id = "%s"
}

resource "huaweicloud_identity_group_role_assignment" "test_project_all" {
  group_id   = huaweicloud_identity_group.test.id
  role_id    = data.huaweicloud_identity_role.role_3.id
  project_id = "all"
}

data "huaweicloud_identity_check_group_role_assignment" "test_domain" {
  group_id  = huaweicloud_identity_group.test.id
  role_id   = data.huaweicloud_identity_role.role_1.id
  domain_id = "%s"

  depends_on = [huaweicloud_identity_group_role_assignment.test_domain]
}

data "huaweicloud_identity_check_group_role_assignment" "test_project" {
  group_id   = huaweicloud_identity_group.test.id
  role_id    = data.huaweicloud_identity_role.role_2.id
  project_id = "%s"

  depends_on = [huaweicloud_identity_group_role_assignment.test_project]
}

data "huaweicloud_identity_check_group_role_assignment" "test_project_all" {
  group_id   = huaweicloud_identity_group.test.id
  role_id    = data.huaweicloud_identity_role.role_3.id
  project_id = "all"

  depends_on = [huaweicloud_identity_group_role_assignment.test_project_all]
}

data "huaweicloud_identity_check_group_role_assignment" "test_domain_not" {
  group_id  = huaweicloud_identity_group.test.id
  role_id   = data.huaweicloud_identity_role.role_4.id
  domain_id = "%s"
}
`, groupName, acceptance.HW_DOMAIN_ID, acceptance.HW_PROJECT_ID,
		acceptance.HW_DOMAIN_ID, acceptance.HW_PROJECT_ID, acceptance.HW_DOMAIN_ID)
}

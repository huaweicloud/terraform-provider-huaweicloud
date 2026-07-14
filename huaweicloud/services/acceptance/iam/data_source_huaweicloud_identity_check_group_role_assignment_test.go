package iam

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataCheckGroupRoleAssignment_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceName()

		byDomain   = "data.huaweicloud_identity_check_group_role_assignment.filter_by_domain"
		dcByDomain = acceptance.InitDataSourceCheck(byDomain)

		byProject   = "data.huaweicloud_identity_check_group_role_assignment.filter_by_project"
		dcByProject = acceptance.InitDataSourceCheck(byProject)

		byProjectAll   = "data.huaweicloud_identity_check_group_role_assignment.filter_by_project_all"
		dcByProjectAll = acceptance.InitDataSourceCheck(byProjectAll)

		byDomainNot   = "data.huaweicloud_identity_check_group_role_assignment.filter_by_domain_not"
		dcByDomainNot = acceptance.InitDataSourceCheck(byDomainNot)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAdminOnly(t)
			acceptance.TestAccPrecheckDomainId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataCheckGroupRoleAssignment_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dcByDomain.CheckResourceExists(),
					resource.TestCheckResourceAttr(byDomain, "result", "true"),
					dcByProject.CheckResourceExists(),
					resource.TestCheckResourceAttr(byProject, "result", "true"),
					dcByProjectAll.CheckResourceExists(),
					resource.TestCheckResourceAttr(byProjectAll, "result", "true"),
					dcByDomainNot.CheckResourceExists(),
					resource.TestCheckResourceAttr(byDomainNot, "result", "false"),
				),
			},
		},
	})
}

func testAccDataCheckGroupRoleAssignment_base(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_identity_projects" "test" {
  name = "%[1]s"
}

resource "huaweicloud_identity_group" "test" {
  name = "%[2]s"
}

data "huaweicloud_identity_role" "aad_full_access" {
  display_name = "AAD FullAccess"
}

data "huaweicloud_identity_role" "aad_read_only_access" {
  display_name = "AAD ReadOnlyAccess"
}

data "huaweicloud_identity_role" "api_administrator" {
  display_name = "APIG Administrator"
}

resource "huaweicloud_identity_group_role_assignment" "test_with_domain" {
  group_id  = huaweicloud_identity_group.test.id
  role_id   = data.huaweicloud_identity_role.aad_full_access.id
  domain_id = "%[3]s"
}

resource "huaweicloud_identity_group_role_assignment" "test_with_project" {
  group_id   = huaweicloud_identity_group.test.id
  role_id    = data.huaweicloud_identity_role.aad_read_only_access.id
  project_id = try(data.huaweicloud_identity_projects.test.projects[0].id, "NOT_FOUND")
}

resource "huaweicloud_identity_group_role_assignment" "test_with_project_all" {
  group_id   = huaweicloud_identity_group.test.id
  role_id    = data.huaweicloud_identity_role.api_administrator.id
  project_id = "all"
}
`, acceptance.HW_REGION_NAME, name, acceptance.HW_DOMAIN_ID)
}

func testAccDataCheckGroupRoleAssignment_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_identity_check_group_role_assignment" "filter_by_domain" {
  group_id  = huaweicloud_identity_group.test.id
  role_id   = data.huaweicloud_identity_role.aad_full_access.id
  domain_id = "%[2]s"

  depends_on = [huaweicloud_identity_group_role_assignment.test_with_domain]
}

data "huaweicloud_identity_check_group_role_assignment" "filter_by_project" {
  group_id   = huaweicloud_identity_group.test.id
  role_id    = data.huaweicloud_identity_role.aad_read_only_access.id
  project_id = try(data.huaweicloud_identity_projects.test.projects[0].id, "NOT_FOUND")

  depends_on = [huaweicloud_identity_group_role_assignment.test_with_project]
}

data "huaweicloud_identity_check_group_role_assignment" "filter_by_project_all" {
  group_id   = huaweicloud_identity_group.test.id
  role_id    = data.huaweicloud_identity_role.api_administrator.id
  project_id = "all"

  depends_on = [huaweicloud_identity_group_role_assignment.test_with_project_all]
}

data "huaweicloud_identity_check_group_role_assignment" "filter_by_domain_not" {
  group_id  = huaweicloud_identity_group.test.id
  role_id   = data.huaweicloud_identity_role.aad_full_access.id
  domain_id = "%[2]s"
}
`, testAccDataCheckGroupRoleAssignment_base(name), acceptance.HW_DOMAIN_ID)
}

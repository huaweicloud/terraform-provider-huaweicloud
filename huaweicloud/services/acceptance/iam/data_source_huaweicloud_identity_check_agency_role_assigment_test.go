package iam

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataCheckAgencyRoleAssignment_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceName()

		byDomain   = "data.huaweicloud_identity_check_agency_role_assignment.filter_by_domain"
		dcByDomain = acceptance.InitDataSourceCheck(byDomain)

		byProject   = "data.huaweicloud_identity_check_agency_role_assignment.filter_by_project"
		dcByProject = acceptance.InitDataSourceCheck(byProject)

		byProjectAll   = "data.huaweicloud_identity_check_agency_role_assignment.filter_by_project_all"
		dcByProjectAll = acceptance.InitDataSourceCheck(byProjectAll)

		byDomainNot   = "data.huaweicloud_identity_check_agency_role_assignment.filter_by_domain_not"
		dcByDomainNot = acceptance.InitDataSourceCheck(byDomainNot)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAdminOnly(t)
			acceptance.TestAccPreCheckProjectID(t)
			acceptance.TestAccPrecheckDomainId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataCheckAgencyRoleAssignment_basic_step1(name),
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

func testAccDataCheckAgencyRoleAssignment_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_agency" "test" {
  name                  = "%[1]s"
  description           = "Test by terraform"
  delegated_domain_name = "%[2]s"

  project_role {
    project = "%[3]s"
    roles   = ["CCE Administrator"]
  }

  domain_roles = [
    "Server Administrator",
    "Anti-DDoS Administrator",
  ]

  all_resources_roles = [
    "VPC Administrator"
  ]
}

data "huaweicloud_identity_role" "server_administrator" {
  display_name = "Server Administrator"
}

data "huaweicloud_identity_role" "cce_administrator" {
  display_name = "CCE Administrator"
}

data "huaweicloud_identity_role" "vpc_administrator" {
  display_name = "VPC Administrator"
}

data "huaweicloud_identity_role" "aad_fullaccess" {
  display_name = "AAD FullAccess"
}
`, name, acceptance.HW_DOMAIN_NAME, acceptance.HW_REGION_NAME)
}

func testAccDataCheckAgencyRoleAssignment_basic_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_identity_check_agency_role_assignment" "filter_by_domain" {
  agency_id = huaweicloud_identity_agency.test.id
  role_id   = data.huaweicloud_identity_role.server_administrator.id
  domain_id = "%[2]s"
}

data "huaweicloud_identity_check_agency_role_assignment" "filter_by_project" {
  agency_id  = huaweicloud_identity_agency.test.id
  role_id    = data.huaweicloud_identity_role.cce_administrator.id
  project_id = "%[3]s"
}

data "huaweicloud_identity_check_agency_role_assignment" "filter_by_project_all" {
  agency_id  = huaweicloud_identity_agency.test.id
  role_id    = data.huaweicloud_identity_role.vpc_administrator.id
  project_id = "all"
}

data "huaweicloud_identity_check_agency_role_assignment" "filter_by_domain_not" {
  agency_id = huaweicloud_identity_agency.test.id
  role_id   = data.huaweicloud_identity_role.aad_fullaccess.id
  domain_id = "%[2]s"
}
`, testAccDataCheckAgencyRoleAssignment_base(name), acceptance.HW_DOMAIN_ID, acceptance.HW_PROJECT_ID)
}

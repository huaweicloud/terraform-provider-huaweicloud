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
data "huaweicloud_identity_projects" "test" {
  name = "%[1]s"
}

resource "huaweicloud_identity_agency" "test" {
  name                  = "%[2]s"
  description           = "Test by terraform"
  delegated_domain_name = "%[3]s"

  project_role {
    project = try(data.huaweicloud_identity_projects.test.projects[0].id, "NOT_FOUND")
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

variable "role_names" {
  type = list(string)
  default = [
    "Server Administrator",
    "CCE Administrator",
    "VPC Administrator",
    "AAD FullAccess",
  ]
}

data "huaweicloud_identity_role" "test" {
  count = length(var.role_names)

  display_name = var.role_names[count.index]
}
`, acceptance.HW_REGION_NAME, name, acceptance.HW_DOMAIN_NAME)
}

func testAccDataCheckAgencyRoleAssignment_basic_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_identity_check_agency_role_assignment" "filter_by_domain" {
  agency_id = huaweicloud_identity_agency.test.id
  role_id   = data.huaweicloud_identity_role.test[0].id
  domain_id = "%[2]s"
}

data "huaweicloud_identity_check_agency_role_assignment" "filter_by_project" {
  agency_id  = huaweicloud_identity_agency.test.id
  role_id    = data.huaweicloud_identity_role.test[1].id
  project_id = try(data.huaweicloud_identity_projects.test.projects[0].id, "NOT_FOUND")
}

data "huaweicloud_identity_check_agency_role_assignment" "filter_by_project_all" {
  agency_id  = huaweicloud_identity_agency.test.id
  role_id    = data.huaweicloud_identity_role.test[2].id
  project_id = "all"
}

data "huaweicloud_identity_check_agency_role_assignment" "filter_by_domain_not" {
  agency_id = huaweicloud_identity_agency.test.id
  role_id   = data.huaweicloud_identity_role.test[3].id
  domain_id = "%[2]s"
}
`, testAccDataCheckAgencyRoleAssignment_base(name), acceptance.HW_DOMAIN_ID)
}

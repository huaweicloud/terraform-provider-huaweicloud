package iam

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccIdentityCheckAgencyRoleAssignment_basic(t *testing.T) {
	var (
		agencyName1          = acceptance.RandomAccResourceName()
		dataSourceDomain     = "data.huaweicloud_identity_check_agency_role_assignment.test_domain"
		dataSourceProject    = "data.huaweicloud_identity_check_agency_role_assignment.test_project"
		dataSourceProjectAll = "data.huaweicloud_identity_check_agency_role_assignment.test_project_all"
		dataSourceDomainNot  = "data.huaweicloud_identity_check_agency_role_assignment.test_domain_not"
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
				Config: testAccIdentityCheckAgencyRoleAssignmentDomain(agencyName1),
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

func testAccIdentityCheckAgencyRoleAssignmentDomain(agencyName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_agency" "test" {
  name                  = "%s"
  description           = "This is a test agency"
  delegated_domain_name = "%s"

  project_role {
    project = "%s"
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

data "huaweicloud_identity_role" "role_1" {
  display_name = "Server Administrator"
}

data "huaweicloud_identity_role" "role_2" {
  display_name = "CCE Administrator"
}

data "huaweicloud_identity_role" "role_3" {
  display_name = "VPC Administrator"
}

data "huaweicloud_identity_role" "role_4" {
  display_name = "AAD FullAccess"
}

data "huaweicloud_identity_check_agency_role_assignment" "test_domain" {
  agency_id = huaweicloud_identity_agency.test.id
  role_id   = data.huaweicloud_identity_role.role_1.id
  domain_id = "%s"
}

data "huaweicloud_identity_check_agency_role_assignment" "test_project" {
  agency_id  = huaweicloud_identity_agency.test.id
  role_id    = data.huaweicloud_identity_role.role_2.id
  project_id = "%s"
}

data "huaweicloud_identity_check_agency_role_assignment" "test_project_all" {
  agency_id  = huaweicloud_identity_agency.test.id
  role_id    = data.huaweicloud_identity_role.role_3.id
  project_id = "all"
}

data "huaweicloud_identity_check_agency_role_assignment" "test_domain_not" {
  agency_id = huaweicloud_identity_agency.test.id
  role_id   = data.huaweicloud_identity_role.role_4.id
  domain_id = "%s"
}
`, agencyName, acceptance.HW_DOMAIN_NAME, acceptance.HW_REGION_NAME,
		acceptance.HW_DOMAIN_ID, acceptance.HW_PROJECT_ID, acceptance.HW_DOMAIN_ID)
}

package iam

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/identity/v3/agency"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getAgencyResourceFunc(c *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := c.IAMV3Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating IAM client: %s", err)
	}
	return agency.Get(client, state.Primary.ID).Extract()
}

func TestAccAgency_basic(t *testing.T) {
	var (
		obj interface{}

		createWithRoleAssignments = "huaweicloud_identity_agency.create_with_role_assignments"
		rcWithRoleAssignments     = acceptance.InitResourceCheck(createWithRoleAssignments, &obj, getAgencyResourceFunc)

		createWithoutRoleAssignments = "huaweicloud_identity_agency.create_without_role_assignments"
		rcWithoutRoleAssignments     = acceptance.InitResourceCheck(createWithoutRoleAssignments, &obj, getAgencyResourceFunc)

		name = acceptance.RandomAccResourceName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAdminOnly(t)
			acceptance.TestAccPreCheckEpsID(t)
			acceptance.TestAccPrecheckDomainName(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			rcWithRoleAssignments.CheckResourceDestroy(),
			rcWithoutRoleAssignments.CheckResourceDestroy(),
		),
		Steps: []resource.TestStep{
			{
				Config: testAccAgency_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rcWithRoleAssignments.CheckResourceExists(),
					resource.TestCheckResourceAttr(createWithRoleAssignments, "name", name+"_create_with_role_assignments"),
					resource.TestCheckResourceAttr(createWithRoleAssignments, "description", "Created by terraform acceptance test"),
					resource.TestCheckResourceAttr(createWithRoleAssignments, "delegated_domain_name", acceptance.HW_DOMAIN_NAME),
					resource.TestCheckResourceAttr(createWithRoleAssignments, "duration", "FOREVER"),
					resource.TestCheckResourceAttr(createWithRoleAssignments, "project_role.#", "1"),
					resource.TestCheckResourceAttr(createWithRoleAssignments, "domain_roles.#", "2"),
					resource.TestCheckResourceAttr(createWithRoleAssignments, "all_resources_roles.#", "1"),
					resource.TestCheckResourceAttr(createWithRoleAssignments, "enterprise_project_roles.#", "1"),
					resource.TestMatchResourceAttr(createWithRoleAssignments, "create_time",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|(\.\d{6}))$`)),
					rcWithoutRoleAssignments.CheckResourceExists(),
					resource.TestCheckResourceAttr(createWithoutRoleAssignments, "name", name+"_create_without_role_assignments"),
					resource.TestCheckResourceAttr(createWithoutRoleAssignments, "delegated_domain_name", acceptance.HW_DOMAIN_NAME),
					resource.TestCheckResourceAttr(createWithoutRoleAssignments, "duration", "FOREVER"),
					resource.TestCheckResourceAttr(createWithoutRoleAssignments, "project_role.#", "0"),
					resource.TestCheckResourceAttr(createWithoutRoleAssignments, "domain_roles.#", "0"),
					resource.TestCheckResourceAttr(createWithoutRoleAssignments, "all_resources_roles.#", "0"),
					resource.TestCheckResourceAttr(createWithoutRoleAssignments, "enterprise_project_roles.#", "0"),
					resource.TestMatchResourceAttr(createWithoutRoleAssignments, "create_time",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|(\.\d{6}))$`)),
				),
			},
			{
				Config: testAccAgency_basic_step2(name),
				Check: resource.ComposeTestCheckFunc(
					rcWithRoleAssignments.CheckResourceExists(),
					resource.TestCheckResourceAttr(createWithRoleAssignments, "name", name+"_create_with_role_assignments"),
					resource.TestCheckResourceAttr(createWithRoleAssignments, "description", ""),
					resource.TestCheckResourceAttr(createWithRoleAssignments, "delegated_domain_name", acceptance.HW_DOMAIN_NAME),
					resource.TestCheckResourceAttr(createWithRoleAssignments, "duration", "1"),
					resource.TestCheckResourceAttr(createWithRoleAssignments, "project_role.#", "1"),
					resource.TestCheckResourceAttr(createWithRoleAssignments, "domain_roles.#", "3"),
					resource.TestCheckResourceAttr(createWithRoleAssignments, "all_resources_roles.#", "1"),
					resource.TestCheckResourceAttr(createWithRoleAssignments, "enterprise_project_roles.#", "1"),
					rcWithoutRoleAssignments.CheckResourceExists(),
					resource.TestCheckResourceAttr(createWithoutRoleAssignments, "name", name+"_create_without_role_assignments"),
					resource.TestCheckResourceAttr(createWithoutRoleAssignments, "description", "Updated by terraform acceptance test"),
					resource.TestCheckResourceAttr(createWithoutRoleAssignments, "delegated_domain_name", acceptance.HW_DOMAIN_NAME),
					resource.TestCheckResourceAttr(createWithoutRoleAssignments, "duration", "30"),
					resource.TestCheckResourceAttr(createWithoutRoleAssignments, "project_role.#", "1"),
					resource.TestCheckResourceAttr(createWithoutRoleAssignments, "domain_roles.#", "2"),
					resource.TestCheckResourceAttr(createWithoutRoleAssignments, "all_resources_roles.#", "1"),
					resource.TestCheckResourceAttr(createWithoutRoleAssignments, "enterprise_project_roles.#", "1"),
				),
			},
			{
				ResourceName:      createWithRoleAssignments,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"all_resources_roles",
					"enterprise_project_roles",
				},
			},
			{
				ResourceName:      createWithoutRoleAssignments,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"all_resources_roles",
					"enterprise_project_roles",
				},
			},
		},
	})
}

func testAccAgency_basic_step1(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_enterprise_projects" "test" {
  enterprise_project_id = "%[1]s"
}

resource "huaweicloud_identity_agency" "create_with_role_assignments" {
  name                  = "%[2]s_create_with_role_assignments"
  description           = "Created by terraform acceptance test"
  delegated_domain_name = "%[3]s"

  project_role {
    project = "%[4]s"
    roles   = ["CCE Administrator"]
  }

  domain_roles = [
    "Server Administrator",
    "Anti-DDoS Administrator",
  ]

  all_resources_roles = [
    "VPC Administrator"
  ]

  enterprise_project_roles {
    enterprise_project = try(data.huaweicloud_enterprise_projects.test.enterprise_projects[0].name, "NOT_FOUND")
    roles              = ["CCE ReadOnlyAccess"]
  }
}

resource "huaweicloud_identity_agency" "create_without_role_assignments" {
  name                  = "%[2]s_create_without_role_assignments"
  delegated_domain_name = "%[3]s"
  duration              = "FOREVER"

  domain_roles = []

  all_resources_roles = []
}
`, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST,
		name,
		acceptance.HW_DOMAIN_NAME,
		acceptance.HW_REGION_NAME,
	)
}

func testAccAgency_basic_step2(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_enterprise_projects" "test" {
  enterprise_project_id = "%[1]s"
}

resource "huaweicloud_identity_agency" "create_with_role_assignments" {
  name                  = "%[2]s_create_with_role_assignments"
  delegated_domain_name = "%[3]s"
  duration              = "1"

  project_role {
    project = "%[4]s"
    roles   = ["RDS Administrator"]
  }

  domain_roles = [
    "Anti-DDoS Administrator",
    "SMN Administrator",
    "OBS Administrator",
  ]

  all_resources_roles = [
    "VPCEndpoint Administrator"
  ]

  enterprise_project_roles {
    enterprise_project = try(data.huaweicloud_enterprise_projects.test.enterprise_projects[0].name, "NOT_FOUND")
    roles              = ["RDS ReadOnlyAccess"]
  }
}

resource "huaweicloud_identity_agency" "create_without_role_assignments" {
  name                  = "%[2]s_create_without_role_assignments"
  description           = "Updated by terraform acceptance test"
  delegated_domain_name = "%[3]s"
  duration              = "30"

  project_role {
    project = "%[4]s"
    roles   = ["CCE Administrator"]
  }

  domain_roles = [
    "Server Administrator",
    "Anti-DDoS Administrator",
  ]

  all_resources_roles = [
    "VPC Administrator"
  ]

  enterprise_project_roles {
    enterprise_project = try(data.huaweicloud_enterprise_projects.test.enterprise_projects[0].name, "NOT_FOUND")
    roles              = ["CCE ReadOnlyAccess"]
  }
}
`, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST,
		name,
		acceptance.HW_DOMAIN_NAME,
		acceptance.HW_REGION_NAME,
	)
}

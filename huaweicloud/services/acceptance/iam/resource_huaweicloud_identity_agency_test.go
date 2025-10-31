package iam

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/identity/v3/agency"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getIdentityAgencyResourceFunc(c *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := c.IAMV3Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating IAM client: %s", err)
	}
	return agency.Get(client, state.Primary.ID).Extract()
}

func TestAccIdentityAgency_basic(t *testing.T) {
	var object agency.Agency
	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_identity_agency.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&object,
		getIdentityAgencyResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAdminOnly(t)
			acceptance.TestAccPrecheckDomainName(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityAgency_domain(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "This is a test agency"),
					resource.TestCheckResourceAttr(resourceName, "delegated_domain_name", acceptance.HW_DOMAIN_NAME),
					resource.TestCheckResourceAttr(resourceName, "duration", "FOREVER"),
					resource.TestCheckResourceAttr(resourceName, "project_role.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "domain_roles.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "all_resources_roles.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_roles.#", "1"),
				),
			},
			{
				Config: testAccIdentityAgency_domainUpdate(rName, "ONEDAY"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "This is a updated test agency"),
					resource.TestCheckResourceAttr(resourceName, "delegated_domain_name", acceptance.HW_DOMAIN_NAME),
					resource.TestCheckResourceAttr(resourceName, "duration", "1"),
					resource.TestCheckResourceAttr(resourceName, "project_role.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "domain_roles.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "all_resources_roles.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_roles.#", "1"),
				),
			},
			{
				Config: testAccIdentityAgency_domainUpdate(rName, "30"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "duration", "30"),
					resource.TestCheckResourceAttr(resourceName, "project_role.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "domain_roles.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "all_resources_roles.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_roles.#", "1"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"all_resources_roles", "enterprise_project_roles"},
			},
		},
	})
}

func testAccIdentityAgency_domain(rName string) string {
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
  enterprise_project_roles {
    enterprise_project = "%s"
    roles              = ["CCE ReadOnlyAccess"]
  }
}
`, rName, acceptance.HW_DOMAIN_NAME, acceptance.HW_REGION_NAME, acceptance.HW_ENTERPRISE_PROJECT_NAME)
}

func testAccIdentityAgency_domainUpdate(rName, duration string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_agency" "test" {
  name                  = "%s"
  description           = "This is a updated test agency"
  duration              = "%s"
  delegated_domain_name = "%s"

  project_role {
    project = "%s"
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
    enterprise_project = "%s"
    roles              = ["RDS ReadOnlyAccess"]
  }
}
`, rName, duration, acceptance.HW_DOMAIN_NAME, acceptance.HW_REGION_NAME, acceptance.HW_ENTERPRISE_PROJECT_NAME)
}

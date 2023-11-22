package swr

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/swr/v2/namespaces"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getResourcePermissions(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	swrClient, err := conf.SwrV2Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating SWR client: %s", err)
	}

	return namespaces.GetAccess(swrClient, state.Primary.ID).Extract()
}

func TestAccSwrOrganizationPermissions_basic(t *testing.T) {
	var permissions namespaces.Access
	organizationName := acceptance.RandomAccResourceName()
	userName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_swr_organization_permissions.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&permissions,
		getResourcePermissions,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccswrOrganizationPermissions_basic(organizationName, userName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					acceptance.TestCheckResourceAttrWithVariable(resourceName, "organization",
						"${huaweicloud_swr_organization.test.name}"),
					resource.TestCheckResourceAttr(resourceName, "users.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "users.0.user_name", userName+"_1"),
					resource.TestCheckResourceAttr(resourceName, "users.0.permission", "Read"),
					resource.TestCheckResourceAttr(resourceName, "users.1.user_name", userName+"_2"),
					resource.TestCheckResourceAttr(resourceName, "users.1.permission", "Write"),
					resource.TestCheckResourceAttr(resourceName, "users.2.user_name", userName+"_3"),
					resource.TestCheckResourceAttr(resourceName, "users.2.permission", "Manage"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccswrOrganizationPermissions_update(organizationName, userName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					acceptance.TestCheckResourceAttrWithVariable(resourceName, "organization",
						"${huaweicloud_swr_organization.test.name}"),
					resource.TestCheckResourceAttr(resourceName, "users.#", "4"),
					resource.TestCheckResourceAttr(resourceName, "users.0.user_name", userName+"_1"),
					resource.TestCheckResourceAttr(resourceName, "users.0.permission", "Write"),
					resource.TestCheckResourceAttr(resourceName, "users.1.user_name", userName+"_2"),
					resource.TestCheckResourceAttr(resourceName, "users.1.permission", "Read"),
					resource.TestCheckResourceAttr(resourceName, "users.2.user_name", userName+"_4"),
					resource.TestCheckResourceAttr(resourceName, "users.2.permission", "Manage"),
					resource.TestCheckResourceAttr(resourceName, "users.3.user_name", userName+"_5"),
					resource.TestCheckResourceAttr(resourceName, "users.3.permission", "Read"),
				),
			},
		},
	})
}

func testAccswrOrganizationPermissions_basic(organizationName, userName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_swr_organization" "test" {
  name = "%[1]s"
}

resource "huaweicloud_identity_user" "user_1" {
  name     = "%[2]s_1"
  enabled  = true
  password = "password12345!"
}

resource "huaweicloud_identity_user" "user_2" {
  name     = "%[2]s_2"
  enabled  = true
  password = "password12345!"
}

resource "huaweicloud_identity_user" "user_3" {
  name     = "%[2]s_3"
  enabled  = true
  password = "password12345!"
}

resource "huaweicloud_swr_organization_permissions" "test" {
  organization = huaweicloud_swr_organization.test.name

  users {
    user_name  = huaweicloud_identity_user.user_1.name
    user_id    = huaweicloud_identity_user.user_1.id
    permission = "Read"
  }
  users {
    user_name  = huaweicloud_identity_user.user_2.name
    user_id    = huaweicloud_identity_user.user_2.id
    permission = "Write"
  }
  users {
    user_id    = huaweicloud_identity_user.user_3.id
    permission = "Manage"
  }
}
`, organizationName, userName)
}

func testAccswrOrganizationPermissions_update(organizationName, userName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_swr_organization" "test" {
  name = "%[1]s"
}

resource "huaweicloud_identity_user" "user_1" {
  name     = "%[2]s_1"
  enabled  = true
  password = "password12345!"
}

resource "huaweicloud_identity_user" "user_2" {
  name     = "%[2]s_2"
  enabled  = true
  password = "password12345!"
}

resource "huaweicloud_identity_user" "user_4" {
  name     = "%[2]s_4"
  enabled  = true
  password = "password12345!"
}

resource "huaweicloud_identity_user" "user_5" {
  name     = "%[2]s_5"
  enabled  = true
  password = "password12345!"
}

resource "huaweicloud_swr_organization_permissions" "test" {
  organization = huaweicloud_swr_organization.test.name

  users {
    user_name  = huaweicloud_identity_user.user_1.name
    user_id    = huaweicloud_identity_user.user_1.id
    permission = "Write"
  }

  users {
    user_name  = huaweicloud_identity_user.user_2.name
    user_id    = huaweicloud_identity_user.user_2.id
    permission = "Read"
  }

  users {
    user_id    = huaweicloud_identity_user.user_4.id
    permission = "Manage"
  }

  users {
    user_id    = huaweicloud_identity_user.user_5.id
    permission = "Read"
  }
}
`, organizationName, userName)
}

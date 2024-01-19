package workspace

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/workspace/v2/groups"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getUserGroupFunc(conf *config.Config, _ *terraform.ResourceState) (interface{}, error) {
	client, err := conf.WorkspaceV2Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating Workspace v2 client: %s", err)
	}
	resp, err := groups.List(client, groups.ListOpts{})
	if err == nil && len(resp) < 1 {
		return nil, golangsdk.ErrDefault404{}
	}

	return resp, err
}

func TestAccUserGroup_basic(t *testing.T) {
	var (
		userGroups   []groups.UserGroup
		resourceName = "huaweicloud_workspace_user_group.test"
		baseConfig   = testAccUserGroup_base()
		name         = acceptance.RandomAccResourceName()
		updateName   = acceptance.RandomAccResourceName()
	)
	rc := acceptance.InitResourceCheck(
		resourceName,
		&userGroups,
		getUserGroupFunc,
	)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccUserGroup_basic(baseConfig, name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "type", "LOCAL"),
					resource.TestCheckResourceAttr(resourceName, "description", "Created by acc test"),
					resource.TestCheckResourceAttr(resourceName, "users.#", "1"),
				),
			},
			{
				Config: testAccUserGroup_update(baseConfig, updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttr(resourceName, "description", "Updated description"),
					resource.TestCheckResourceAttr(resourceName, "users.#", "0"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccUserGroup_base() string {
	name := acceptance.RandomAccResourceNameWithDash()

	return fmt.Sprintf(`
resource "huaweicloud_vpc" "test" {
  name = "%[1]s"
  cidr = "192.168.0.0/20"
}

resource "huaweicloud_vpc_subnet" "test" {
  vpc_id = huaweicloud_vpc.test.id

  name       = "%[1]s"
  cidr       = cidrsubnet(huaweicloud_vpc.test.cidr, 4, 1)
  gateway_ip = cidrhost(cidrsubnet(huaweicloud_vpc.test.cidr, 4, 1), 1)
}

resource "huaweicloud_workspace_service" "test" {
  access_mode = "INTERNET"
  vpc_id      = huaweicloud_vpc.test.id
  network_ids = [
    huaweicloud_vpc_subnet.test.id,
  ]
}

resource "huaweicloud_workspace_user" "test" {
  depends_on = [huaweicloud_workspace_service.test]

  name  = "%[1]s"
  email = "basic@example.com"

  password_never_expires = false
  disabled               = false
}
`, name)
}

func testAccUserGroup_basic(baseConfig, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_workspace_user_group" "test" {
  name        = "%[2]s"
  type        = "LOCAL"
  description = "Created by acc test"

  users {
    id = huaweicloud_workspace_user.test.id
  }
}
`, baseConfig, name)
}

func testAccUserGroup_update(baseConfig, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_workspace_user_group" "test" {
  name        = "%[2]s"
  type        = "LOCAL"
  description = "Updated description"
}
`, baseConfig, name)
}

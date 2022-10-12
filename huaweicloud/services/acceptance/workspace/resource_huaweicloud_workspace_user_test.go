package workspace

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/workspace/v2/users"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getUserFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.WorkspaceV2Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("Error creating Workspace v2 client: %s", err)
	}
	return users.Get(client, state.Primary.ID)
}

func TestAccUser_basic(t *testing.T) {
	var (
		user         users.UserDetail
		resourceName = "huaweicloud_workspace_user.test"
		rName        = acceptance.RandomAccResourceNameWithDash()
		currentTime  = time.Now().Format("2006-01-02T15:04:05Z")
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&user,
		getUserFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccUser_basic(rName, currentTime),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "email", "basic@example.com"),
					resource.TestCheckResourceAttr(resourceName, "description", "Created by acc test"),
					resource.TestCheckResourceAttrSet(resourceName, "account_expires"),
					resource.TestCheckResourceAttr(resourceName, "password_never_expires", "false"),
					resource.TestCheckResourceAttr(resourceName, "enable_change_password", "true"),
					resource.TestCheckResourceAttr(resourceName, "next_login_change_password", "true"),
					resource.TestCheckResourceAttr(resourceName, "disabled", "false"),
				),
			},
			{
				Config: testAccUser_update(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "email", "update@example.com"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttrSet(resourceName, "account_expires"),
					resource.TestCheckResourceAttr(resourceName, "password_never_expires", "true"),
					resource.TestCheckResourceAttr(resourceName, "enable_change_password", "false"),
					resource.TestCheckResourceAttr(resourceName, "next_login_change_password", "false"),
					resource.TestCheckResourceAttr(resourceName, "disabled", "true"),
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

func testAccUser_base(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc" "test" {
  name = "%[1]s"
  cidr = "192.168.128.0/18"
}

resource "huaweicloud_vpc_subnet" "test" {
  vpc_id = huaweicloud_vpc.test.id

  name       = "%[1]s"
  cidr       = cidrsubnet(huaweicloud_vpc.test.cidr, 4, 1)
  gateway_ip = cidrhost(cidrsubnet(huaweicloud_vpc.test.cidr, 4, 1), 1)
}

resource "huaweicloud_networking_secgroup" "test" {
  name = "%[1]s"
}

resource "huaweicloud_workspace_service" "test" {
  access_mode = "INTERNET"
  vpc_id      = huaweicloud_vpc.test.id
  network_ids = [
    huaweicloud_vpc_subnet.test.id,
  ]
}
`, rName)
}

func testAccUser_basic(rName, currentTime string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_workspace_user" "test" {
  depends_on = [huaweicloud_workspace_service.test]

  name        = "%[2]s"
  email       = "basic@example.com"
  description = "Created by acc test"

  account_expires        = timeadd("%[3]s", "1h")
  password_never_expires = false
  disabled               = false
}
`, testAccUser_base(rName), rName, currentTime)
}

func testAccUser_update(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_workspace_user" "test" {
  depends_on = [huaweicloud_workspace_service.test]

  name  = "%[2]s"
  email = "update@example.com"

  account_expires            = "0"
  password_never_expires     = true
  enable_change_password     = false
  next_login_change_password = false
  disabled                   = true
}
`, testAccUser_base(rName), rName)
}

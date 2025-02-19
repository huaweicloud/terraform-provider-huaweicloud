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
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
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
		user              users.UserDetail
		resourceName      = "huaweicloud_workspace_user.test"
		withAdminActive   = "huaweicloud_workspace_user.with_admin_active"
		rc                = acceptance.InitResourceCheck(resourceName, &user, getUserFunc)
		rcWithAdminActive = acceptance.InitResourceCheck(withAdminActive, &user, getUserFunc)

		rName       = acceptance.RandomAccResourceNameWithDash()
		currentTime = time.Now().Format("2006-01-02T15:04:05Z")
		password    = acceptance.RandomPassword()
	)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccUser_basic_step1(rName, password),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "active_type", "USER_ACTIVATE"),
					resource.TestCheckResourceAttr(resourceName, "email", "basic@example.com"),
					resource.TestCheckResourceAttr(resourceName, "phone", "+8612345678987"),
					resource.TestCheckResourceAttr(resourceName, "description", "Created by acc test"),
					resource.TestCheckResourceAttr(resourceName, "account_expires", "0"),
					resource.TestCheckResourceAttr(resourceName, "password_never_expires", "false"),
					resource.TestCheckResourceAttr(resourceName, "enable_change_password", "true"),
					resource.TestCheckResourceAttr(resourceName, "next_login_change_password", "true"),
					resource.TestCheckResourceAttr(resourceName, "disabled", "false"),
					rcWithAdminActive.CheckResourceExists(),
					resource.TestCheckResourceAttr(withAdminActive, "active_type", "ADMIN_ACTIVATE"),
					resource.TestCheckResourceAttr(withAdminActive, "disabled", "true"),
					resource.TestCheckResourceAttr(withAdminActive, "email", ""),
					resource.TestCheckResourceAttr(withAdminActive, "phone", "+8612345678978"),
				),
			},
			{
				Config: testAccUser_basic_step2(rName, currentTime, password),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "active_type", "ADMIN_ACTIVATE"),
					resource.TestCheckResourceAttr(resourceName, "email", ""),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttrSet(resourceName, "account_expires"),
					resource.TestCheckResourceAttr(resourceName, "password_never_expires", "true"),
					resource.TestCheckResourceAttr(resourceName, "enable_change_password", "false"),
					resource.TestCheckResourceAttr(resourceName, "next_login_change_password", "false"),
					resource.TestCheckResourceAttr(resourceName, "disabled", "true"),
					rcWithAdminActive.CheckResourceExists(),
					resource.TestCheckResourceAttr(withAdminActive, "active_type", "USER_ACTIVATE"),
					resource.TestCheckResourceAttr(withAdminActive, "disabled", "false"),
					resource.TestCheckResourceAttr(withAdminActive, "phone", ""),
					resource.TestCheckResourceAttr(withAdminActive, "email", "update@example.com"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}

func testAccUser_base(rName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_workspace_service" "test" {
  access_mode = "INTERNET"
  vpc_id      = huaweicloud_vpc.test.id
  network_ids = [
    huaweicloud_vpc_subnet.test.id,
  ]
}
`, common.TestBaseNetwork(rName))
}

func testAccUser_basic_step1(rName, password string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_workspace_user" "test" {
  depends_on = [huaweicloud_workspace_service.test]

  name        = "%[2]s"
  email       = "basic@example.com"
  phone       = "+8612345678987"
  description = "Created by acc test"

  password_never_expires = false
  disabled               = false
}

resource "huaweicloud_workspace_user" "with_admin_active" {
  depends_on = [huaweicloud_workspace_service.test]

  name        = "%[2]s_with_admin"
  active_type = "ADMIN_ACTIVATE"
  password    = "%[3]s"
  phone       = "+8612345678978"

  password_never_expires = false
  disabled               = true
}
`, testAccUser_base(rName), rName, password)
}

func testAccUser_basic_step2(rName, currentTime, password string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_workspace_user" "test" {
  depends_on = [huaweicloud_workspace_service.test]

  name        = "%[2]s"
  active_type = "ADMIN_ACTIVATE"
  phone       = "+8612345678987"
  password    = "%[4]s"

  account_expires            = timeadd("%[3]s", "1h")
  password_never_expires     = true
  enable_change_password     = false
  next_login_change_password = false
  disabled                   = true
}

resource "huaweicloud_workspace_user" "with_admin_active" {
  depends_on = [huaweicloud_workspace_service.test]

  name        = "%[2]s_with_admin"
  active_type = "USER_ACTIVATE"
  email       = "update@example.com"

  password_never_expires = false
  disabled               = false
}
`, testAccUser_base(rName), rName, currentTime, password)
}

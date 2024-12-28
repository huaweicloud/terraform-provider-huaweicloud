package workspace

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccResourceAppGroupAuthorization_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckWorkspaceAppServerGroup(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccAppGroupAuthorization_basic(name),
			},
		},
	})
}

func testAccAppGroupAuthorization_base(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_workspace_service" "test" {}

resource "huaweicloud_workspace_app_server_group" "test" {
  name             = "%[1]s"
  os_type          = "Windows"
  flavor_id        = "%[2]s"
  vpc_id           = data.huaweicloud_workspace_service.test.vpc_id
  subnet_id        = data.huaweicloud_workspace_service.test.network_ids[0]
  system_disk_type = "SAS"
  system_disk_size = 80
  is_vdi           = true
  app_type         = "SESSION_DESKTOP_APP"
  image_id         = "%[3]s"
  image_type       = "gold"
  image_product_id = "%[4]s"
}

resource "huaweicloud_workspace_app_group" "test" {
  server_group_id = huaweicloud_workspace_app_server_group.test.id
  name            = "%[1]s"
  type            = "SESSION_DESKTOP_APP"
  description     = "Created APP group by script"
}

resource "huaweicloud_workspace_user" "test" {
  name  = "%[1]s"
  email = "tf@example.com"
}

resource "huaweicloud_workspace_user_group" "test" {
  count = 2

  name  = "%[1]s${count.index}"
  type  = "LOCAL"
}
`, name, acceptance.HW_WORKSPACE_APP_SERVER_GROUP_FLAVOR_ID,
		acceptance.HW_WORKSPACE_APP_SERVER_GROUP_IMAGE_ID,
		acceptance.HW_WORKSPACE_APP_SERVER_GROUP_IMAGE_PRODUCT_ID)
}

func testAccAppGroupAuthorization_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_workspace_app_group_authorization" "test" {
  app_group_id = huaweicloud_workspace_app_group.test.id

  accounts {
    account = huaweicloud_workspace_user.test.name
    type    = "USER"
  }

  dynamic "accounts" {
    for_each = huaweicloud_workspace_user_group.test[*]

    content {
      id      = accounts.value.id
      account = accounts.value.name
      type    = "USER_GROUP"
    }
  }
}
`, testAccAppGroupAuthorization_base(name))
}

func TestAccResourceAppGroupAuthorization_expectErr(t *testing.T) {
	name := acceptance.RandomAccResourceName()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckWorkspaceAppServerGroup(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config:      testAccAppGroupAuthorization_expectErr(name),
				ExpectError: regexp.MustCompile(`unable to authorize for some accounts: not_exist_user_group_tf | USER`),
			},
		},
	})
}

func testAccAppGroupAuthorization_expectErr(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_workspace_user" "test" {
  count = 2
  name  = "%[2]s${count.index}"
  email = "tf@example.com"
}

resource "huaweicloud_workspace_app_group_authorization" "test" {
  app_group_id = huaweicloud_workspace_app_group.test.id

  dynamic "accounts" {
    for_each = huaweicloud_workspace_user.test[*]
  
    content {
      account = accounts.value.name
      type    = "USER"
    }
  }

  accounts {
    account = "not_exist_user_group_tf"
    type    = "USER"
  }
}
`, testResourceWorkspaceAppGroup_basic_step1(testResourceWorkspaceAppGroup_base(name), name), name)
}

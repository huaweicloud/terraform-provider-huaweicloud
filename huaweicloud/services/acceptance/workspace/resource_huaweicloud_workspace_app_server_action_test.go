package workspace

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccAppServerAction_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceName()

		changeImageRName = "huaweicloud_workspace_app_server_action.changeImage"
		reinstallRName   = "huaweicloud_workspace_app_server_action.reinstall"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckWorkspaceAppServerGroupId(t)
			acceptance.TestAccPreCheckWorkspaceAppServerImageId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		// This resource is a one-time action resource and there is no logic in the delete method.
		// lintignore:AT001
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccAppServerAction_changeImage(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(changeImageRName, "type", "change-image"),
				),
			},
			{
				Config: testAccAppServerAction_reinstall(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(reinstallRName, "type", "reinstall"),
				),
			},
		},
	})
}

func testAccAppServerAction_base(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_workspace_app_server_groups" "test" {
  server_group_id = "%[1]s"
}

data "huaweicloud_vpc_subnets" "test" {
  id = try(data.huaweicloud_workspace_app_server_groups.test.server_groups[0].subnet_id, null)
}

resource "huaweicloud_workspace_app_server" "test" {
  name                = "%[2]s" 
  server_group_id     = try(data.huaweicloud_workspace_app_server_groups.test.server_groups[0].id, null)
  type                = "createApps"
  flavor_id           = try(data.huaweicloud_workspace_app_server_groups.test.server_groups[0].product_id, null)
  vpc_id              = try(data.huaweicloud_vpc_subnets.test.subnets[0].vpc_id, null)
  subnet_id           = try(data.huaweicloud_workspace_app_server_groups.test.server_groups[0].subnet_id, null)
  update_access_agent = false
  description         = "Created by terraform script"
  maintain_status     = true

  root_volume {
    type = try(data.huaweicloud_workspace_app_server_groups.test.server_groups[0].system_disk_type, null)
    size = try(data.huaweicloud_workspace_app_server_groups.test.server_groups[0].system_disk_size, null)
  }

  lifecycle {
    ignore_changes = [
      maintain_status
    ]
  }
}
`, acceptance.HW_WORKSPACE_APP_SERVER_GROUP_ID, name)
}

func testAccAppServerAction_changeImage(name string) string {
	return fmt.Sprintf(`
%[1]s

variable "image_product_id" {
  default = "%[2]s"
}

resource "huaweicloud_workspace_app_server_action" "changeImage" {
  type      = "change-image"
  server_id = huaweicloud_workspace_app_server.test.id
  content   = jsonencode({
    image_id            = "%[3]s"
    os_type             = "Windows"
    image_type          = var.image_product_id != "" ? "gold" : null
	image_product_id    = var.image_product_id != "" ? var.image_product_id : null
    update_access_agent = true
  })

  max_retries = 3
}
`, testAccAppServerAction_base(name),
		acceptance.HW_WORKSPACE_APP_SERVER_GROUP_IMAGE_PRODUCT_ID,
		acceptance.HW_WORKSPACE_APP_SERVER_GROUP_IMAGE_ID)
}

func testAccAppServerAction_reinstall(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_workspace_app_server_action" "reinstall" {
  type      = "reinstall"
  server_id = huaweicloud_workspace_app_server.test.id
  content   = jsonencode({
    update_access_agent = false
  })

  max_retries = 3
}
`, testAccAppServerAction_base(name))
}

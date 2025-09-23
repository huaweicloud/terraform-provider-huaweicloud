package workspace

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccAppHdaBatchUpgrade_basic(t *testing.T) {
	var (
		resourceName = "huaweicloud_workspace_app_hda_batch_upgrade.test"
		name         = acceptance.RandomAccResourceName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckWorkspaceAppServerGroupId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		// This resource is a one-time action resource and there is no logic in the delete method.
		// lintignore:AT001
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				// The first step is to upgrade the first server.
				Config: testAccAppHdaBatchUpgrade_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "server_sids.#", "1"),
				),
			},
			{
				// The second step is to upgrade the remaining two servers.
				Config: testAccAppHdaBatchUpgrade_basic_step2(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "server_sids.#", "2"),
				),
			},
		},
	})
}

func testAccAppHdaBatchUpgrade_base(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_workspace_app_server_groups" "test" {
  server_group_id = "%[1]s"
}

data "huaweicloud_vpc_subnets" "test" {
  id = try(data.huaweicloud_workspace_app_server_groups.test.server_groups[0].subnet_id, null)
}

resource "huaweicloud_workspace_app_server" "test" {
  count = 3

  name                = format("%[2]s_%%d", count.index)
  server_group_id     = try(data.huaweicloud_workspace_app_server_groups.test.server_groups[0].id, null)
  type                = "createApps"
  flavor_id           = try(data.huaweicloud_workspace_app_server_groups.test.server_groups[0].product_id, null)
  vpc_id              = try(data.huaweicloud_vpc_subnets.test.subnets[0].vpc_id, null)
  subnet_id           = try(data.huaweicloud_workspace_app_server_groups.test.server_groups[0].subnet_id, null)
  update_access_agent = false
  description         = "Created by terraform script for HDA upgrade test"
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

func testAccAppHdaBatchUpgrade_basic_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_workspace_app_hda_batch_upgrade" "test" {
  depends_on = [
    huaweicloud_workspace_app_server.test
  ]

  server_ids       = slice(huaweicloud_workspace_app_server.test[*].id, 0, 1)
  enable_force_new = true

  provisioner "local-exec" {
    command = "sleep 300"
  }
}
`, testAccAppHdaBatchUpgrade_base(name))
}

func testAccAppHdaBatchUpgrade_basic_step2(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_workspace_app_hda_batch_upgrade" "test" {
  depends_on = [
    huaweicloud_workspace_app_server.test
  ]

  server_ids       = slice(huaweicloud_workspace_app_server.test[*].id, 1, 3)
  enable_force_new = true

  provisioner "local-exec" {
    command = "sleep 300"
  }
}
`, testAccAppHdaBatchUpgrade_base(name))
}

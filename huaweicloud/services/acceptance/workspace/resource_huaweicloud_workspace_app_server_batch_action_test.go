package workspace

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccAppServerBatchAction_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceName()

		changeImageRName = "huaweicloud_workspace_app_server_batch_action.changeImage"
		reinstallRName   = "huaweicloud_workspace_app_server_batch_action.reinstall"
		maintainRName    = "huaweicloud_workspace_app_server_batch_action.maintain"
		rebootRName      = "huaweicloud_workspace_app_server_batch_action.reboot"
		stopRName        = "huaweicloud_workspace_app_server_batch_action.stop"
		startRName       = "huaweicloud_workspace_app_server_batch_action.start"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckWorkspaceAppServerGroupId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		// This resource is a one-time batch action resource and there is no logic in the delete method.
		// lintignore:AT001
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccAppServerBatchAction_batchChangeImage(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(changeImageRName, "type", "batch-change-image"),
				),
			},
			{
				Config: testAccAppServerBatchAction_batchReinstall(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(reinstallRName, "type", "batch-reinstall"),
				),
			},
			{
				Config: testAccAppServerBatchAction_maintain(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(maintainRName, "type", "batch-maint"),
				),
			},
			{
				Config: testAccAppServerBatchAction_batchReboot(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rebootRName, "type", "batch-reboot"),
				),
			},
			{
				Config: testAccAppServerBatchAction_batchStop(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(stopRName, "type", "batch-stop"),
				),
			},
			{
				Config: testAccAppServerBatchAction_batchStart(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(startRName, "type", "batch-start"),
				),
			},
		},
	})
}

func testAccAppServerBatchAction_base(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_workspace_app_server_groups" "test" {
  server_group_id = "%[1]s"
}

data "huaweicloud_vpc_subnets" "test" {
  id = try(data.huaweicloud_workspace_app_server_groups.test.server_groups[0].subnet_id, null)
}

resource "huaweicloud_workspace_app_server" "test" {
  count = 2

  name                = "%[2]s${count.index}" 
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

func testAccAppServerBatchAction_batchChangeImage(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_images_images" "test" {
  visibility = "market"
  os         = "Windows"
}

locals {
  available_images = try([
    for o in data.huaweicloud_images_images.test.images : o if 
      strcontains(lower(o.name), "appstream") && o.id != try(data.huaweicloud_workspace_app_server_groups.test.server_groups[0].image_id, "")
    ], [])
}

resource "huaweicloud_workspace_app_server_batch_action" "changeImage" {
  type    = "batch-change-image"
  content = jsonencode({
    server_ids          = huaweicloud_workspace_app_server.test[*].id
    image_id            = try(local.available_images[0].id, "NOT_FOUND")
    image_type          = "gold"
    os_type             = "Windows"
    update_access_agent = true
  })

  max_retries = 3
}
`, testAccAppServerBatchAction_base(name))
}

func testAccAppServerBatchAction_batchReinstall(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_workspace_app_server_batch_action" "reinstall" {
  type    = "batch-reinstall"
  content = jsonencode({
    server_ids          = huaweicloud_workspace_app_server.test[*].id
    update_access_agent = false
  })

  max_retries = 3
}
`, testAccAppServerBatchAction_base(name))
}

func testAccAppServerBatchAction_maintain(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_workspace_app_server_batch_action" "maintain" {
  type    = "batch-maint"
  content = jsonencode({
    items           = huaweicloud_workspace_app_server.test[*].id
    maintain_status = true
  })

  max_retries = 3
}
`, testAccAppServerBatchAction_base(name))
}

func testAccAppServerBatchAction_batchReboot(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_workspace_app_server_batch_action" "reboot" {
  type    = "batch-reboot"
  content = jsonencode({
    items = huaweicloud_workspace_app_server.test[*].id
    type  = "SOFT"
  })

  max_retries = 3
}
`, testAccAppServerBatchAction_base(name))
}

func testAccAppServerBatchAction_batchStop(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_workspace_app_server_batch_action" "stop" {
  type    = "batch-stop"
  content = jsonencode({
    items = huaweicloud_workspace_app_server.test[*].id
    type  = "SOFT"
  })

  max_retries = 3
}
`, testAccAppServerBatchAction_base(name))
}

func testAccAppServerBatchAction_batchStart(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_workspace_app_server_batch_action" "start" {
  type    = "batch-start"
  content = jsonencode({
    items = huaweicloud_workspace_app_server.test[*].id
  })

  max_retries = 3
}
`, testAccAppServerBatchAction_base(name))
}

// The Workspace serivice must connect to the Active Directory server.
func TestAccAppServerBatchAction_batchRejoinDomain(t *testing.T) {
	var (
		resourceName = "huaweicloud_workspace_app_server_batch_action.test"
		name         = acceptance.RandomAccResourceName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckWorkspaceAppServerGroupId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		// This resource is a one-time batch action resource and there is no logic in the delete method.
		// lintignore:AT001
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccAppServerBatchAction_batchRejoinDomain(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "type", "batch-rejoin-domain"),
				),
			},
		},
	})
}

func testAccAppServerBatchAction_batchRejoinDomain(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_workspace_app_server_batch_action" "test" {
  type    = "batch-rejoin-domain"
  content = jsonencode({
    items = [huaweicloud_workspace_app_server.test[0].id]
  })

  max_retries = 3
}
`, testAccAppServerBatchAction_base(name))
}

func TestAccAppServerBatchAction_batchUpdateTsvi(t *testing.T) {
	var (
		resourceName = "huaweicloud_workspace_app_server_batch_action.test"
		name         = acceptance.RandomAccResourceName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckWorkspaceAppServerGroupId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		// This resource is a one-time batch action resource and there is no logic in the delete method.
		// lintignore:AT001
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccAppServerBatchAction_batchUpdateTsvi(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "type", "batch-update-tsvi"),
				),
			},
		},
	})
}

func testAccAppServerBatchAction_batchUpdateTsvi(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_workspace_app_server_batch_action" "test" {
  type    = "batch-update-tsvi"
  content = jsonencode({
    items = [
      {
        id     = huaweicloud_workspace_app_server.test.id
        enable = true
      }
    ]
  })

  max_retries = 3
}
`, testAccAppServerBatchAction_base(name))
}

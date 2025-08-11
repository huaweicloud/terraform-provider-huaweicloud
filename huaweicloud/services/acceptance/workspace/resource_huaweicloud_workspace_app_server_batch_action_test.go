package workspace

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccAppServerBatchAction_batchChangeImage(t *testing.T) {
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
				Config: testAccAppServerBatchAction_batchChangeImage(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "type", "batch-change-image"),
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

resource "huaweicloud_workspace_app_server_batch_action" "test" {
  type    = "batch-change-image"
  content = jsonencode({
    server_ids          = [huaweicloud_workspace_app_server.test.id]
    image_id            = try(local.available_images[0].id, "NOT_FOUND")
    image_type          = "gold"
    os_type             = "Windows"
    update_access_agent = true
  })

  max_retries = 3

  provisioner "local-exec" {
    command = "sleep 600"
  }
}
`, testAccAppServerBatchAction_base(name))
}

func TestAccAppServerBatchAction_batchReinstall(t *testing.T) {
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
				Config: testAccAppServerBatchAction_batchReinstall(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "type", "batch-reinstall"),
				),
			},
		},
	})
}

func testAccAppServerBatchAction_batchReinstall(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_workspace_app_server_batch_action" "test" {
  type    = "batch-reinstall"
  content = jsonencode({
    server_ids          = [huaweicloud_workspace_app_server.test.id]
    update_access_agent = false
  })

  max_retries = 3

  provisioner "local-exec" {
    command = "sleep 1200"
  }
}
`, testAccAppServerBatchAction_base(name))
}

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
    items = [huaweicloud_workspace_app_server.test.id]
  })

  max_retries = 3

  provisioner "local-exec" {
    command = "sleep 300"
  }
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

  provisioner "local-exec" {
    command = "sleep 600"
  }
}
`, testAccAppServerBatchAction_base(name))
}

package workspace

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataDesktopPoolAssociatedDesktops_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceNameWithDash()

		dataSourceName = "data.huaweicloud_workspace_desktop_pool_associated_desktops.all"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckWorkspaceDesktopPoolImageId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataDesktopPoolAssociatedDesktops_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dataSourceName, "desktops.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(dataSourceName, "desktops.0.desktop_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "desktops.0.computer_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "desktops.0.os_host_name"),
					resource.TestMatchResourceAttr(dataSourceName, "desktops.0.ip_addresses.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(dataSourceName, "desktops.0.ipv4"),
					resource.TestCheckResourceAttrSet(dataSourceName, "desktops.0.desktop_type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "desktops.0.status"),
					resource.TestCheckResourceAttrSet(dataSourceName, "desktops.0.in_maintenance_mode"),
					resource.TestCheckResourceAttrSet(dataSourceName, "desktops.0.created"),
					resource.TestCheckResourceAttrSet(dataSourceName, "desktops.0.login_status"),
					resource.TestCheckResourceAttrSet(dataSourceName, "desktops.0.product_id"),
					resource.TestMatchResourceAttr(dataSourceName, "desktops.0.root_volume.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(dataSourceName, "desktops.0.root_volume.0.type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "desktops.0.root_volume.0.size"),
					resource.TestCheckResourceAttrSet(dataSourceName, "desktops.0.root_volume.0.device"),
					resource.TestCheckResourceAttrSet(dataSourceName, "desktops.0.root_volume.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "desktops.0.root_volume.0.volume_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "desktops.0.root_volume.0.bill_resource_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "desktops.0.root_volume.0.create_time"),
					resource.TestCheckResourceAttrSet(dataSourceName, "desktops.0.root_volume.0.display_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "desktops.0.root_volume.0.resource_spec_code"),
					resource.TestMatchResourceAttr(dataSourceName, "desktops.0.data_volumes.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(dataSourceName, "desktops.0.data_volumes.0.type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "desktops.0.data_volumes.0.size"),
					resource.TestCheckResourceAttrSet(dataSourceName, "desktops.0.data_volumes.0.device"),
					resource.TestCheckResourceAttrSet(dataSourceName, "desktops.0.data_volumes.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "desktops.0.data_volumes.0.volume_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "desktops.0.data_volumes.0.bill_resource_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "desktops.0.data_volumes.0.create_time"),
					resource.TestCheckResourceAttrSet(dataSourceName, "desktops.0.data_volumes.0.display_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "desktops.0.data_volumes.0.resource_spec_code"),
					resource.TestCheckResourceAttrSet(dataSourceName, "desktops.0.availability_zone"),
					resource.TestCheckResourceAttrSet(dataSourceName, "desktops.0.site_type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "desktops.0.site_name"),
					resource.TestMatchResourceAttr(dataSourceName, "desktops.0.product.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(dataSourceName, "desktops.0.product.0.product_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "desktops.0.product.0.flavor_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "desktops.0.product.0.type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "desktops.0.product.0.cpu"),
					resource.TestCheckResourceAttrSet(dataSourceName, "desktops.0.product.0.memory"),
					resource.TestCheckResourceAttrSet(dataSourceName, "desktops.0.product.0.descriptions"),
					resource.TestCheckResourceAttrSet(dataSourceName, "desktops.0.product.0.charge_mode"),
					resource.TestCheckResourceAttrSet(dataSourceName, "desktops.0.product.0.architecture"),
					resource.TestCheckResourceAttrSet(dataSourceName, "desktops.0.product.0.is_gpu"),
					resource.TestCheckResourceAttrSet(dataSourceName, "desktops.0.product.0.package_type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "desktops.0.product.0.system_disk_type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "desktops.0.product.0.system_disk_size"),
					resource.TestCheckResourceAttrSet(dataSourceName, "desktops.0.product.0.contain_data_disk"),
					resource.TestCheckResourceAttrSet(dataSourceName, "desktops.0.product.0.resource_type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "desktops.0.product.0.cloud_service_type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "desktops.0.product.0.volume_product_type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "desktops.0.product.0.status"),
					resource.TestCheckResourceAttrSet(dataSourceName, "desktops.0.os_version"),
					resource.TestCheckResourceAttrSet(dataSourceName, "desktops.0.sid"),
					resource.TestMatchResourceAttr(dataSourceName, "desktops.0.tags.%", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttr(dataSourceName, "desktops.0.tags.foo", "bar"),
					resource.TestCheckResourceAttr(dataSourceName, "desktops.0.tags.owner", "terraform"),
					resource.TestCheckResourceAttrSet(dataSourceName, "desktops.0.is_support_internet"),
					resource.TestCheckResourceAttrSet(dataSourceName, "desktops.0.is_attaching_eip"),
					resource.TestCheckResourceAttrSet(dataSourceName, "desktops.0.attach_state"),
					resource.TestCheckResourceAttrSet(dataSourceName, "desktops.0.enterprise_project_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "desktops.0.subnet_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "desktops.0.bill_resource_id"),
				),
			},
		},
	})
}

func testAccDataDesktopPoolAssociatedDesktops_basic(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_workspace_service" "test" {}

data "huaweicloud_workspace_flavors" "test" {
  availability_zone = try(data.huaweicloud_availability_zones.test.names[0], null)
  os_type           = "Windows"
}

resource "huaweicloud_workspace_user" "test" {
  name  = "%[1]s"
  email = "test@example.com"
}

resource "huaweicloud_workspace_desktop_pool" "test" {
  name                          = "%[1]s"
  type                          = "DYNAMIC"
  size                          = 2
  product_id                    = try(data.huaweicloud_workspace_flavors.test.flavors[0].id, "")
  image_type                    = "gold"
  image_id                      = "%[2]s"
  subnet_ids                    = try(slice(data.huaweicloud_workspace_service.test.network_ids, 0, 1), [])
  vpc_id                        = data.huaweicloud_workspace_service.test.vpc_id
  availability_zone             = data.huaweicloud_availability_zones.test.names[0]
  in_maintenance_mode           = true
  enable_autoscale              = true
  disconnected_retention_period = 10

  security_groups {
    id = data.huaweicloud_workspace_service.test.desktop_security_group[0].id
  }

  root_volume {
    type = "SAS"
    size = 80
  }

  data_volumes {
    type = "SAS"
    size = 50
  }

  authorized_objects {
    object_id   = huaweicloud_workspace_user.test.id
    object_type = "USER"
    object_name = huaweicloud_workspace_user.test.name
    user_group  = "administrators"
  }

  autoscale_policy {
    autoscale_type    = "AUTO_CREATED"
    min_idle          = 1
    max_auto_created  = 2
    once_auto_created = 1
  }

  tags = {
    foo   = "bar"
    owner = "terraform"
  }

  lifecycle {
    ignore_changes = [
      size,
    ]
  }
}

data "huaweicloud_workspace_desktop_pool_associated_desktops" "all" {
  pool_id = huaweicloud_workspace_desktop_pool.test.id
}
`, name, acceptance.HW_WORKSPACE_DESKTOP_POOL_IMAGE_ID)
}

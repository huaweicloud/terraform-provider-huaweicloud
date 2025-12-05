package workspace

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataDesktopPools_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceNameWithDash()

		all = "data.huaweicloud_workspace_desktop_pools.all"
		dc  = acceptance.InitDataSourceCheck(all)

		filterByName   = "data.huaweicloud_workspace_desktop_pools.filter_by_name"
		dcFilterByName = acceptance.InitDataSourceCheck(filterByName)

		filterByType   = "data.huaweicloud_workspace_desktop_pools.filter_by_type"
		dcFilterByType = acceptance.InitDataSourceCheck(filterByType)

		filterByEnterpriseProjectId   = "data.huaweicloud_workspace_desktop_pools.filter_by_enterprise_project_id"
		dcFilterByEnterpriseProjectId = acceptance.InitDataSourceCheck(filterByEnterpriseProjectId)

		filterByMaintenanceMode   = "data.huaweicloud_workspace_desktop_pools.filter_by_in_maintenance_mode"
		dcFilterByMaintenanceMode = acceptance.InitDataSourceCheck(filterByMaintenanceMode)
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckWorkspaceDesktopPoolImageId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataDesktopPools_basic(name),
				Check: resource.ComposeTestCheckFunc(
					// Without any filter parameter.
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "desktop_pools.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(all, "desktop_pools.0.id"),
					resource.TestCheckResourceAttrSet(all, "desktop_pools.0.name"),
					resource.TestMatchResourceAttr(all, "desktop_pools.0.autoscale_policy.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(all, "desktop_pools.0.autoscale_policy.0.autoscale_type"),
					resource.TestCheckResourceAttrSet(all, "desktop_pools.0.autoscale_policy.0.max_auto_created"),
					resource.TestCheckResourceAttrSet(all, "desktop_pools.0.autoscale_policy.0.min_idle"),
					resource.TestCheckResourceAttrSet(all, "desktop_pools.0.autoscale_policy.0.once_auto_created"),
					resource.TestCheckResourceAttrSet(all, "desktop_pools.0.availability_zone"),
					resource.TestCheckResourceAttrSet(all, "desktop_pools.0.charging_mode"),
					resource.TestCheckResourceAttrSet(all, "desktop_pools.0.created_time"),
					resource.TestMatchResourceAttr(all, "desktop_pools.0.data_volumes.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(all, "desktop_pools.0.desktop_count"),
					resource.TestCheckResourceAttrSet(all, "desktop_pools.0.desktop_used"),
					resource.TestCheckResourceAttrSet(all, "desktop_pools.0.disconnected_retention_period"),
					resource.TestCheckResourceAttrSet(all, "desktop_pools.0.enable_autoscale"),
					resource.TestCheckResourceAttrSet(all, "desktop_pools.0.enterprise_project_id"),
					resource.TestCheckResourceAttrSet(all, "desktop_pools.0.image_id"),
					resource.TestCheckResourceAttrSet(all, "desktop_pools.0.image_name"),
					resource.TestCheckResourceAttrSet(all, "desktop_pools.0.image_os_platform"),
					resource.TestCheckResourceAttrSet(all, "desktop_pools.0.image_os_type"),
					resource.TestCheckResourceAttrSet(all, "desktop_pools.0.image_os_version"),
					resource.TestCheckResourceAttrSet(all, "desktop_pools.0.in_maintenance_mode"),
					resource.TestMatchResourceAttr(all, "desktop_pools.0.product.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestMatchResourceAttr(all, "desktop_pools.0.root_volume.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestMatchResourceAttr(all, "desktop_pools.0.security_groups.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(all, "desktop_pools.0.status"),
					resource.TestCheckResourceAttrSet(all, "desktop_pools.0.subnet_id"),
					resource.TestCheckResourceAttrSet(all, "desktop_pools.0.type"),
					// Filter by 'name' parameter.
					dcFilterByName.CheckResourceExists(),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
					// Filter by 'type' parameter.
					dcFilterByType.CheckResourceExists(),
					resource.TestCheckOutput("is_type_filter_useful", "true"),
					// Filter by 'enterprise_project_id' parameter.
					dcFilterByEnterpriseProjectId.CheckResourceExists(),
					resource.TestCheckOutput("is_enterprise_project_id_filter_useful", "true"),
					// Filter by 'in_maintenance_mode' parameter.
					dcFilterByMaintenanceMode.CheckResourceExists(),
					resource.TestCheckOutput("is_in_maintenance_mode_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataDesktopPools_base(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_workspace_service" "test" {}

data "huaweicloud_workspace_flavors" "test" {
  os_type = "Windows"
}

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_workspace_desktop_pool" "test" {
  name                          = "%[1]s"
  type                          = "DYNAMIC"
  size                          = 2
  product_id                    = try(data.huaweicloud_workspace_flavors.test.flavors[0].id, "")
  image_type                    = "gold"
  image_id                      = "%[2]s"
  subnet_ids                    = data.huaweicloud_workspace_service.test.network_ids
  vpc_id                        = data.huaweicloud_workspace_service.test.vpc_id
  availability_zone             = data.huaweicloud_availability_zones.test.names[0]
  disconnected_retention_period = 10
  enable_autoscale              = true
  in_maintenance_mode           = true

  security_groups {
    id = data.huaweicloud_workspace_service.test.desktop_security_group.0.id
  }
  security_groups {
    id = data.huaweicloud_workspace_service.test.infrastructure_security_group.0.id
  }

  root_volume {
    type = "SAS"
    size = 80
  }

  data_volumes {
    type = "SAS"
    size = 20
  }

  autoscale_policy {
    autoscale_type    = "AUTO_CREATED"
    min_idle          = 1
    max_auto_created  = 2
    once_auto_created = 1
  }
}
`, name, acceptance.HW_WORKSPACE_DESKTOP_POOL_IMAGE_ID)
}

func testAccDataDesktopPools_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

locals {
  name                  = huaweicloud_workspace_desktop_pool.test.name
  type                  = huaweicloud_workspace_desktop_pool.test.type
  enterprise_project_id = huaweicloud_workspace_desktop_pool.test.enterprise_project_id
  in_maintenance_mode 	= huaweicloud_workspace_desktop_pool.test.in_maintenance_mode
}

# Without any filter parameter.
data "huaweicloud_workspace_desktop_pools" "all" {
  depends_on = [
    huaweicloud_workspace_desktop_pool.test
  ]
}

# Filter by 'name' parameter.
data "huaweicloud_workspace_desktop_pools" "filter_by_name" {
  name = local.name

  depends_on = [
    huaweicloud_workspace_desktop_pool.test
  ]
}

locals {
  name_filter_result = [
    for v in data.huaweicloud_workspace_desktop_pools.filter_by_name.desktop_pools[*].name : v == local.name
  ]
}

output "is_name_filter_useful" {
  value = length(local.name_filter_result) > 0 && alltrue(local.name_filter_result)
}

# Filter by 'type' parameter.
data "huaweicloud_workspace_desktop_pools" "filter_by_type" {
  type = local.type

  depends_on = [
    huaweicloud_workspace_desktop_pool.test
  ]
}

locals {
  type_filter_result = [
    for v in data.huaweicloud_workspace_desktop_pools.filter_by_type.desktop_pools[*].type : v == local.type
  ]
}

output "is_type_filter_useful" {
  value = length(local.type_filter_result) > 0 && alltrue(local.type_filter_result)
}

# Filter by 'enterprise_project_id' parameter.
data "huaweicloud_workspace_desktop_pools" "filter_by_enterprise_project_id" {
  enterprise_project_id = local.enterprise_project_id

  depends_on = [
    huaweicloud_workspace_desktop_pool.test
  ]
}

locals {
  enterprise_project_id_filter_result = [
    for v in data.huaweicloud_workspace_desktop_pools.filter_by_enterprise_project_id.desktop_pools[*].enterprise_project_id : 
      v == local.enterprise_project_id
  ]
}

output "is_enterprise_project_id_filter_useful" {
  value = length(local.enterprise_project_id_filter_result) > 0 && alltrue(local.enterprise_project_id_filter_result)
}

# Filter by 'in_maintenance_mode' parameter.
data "huaweicloud_workspace_desktop_pools" "filter_by_in_maintenance_mode" {
  in_maintenance_mode = local.in_maintenance_mode

  depends_on = [
    huaweicloud_workspace_desktop_pool.test
  ]
}

locals {
  in_maintenance_mode_filter_result = [
    for v in data.huaweicloud_workspace_desktop_pools.filter_by_in_maintenance_mode.desktop_pools[*].in_maintenance_mode : 
      v == local.in_maintenance_mode
  ]
}

output "is_in_maintenance_mode_filter_useful" {
  value = length(local.in_maintenance_mode_filter_result) > 0 && alltrue(local.in_maintenance_mode_filter_result)
}
`, testAccDataDesktopPools_base(name))
}

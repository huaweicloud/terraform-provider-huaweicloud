package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceHssImageRegistries_basic(t *testing.T) {
	dataSource := "data.huaweicloud_hss_image_registries.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceHssImageRegistries_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.registry_name"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.registry_type"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.api_version"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.protocol"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.registry_addr"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.registry_username"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.namespace"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.connect_cluster_id"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.get_scan_image_channel"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.fail_reason"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.images_num"),

					resource.TestCheckOutput("is_registry_name_filter_useful", "true"),
					resource.TestCheckOutput("is_registry_type_filter_useful", "true"),
					resource.TestCheckOutput("is_registry_addr_filter_useful", "true"),
					resource.TestCheckOutput("is_status_filter_useful", "true"),
				),
			},
		},
	})
}

const testDataSourceDataSourceHssImageRegistries_basic = `
data "huaweicloud_hss_image_registries" "test" {
  enterprise_project_id = "all_granted_eps"
}

# Filter by registry_name
locals {
  registry_name = data.huaweicloud_hss_image_registries.test.data_list[0].registry_name
}

data "huaweicloud_hss_image_registries" "test_registry_name" {
  enterprise_project_id = "all_granted_eps"
  registry_name         = local.registry_name
}

output "is_registry_name_filter_useful" {
  value = length(data.huaweicloud_hss_image_registries.test_registry_name.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_image_registries.test_registry_name.data_list[*].registry_name : v == local.registry_name]
  )
}

# Filter by registry_type
locals {
  registry_type = data.huaweicloud_hss_image_registries.test.data_list[0].registry_type
}

data "huaweicloud_hss_image_registries" "test_registry_type" {
  enterprise_project_id = "all_granted_eps"
  registry_type         = local.registry_type
}

output "is_registry_type_filter_useful" {
  value = length(data.huaweicloud_hss_image_registries.test_registry_type.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_image_registries.test_registry_type.data_list[*].registry_type : v == local.registry_type]
  )
}

# Filter by registry_addr
locals {
  registry_addr = data.huaweicloud_hss_image_registries.test.data_list[0].registry_addr
}

data "huaweicloud_hss_image_registries" "test_registry_addr" {
  enterprise_project_id = "all_granted_eps"
  registry_addr         = local.registry_addr
}

output "is_registry_addr_filter_useful" {
  value = length(data.huaweicloud_hss_image_registries.test_registry_addr.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_image_registries.test_registry_addr.data_list[*].registry_addr : v == local.registry_addr]
  )
}

# Filter by status
locals {
  status = data.huaweicloud_hss_image_registries.test.data_list[0].status
}

data "huaweicloud_hss_image_registries" "test_status" {
  enterprise_project_id = "all_granted_eps"
  status                = local.status
}

output "is_status_filter_useful" {
  value = length(data.huaweicloud_hss_image_registries.test_status.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_image_registries.test_status.data_list[*].status : v == local.status]
  )
}
`

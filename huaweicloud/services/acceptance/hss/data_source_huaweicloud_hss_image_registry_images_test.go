package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceHssImageRegistryImages_basic(t *testing.T) {
	dataSource := "data.huaweicloud_hss_image_registry_images.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceHssImageRegistryImages_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.namespace"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.image_name"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.image_id"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.image_digest"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.image_version"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.image_type"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.registry_id"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.registry_name"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.registry_type"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.latest_version"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.scan_status"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.image_size"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.latest_update_time"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.latest_sync_time"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.vul_num"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.unsafe_setting_num"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.malicious_file_num"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.scannable"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.is_multiple_arch"),

					resource.TestCheckOutput("is_namespace_filter_useful", "true"),
					resource.TestCheckOutput("is_image_name_filter_useful", "true"),
					resource.TestCheckOutput("is_image_version_filter_useful", "true"),
					resource.TestCheckOutput("is_registry_name_filter_useful", "true"),
					resource.TestCheckOutput("is_image_type_filter_useful", "true"),
					resource.TestCheckOutput("is_image_size_filter_useful", "true"),
					resource.TestCheckOutput("is_scan_status_filter_useful", "true"),
				),
			},
		},
	})
}

const testDataSourceDataSourceHssImageRegistryImages_basic string = `
data "huaweicloud_hss_image_registry_images" "test" {
  enterprise_project_id = "all_granted_eps"
}

# Filter by namespace
locals {
  namespace = data.huaweicloud_hss_image_registry_images.test.data_list[0].namespace
}

data "huaweicloud_hss_image_registry_images" "namespace_filter" {
  enterprise_project_id = "all_granted_eps"
  namespace             = local.namespace
}

output "is_namespace_filter_useful" {
  value = length(data.huaweicloud_hss_image_registry_images.namespace_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_image_registry_images.namespace_filter.data_list[*].namespace : v == local.namespace]
  )
}

# Filter by image_name
locals {
  image_name = data.huaweicloud_hss_image_registry_images.test.data_list[0].image_name
}

data "huaweicloud_hss_image_registry_images" "image_name_filter" {
  enterprise_project_id = "all_granted_eps"
  image_name            = local.image_name
}

output "is_image_name_filter_useful" {
  value = length(data.huaweicloud_hss_image_registry_images.image_name_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_image_registry_images.image_name_filter.data_list[*].image_name : v == local.image_name]
  )
}

# Filter by image_version
locals {
  image_version = data.huaweicloud_hss_image_registry_images.test.data_list[0].image_version
}

data "huaweicloud_hss_image_registry_images" "image_version_filter" {
  enterprise_project_id = "all_granted_eps"
  image_version         = local.image_version
}

output "is_image_version_filter_useful" {
  value = length(data.huaweicloud_hss_image_registry_images.image_version_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_image_registry_images.image_version_filter.data_list[*].image_version : v == local.image_version]
  )
}

# Filter by registry_name
locals {
  registry_name = data.huaweicloud_hss_image_registry_images.test.data_list[0].registry_name
}

data "huaweicloud_hss_image_registry_images" "registry_name_filter" {
  enterprise_project_id = "all_granted_eps"
  registry_name         = local.registry_name
}

output "is_registry_name_filter_useful" {
  value = length(data.huaweicloud_hss_image_registry_images.registry_name_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_image_registry_images.registry_name_filter.data_list[*].registry_name : v == local.registry_name]
  )
}

# Filter by image_type
locals {
  image_type = data.huaweicloud_hss_image_registry_images.test.data_list[0].image_type
}

data "huaweicloud_hss_image_registry_images" "image_type_filter" {
  enterprise_project_id = "all_granted_eps"
  image_type            = local.image_type
}

output "is_image_type_filter_useful" {
  value = length(data.huaweicloud_hss_image_registry_images.image_type_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_image_registry_images.image_type_filter.data_list[*].image_type : v == local.image_type]
  )
}

# Filter by image_size
locals {
  image_size = data.huaweicloud_hss_image_registry_images.test.data_list[0].image_size
}

data "huaweicloud_hss_image_registry_images" "image_size_filter" {
  enterprise_project_id = "all_granted_eps"
  image_size            = local.image_size
}

output "is_image_size_filter_useful" {
  value = length(data.huaweicloud_hss_image_registry_images.image_size_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_image_registry_images.image_size_filter.data_list[*].image_size : v == local.image_size]
  )
}

# Filter by scan_status
locals {
  scan_status = data.huaweicloud_hss_image_registry_images.test.data_list[0].scan_status
}

data "huaweicloud_hss_image_registry_images" "scan_status_filter" {
  enterprise_project_id = "all_granted_eps"
  scan_status           = local.scan_status
}

output "is_scan_status_filter_useful" {
  value = length(data.huaweicloud_hss_image_registry_images.scan_status_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_image_registry_images.scan_status_filter.data_list[*].scan_status : v == local.scan_status]
  )
}
`

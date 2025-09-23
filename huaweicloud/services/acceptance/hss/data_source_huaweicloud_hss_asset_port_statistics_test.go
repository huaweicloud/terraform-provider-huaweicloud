package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceAssetPortStatistics_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_hss_asset_port_statistics.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// This test case requires setting a host ID with host protection enabled.
			acceptance.TestAccPreCheckHSSHostProtectionHostId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceAssetPortStatistics_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.port"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.num"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.status"),

					resource.TestCheckOutput("is_port_filter_useful", "true"),
					resource.TestCheckOutput("is_port_string_filter_useful", "true"),
					resource.TestCheckOutput("is_type_filter_useful", "true"),
					resource.TestCheckOutput("is_category_filter_useful", "true"),
					resource.TestCheckOutput("is_sort_key_filter_useful", "true"),
					resource.TestCheckOutput("is_sort_dir_filter_useful", "true"),
					resource.TestCheckOutput("not_found_validation_pass", "true"),
				),
			},
		},
	})
}

const testDataSourceAssetPortStatistics_basic string = `
data "huaweicloud_hss_asset_port_statistics" "test" {}

# Filter using port.
locals {
  port = data.huaweicloud_hss_asset_port_statistics.test.data_list[0].port
}

data "huaweicloud_hss_asset_port_statistics" "port_filter" {
  port = local.port
}

output "is_port_filter_useful" {
  value = length(data.huaweicloud_hss_asset_port_statistics.port_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_asset_port_statistics.port_filter.data_list[*].port : v == local.port]
  )
}

# Filter using port_string.
locals {
  port_string = tostring(local.port)
}

data "huaweicloud_hss_asset_port_statistics" "port_string_filter" {
  port_string = local.port_string
}

output "is_port_string_filter_useful" {
  value = length(data.huaweicloud_hss_asset_port_statistics.port_string_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_asset_port_statistics.port_string_filter.data_list[*].port : v == local.port]
  )
}

# Filter using type.
locals {
  type = data.huaweicloud_hss_asset_port_statistics.test.data_list[0].type
}

data "huaweicloud_hss_asset_port_statistics" "type_filter" {
  type = local.type
}

output "is_type_filter_useful" {
  value = length(data.huaweicloud_hss_asset_port_statistics.type_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_asset_port_statistics.type_filter.data_list[*].type : v == local.type]
  )
}

# Filter using category.
locals {
  category = "host"
}

data "huaweicloud_hss_asset_port_statistics" "category_filter" {
  category = local.category
}

output "is_category_filter_useful" {
  value = length(data.huaweicloud_hss_asset_port_statistics.category_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_asset_port_statistics.category_filter.data_list[*].port : v != 0]
  )
}

# Filter using sort_key and sort_dir.
data "huaweicloud_hss_asset_port_statistics" "sort_desc_filter" {
  sort_key = "port"
  sort_dir = "desc"
}

data "huaweicloud_hss_asset_port_statistics" "sort_asc_filter" {
  sort_key = "port"
  sort_dir = "asc"
}

locals {
  asc_first_port = data.huaweicloud_hss_asset_port_statistics.sort_asc_filter.data_list[0].port
  desc_length = length(data.huaweicloud_hss_asset_port_statistics.sort_asc_filter.data_list)
  desc_last_port = data.huaweicloud_hss_asset_port_statistics.sort_desc_filter.data_list[local.desc_length - 1].port
}

output "is_sort_key_filter_useful" {
  value = length(data.huaweicloud_hss_asset_port_statistics.sort_desc_filter.data_list) > 0
}

output "is_sort_dir_filter_useful" {
  value = length(data.huaweicloud_hss_asset_port_statistics.sort_desc_filter.data_list) > 0 && local.asc_first_port == local.desc_last_port
}

# Filter using non existent port.
data "huaweicloud_hss_asset_port_statistics" "not_found" {
  port = 99999
}

output "not_found_validation_pass" {
  value = length(data.huaweicloud_hss_asset_port_statistics.not_found.data_list) == 0
}
`

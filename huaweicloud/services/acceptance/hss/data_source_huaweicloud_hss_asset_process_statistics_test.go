package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceAssetProcessStatistics_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_hss_asset_process_statistics.test"
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
				Config: testDataSourceAssetProcessStatistics_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.path"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.num"),

					resource.TestCheckOutput("is_path_filter_useful", "true"),
					resource.TestCheckOutput("is_category_filter_useful", "true"),
					resource.TestCheckOutput("not_found_validation_pass", "true"),
				),
			},
		},
	})
}

const testDataSourceAssetProcessStatistics_basic string = `
data "huaweicloud_hss_asset_process_statistics" "test" {}

# Filter using path.
locals {
  path = data.huaweicloud_hss_asset_process_statistics.test.data_list[0].path
}

data "huaweicloud_hss_asset_process_statistics" "path_filter" {
  path = local.path
}

output "is_path_filter_useful" {
  value = length(data.huaweicloud_hss_asset_process_statistics.path_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_asset_process_statistics.path_filter.data_list[*].path : v == local.path]
  )
}

# Filter using category.
locals {
  category = "host"
}

data "huaweicloud_hss_asset_process_statistics" "category_filter" {
  category = local.category
}

output "is_category_filter_useful" {
  value = length(data.huaweicloud_hss_asset_process_statistics.category_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_asset_process_statistics.category_filter.data_list[*].path : v != ""]
  )
}

# Filter using non existent path.
data "huaweicloud_hss_asset_process_statistics" "not_found" {
  path = "resource_not_found"
}

output "not_found_validation_pass" {
  value = length(data.huaweicloud_hss_asset_process_statistics.not_found.data_list) == 0
}
`

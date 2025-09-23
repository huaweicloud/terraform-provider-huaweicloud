package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceAssetUserStatistics_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_hss_asset_user_statistics.test"
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
				Config: testDataSourceAssetUserStatistics_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.user_name"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.num"),

					resource.TestCheckOutput("is_user_name_filter_useful", "true"),
					resource.TestCheckOutput("is_category_filter_useful", "true"),
					resource.TestCheckOutput("not_found_validation_pass", "true"),
				),
			},
		},
	})
}

const testDataSourceAssetUserStatistics_basic string = `
data "huaweicloud_hss_asset_user_statistics" "test" {}

# Filter using user_name.
locals {
  user_name = data.huaweicloud_hss_asset_user_statistics.test.data_list[0].user_name
}

data "huaweicloud_hss_asset_user_statistics" "user_name_filter" {
  user_name = local.user_name
}

output "is_user_name_filter_useful" {
  value = length(data.huaweicloud_hss_asset_user_statistics.user_name_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_asset_user_statistics.user_name_filter.data_list[*].user_name : v == local.user_name]
  )
}

# Filter using category.
locals {
  category = "host"
}

data "huaweicloud_hss_asset_user_statistics" "category_filter" {
  category = local.category
}

output "is_category_filter_useful" {
  value = length(data.huaweicloud_hss_asset_user_statistics.category_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_asset_user_statistics.category_filter.data_list[*].user_name : v != ""]
  )
}

# Filter using non existent user_name.
data "huaweicloud_hss_asset_user_statistics" "not_found" {
  user_name = "resource_not_found"
}

output "not_found_validation_pass" {
  value = length(data.huaweicloud_hss_asset_user_statistics.not_found.data_list) == 0
}
`

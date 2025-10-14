package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceAssetKernelModuleStatistics_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_hss_asset_kernel_module_statistics.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// This test case requires setting a host ID with host protection enabled.
			acceptance.TestAccPreCheckHSSHostProtectionHostId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAssetKernelModuleStatistics_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.num"),

					resource.TestCheckOutput("name_filter_useful", "true"),
					resource.TestCheckOutput("eps_filter_useful", "true"),
				),
			},
		},
	})
}

const testAccDataSourceAssetKernelModuleStatistics_basic = `
data "huaweicloud_hss_asset_kernel_module_statistics" "test" {}

locals {
  name = data.huaweicloud_hss_asset_kernel_module_statistics.test.data_list[0].name
}

data "huaweicloud_hss_asset_kernel_module_statistics" "name_filter" {
  name = local.name
}

output "name_filter_useful" {
  value = length(data.huaweicloud_hss_asset_kernel_module_statistics.name_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_asset_kernel_module_statistics.name_filter.data_list[*].name : v == local.name]
  )
}

data "huaweicloud_hss_asset_kernel_module_statistics" "eps_filter" {
  enterprise_project_id = "all_granted_eps"
}

output "eps_filter_useful" {
  value = length(data.huaweicloud_hss_asset_kernel_module_statistics.eps_filter.data_list) > 0
}
`

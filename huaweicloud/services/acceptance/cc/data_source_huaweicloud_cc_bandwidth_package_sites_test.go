package cc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCcBandwidthPackageSites_basic(t *testing.T) {
	dataSource := "data.huaweicloud_cc_bandwidth_package_sites.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceCcBandwidthPackageSites_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "bandwidth_package_sites.#"),
					resource.TestCheckResourceAttrSet(dataSource, "bandwidth_package_sites.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "bandwidth_package_sites.0.site_code"),
					resource.TestCheckResourceAttrSet(dataSource, "bandwidth_package_sites.0.site_type"),
					resource.TestCheckResourceAttrSet(dataSource, "bandwidth_package_sites.0.name_cn"),
					resource.TestCheckResourceAttrSet(dataSource, "bandwidth_package_sites.0.name_en"),
					resource.TestCheckResourceAttrSet(dataSource, "bandwidth_package_sites.0.region_id"),
					resource.TestCheckResourceAttrSet(dataSource, "bandwidth_package_sites.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSource, "bandwidth_package_sites.0.updated_at"),
					resource.TestCheckOutput("site_code_filter_is_useful", "true"),
					resource.TestCheckOutput("region_id_filter_is_useful", "true"),
					resource.TestCheckOutput("name_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceCcBandwidthPackageSites_basic() string {
	return `
data "huaweicloud_cc_bandwidth_package_sites" "test" {}

data "huaweicloud_cc_bandwidth_package_sites" "site_code_filter" {
  site_code = data.huaweicloud_cc_bandwidth_package_sites.test.bandwidth_package_sites[0].site_code
}
locals {
  site_code = data.huaweicloud_cc_bandwidth_package_sites.test.bandwidth_package_sites[0].site_code
}
output "site_code_filter_is_useful" {
  value = length(data.huaweicloud_cc_bandwidth_package_sites.site_code_filter.bandwidth_package_sites) > 0 && alltrue(
  [for v in data.huaweicloud_cc_bandwidth_package_sites.site_code_filter.bandwidth_package_sites[*].site_code :
  v == local.site_code]
  )
}

data "huaweicloud_cc_bandwidth_package_sites" "region_id_filter" {
  region_id = data.huaweicloud_cc_bandwidth_package_sites.test.bandwidth_package_sites[0].region_id
}
locals {
  region_id = data.huaweicloud_cc_bandwidth_package_sites.test.bandwidth_package_sites[0].region_id
}
output "region_id_filter_is_useful" {
  value = length(data.huaweicloud_cc_bandwidth_package_sites.region_id_filter.bandwidth_package_sites) > 0 && alltrue(
  [for v in data.huaweicloud_cc_bandwidth_package_sites.region_id_filter.bandwidth_package_sites[*].region_id :
  v == local.region_id]
  )
}

data "huaweicloud_cc_bandwidth_package_sites" "name_filter" {
  name = data.huaweicloud_cc_bandwidth_package_sites.test.bandwidth_package_sites[0].name_cn
}
locals {
  name = data.huaweicloud_cc_bandwidth_package_sites.test.bandwidth_package_sites[0].name_cn
}
output "name_filter_is_useful" {
  value = length(data.huaweicloud_cc_bandwidth_package_sites.name_filter.bandwidth_package_sites) > 0 && alltrue(
  [for v in data.huaweicloud_cc_bandwidth_package_sites.name_filter.bandwidth_package_sites[*].name_cn : v == local.name]
  )
}
`
}

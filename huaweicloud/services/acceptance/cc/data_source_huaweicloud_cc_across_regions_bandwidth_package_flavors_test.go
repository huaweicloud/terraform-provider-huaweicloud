package cc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCcAcrossRegionsBandwidthPackageFlavors_basic(t *testing.T) {
	dataSource := "data.huaweicloud_cc_across_regions_bandwidth_package_flavors.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceCcAcrossRegionsBandwidthPackageFlavors_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "region_specifications.#"),
					resource.TestCheckResourceAttrSet(dataSource, "region_specifications.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "region_specifications.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "region_specifications.0.en_name"),
					resource.TestCheckResourceAttrSet(dataSource, "region_specifications.0.local_region_id"),
					resource.TestCheckResourceAttrSet(dataSource, "region_specifications.0.remote_region_id"),
					resource.TestCheckResourceAttrSet(dataSource, "region_specifications.0.spec_codes.#"),
					resource.TestCheckResourceAttrSet(dataSource, "region_specifications.0.spec_codes.0.spec_code"),
					resource.TestCheckResourceAttrSet(dataSource, "region_specifications.0.spec_codes.0.billing_mode"),
					resource.TestCheckResourceAttrSet(dataSource, "region_specifications.0.spec_codes.0.max_bandwidth"),
					resource.TestCheckResourceAttrSet(dataSource, "region_specifications.0.spec_codes.0.mim_bandwidth"),
					resource.TestCheckOutput("local_region_id_filter_is_useful", "true"),
					resource.TestCheckOutput("remote_region_id_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceCcAcrossRegionsBandwidthPackageFlavors_basic() string {
	return `
data "huaweicloud_cc_across_regions_bandwidth_package_flavors" "test" {}

locals {
  flavors = data.huaweicloud_cc_across_regions_bandwidth_package_flavors.test.region_specifications
}

data "huaweicloud_cc_across_regions_bandwidth_package_flavors" "local_region_id_filter" {
  local_region_id = local.flavors[0].local_region_id
}
locals {
  local_region_id                = local.flavors[0].local_region_id
  local_region_id_filter_flavors = data.huaweicloud_cc_across_regions_bandwidth_package_flavors.local_region_id_filter
}
output "local_region_id_filter_is_useful" {
  value = length(local.local_region_id_filter_flavors.region_specifications) > 0 && alltrue(
  [for v in local.local_region_id_filter_flavors.region_specifications[*].local_region_id : v == local.local_region_id]
  )
}

data "huaweicloud_cc_across_regions_bandwidth_package_flavors" "remote_region_id_filter" {
  remote_region_id = local.flavors[0].remote_region_id
}
locals {
  remote_region_id                = local.flavors[0].remote_region_id
  remote_region_id_filter_flavors = data.huaweicloud_cc_across_regions_bandwidth_package_flavors.remote_region_id_filter
}
output "remote_region_id_filter_is_useful" {
  value = length(local.remote_region_id_filter_flavors.region_specifications) > 0 && alltrue(
  [for v in local.remote_region_id_filter_flavors.region_specifications[*].remote_region_id : v == local.remote_region_id]
  )
}
`
}

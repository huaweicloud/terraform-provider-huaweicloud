package eip

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceGlobalEipInternetBandwidthLimits_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_global_eip_internet_bandwidth_limits.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceGlobalEipInternetBandwidthLimits_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "internet_bandwidth_limits.#"),
					resource.TestCheckResourceAttrSet(dataSource, "internet_bandwidth_limits.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "internet_bandwidth_limits.0.charge_mode"),
					resource.TestCheckResourceAttrSet(dataSource, "internet_bandwidth_limits.0.min_size"),
					resource.TestCheckResourceAttrSet(dataSource, "internet_bandwidth_limits.0.max_size"),
					resource.TestCheckResourceAttrSet(dataSource, "internet_bandwidth_limits.0.type"),

					resource.TestCheckOutput("is_charge_mode_filter_useful", "true"),
					resource.TestCheckOutput("is_type_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceGlobalEipInternetBandwidthLimits_basic() string {
	return `
data "huaweicloud_global_eip_internet_bandwidth_limits" "test" {
  fields   = ["id", "charge_mode", "min_size", "ext_limit", "max_size", "type"]
  sort_key = "id"
  sort_dir = "asc"
}

# Filter by charge_mode.
locals {
  charge_mode = data.huaweicloud_global_eip_internet_bandwidth_limits.test.internet_bandwidth_limits[0].charge_mode
}

data "huaweicloud_global_eip_internet_bandwidth_limits" "charge_mode_filter" {
  charge_mode = local.charge_mode
}

locals {
  charge_mode_filter_result = [
    for v in data.huaweicloud_global_eip_internet_bandwidth_limits.charge_mode_filter.internet_bandwidth_limits[*].charge_mode
    : v == local.charge_mode
  ]
}

output "is_charge_mode_filter_useful" {
  value = alltrue(local.charge_mode_filter_result) && length(local.charge_mode_filter_result) > 0
}

# Filter by type.
locals {
  type = data.huaweicloud_global_eip_internet_bandwidth_limits.test.internet_bandwidth_limits[0].type
}

data "huaweicloud_global_eip_internet_bandwidth_limits" "type_filter" {
  type = local.type
}

locals {
  type_filter_result = [
    for v in data.huaweicloud_global_eip_internet_bandwidth_limits.type_filter.internet_bandwidth_limits[*].type
    : v == local.type
  ]
}

output "is_type_filter_useful" {
  value = alltrue(local.type_filter_result) && length(local.type_filter_result) > 0
}
`
}

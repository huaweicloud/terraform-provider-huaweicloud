package cc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCcAcrossAreaBandwidthPackageFlavors_basic(t *testing.T) {
	dataSource := "data.huaweicloud_cc_across_area_bandwidth_package_flavors.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceCcAcrossAreaBandwidthPackageFlavors_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "area_specifications.#"),
					resource.TestCheckResourceAttrSet(dataSource, "area_specifications.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "area_specifications.0.local_area_id"),
					resource.TestCheckResourceAttrSet(dataSource, "area_specifications.0.remote_area_id"),
					resource.TestCheckResourceAttrSet(dataSource, "area_specifications.0.spec_codes.#"),
					resource.TestCheckResourceAttrSet(dataSource, "area_specifications.0.spec_codes.0.spec_code"),
					resource.TestCheckResourceAttrSet(dataSource, "area_specifications.0.spec_codes.0.billing_mode"),
					resource.TestCheckResourceAttrSet(dataSource, "area_specifications.0.spec_codes.0.max_bandwidth"),
					resource.TestCheckResourceAttrSet(dataSource, "area_specifications.0.spec_codes.0.mim_bandwidth"),
					resource.TestCheckOutput("local_area_id_filter_is_useful", "true"),
					resource.TestCheckOutput("remote_area_id_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceCcAcrossAreaBandwidthPackageFlavors_basic() string {
	return `
data "huaweicloud_cc_across_area_bandwidth_package_flavors" "test" {}

locals {
  flavors = data.huaweicloud_cc_across_area_bandwidth_package_flavors.test.area_specifications
}

data "huaweicloud_cc_across_area_bandwidth_package_flavors" "local_area_id_filter" {
  local_area_id = [local.flavors[0].local_area_id]
}
locals {
  local_area_id                = local.flavors[0].local_area_id
  local_area_id_filter_flavors = data.huaweicloud_cc_across_area_bandwidth_package_flavors.local_area_id_filter
}
output "local_area_id_filter_is_useful"{
  value = length(local.local_area_id_filter_flavors.area_specifications) > 0 && alltrue(
  [for v in local.local_area_id_filter_flavors.area_specifications[*].local_area_id : v == local.local_area_id]
  )
}

data "huaweicloud_cc_across_area_bandwidth_package_flavors" "remote_area_id_filter" {
  remote_area_id = [local.flavors[0].remote_area_id]
}
locals {
  remote_area_id                = local.flavors[0].remote_area_id
  remote_area_id_filter_flavors = data.huaweicloud_cc_across_area_bandwidth_package_flavors.remote_area_id_filter
}
output "remote_area_id_filter_is_useful"{
  value = length(local.remote_area_id_filter_flavors.area_specifications) > 0 && alltrue(
  [for v in local.remote_area_id_filter_flavors.area_specifications[*].remote_area_id : v == local.remote_area_id]
  )
}
`
}

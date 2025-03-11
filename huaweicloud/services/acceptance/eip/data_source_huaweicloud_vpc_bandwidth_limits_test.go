package eip

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceVpcBandwidthLimits_basic(t *testing.T) {
	dataSource := "data.huaweicloud_vpc_bandwidth_limits.all"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceVpcBandwidthLimits_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "eip_bandwidth_limits.#"),
					resource.TestCheckResourceAttrSet(dataSource, "eip_bandwidth_limits.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "eip_bandwidth_limits.0.charge_mode"),
					resource.TestCheckResourceAttrSet(dataSource, "eip_bandwidth_limits.0.min_size"),
					resource.TestCheckResourceAttrSet(dataSource, "eip_bandwidth_limits.0.max_size"),

					resource.TestCheckOutput("charge_mode_filter_validation", "true"),
				),
			},
		},
	})
}

const testDataSourceVpcBandwidthLimits_basic = `
data "huaweicloud_vpc_bandwidth_limits" "all" {}

data "huaweicloud_vpc_bandwidth_limits" "test" {
  charge_mode = data.huaweicloud_vpc_bandwidth_limits.all.eip_bandwidth_limits[0].charge_mode
}

output "charge_mode_filter_validation" {
  value = alltrue([for v in data.huaweicloud_vpc_bandwidth_limits.test.eip_bandwidth_limits[*].charge_mode :
    v == data.huaweicloud_vpc_bandwidth_limits.all.eip_bandwidth_limits[0].charge_mode])
}`

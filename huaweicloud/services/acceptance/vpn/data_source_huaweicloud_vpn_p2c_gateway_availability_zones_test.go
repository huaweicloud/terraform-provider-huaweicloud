package vpn

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceVpnP2cGatewayAvailabilityZones_basic(t *testing.T) {
	dataSource := "data.huaweicloud_vpn_p2c_gateway_availability_zones.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceVpnP2cGatewayAvailabilityZones_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "availability_zones.#"),
				),
			},
		},
	})
}

func testDataSourceVpnP2cGatewayAvailabilityZones_basic() string {
	return `
data "huaweicloud_vpn_p2c_gateway_availability_zones" "test" {
  flavor = "Professional1"
}
`
}

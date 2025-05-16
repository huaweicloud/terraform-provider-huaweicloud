package vpn

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceVpnv51GatewayAvailabilityZones_basic(t *testing.T) {
	dataSource := "data.huaweicloud_vpnv51_gateway_availability_zones.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceVpnv51GatewayAvailabilityZones_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "availability_zones.#"),
					resource.TestCheckResourceAttrSet(dataSource, "availability_zones.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "availability_zones.0.public_border_group"),
					resource.TestCheckResourceAttrSet(dataSource, "availability_zones.0.available_specs.#"),
					resource.TestCheckResourceAttrSet(dataSource, "availability_zones.0.available_specs.0.flavor"),
					resource.TestCheckResourceAttrSet(dataSource, "availability_zones.0.available_specs.0.attachment_type"),
					resource.TestCheckResourceAttrSet(dataSource, "availability_zones.0.available_specs.0.ip_version"),
				),
			},
		},
	})
}

const testDataSourceVpnv51GatewayAvailabilityZones_basic = `data "huaweicloud_vpnv51_gateway_availability_zones" "test" {}`

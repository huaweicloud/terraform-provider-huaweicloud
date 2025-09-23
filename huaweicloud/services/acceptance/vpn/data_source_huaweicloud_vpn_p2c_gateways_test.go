package vpn

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceVpnP2cGateways_basic(t *testing.T) {
	dataSource := "data.huaweicloud_vpn_p2c_gateways.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckVPNP2cGatewayId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceVpnP2cGateways_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "p2c_vpn_gateways.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "p2c_vpn_gateways.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "p2c_vpn_gateways.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "p2c_vpn_gateways.0.vpc_id"),
					resource.TestCheckResourceAttrSet(dataSource, "p2c_vpn_gateways.0.connect_subnet"),
					resource.TestCheckResourceAttrSet(dataSource, "p2c_vpn_gateways.0.flavor"),
					resource.TestCheckResourceAttrSet(dataSource, "p2c_vpn_gateways.0.availability_zone_ids.#"),
					resource.TestCheckResourceAttrSet(dataSource, "p2c_vpn_gateways.0.eip.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "p2c_vpn_gateways.0.eip.0.ip_version"),
					resource.TestCheckResourceAttrSet(dataSource, "p2c_vpn_gateways.0.eip.0.ip_billing_info"),
					resource.TestCheckResourceAttrSet(dataSource, "p2c_vpn_gateways.0.eip.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "p2c_vpn_gateways.0.eip.0.ip_address"),
					resource.TestCheckResourceAttrSet(dataSource, "p2c_vpn_gateways.0.eip.0.charge_mode"),
					resource.TestCheckResourceAttrSet(dataSource, "p2c_vpn_gateways.0.eip.0.bandwidth_id"),
					resource.TestCheckResourceAttrSet(dataSource, "p2c_vpn_gateways.0.eip.0.bandwidth_size"),
					resource.TestCheckResourceAttrSet(dataSource, "p2c_vpn_gateways.0.eip.0.bandwidth_name"),
					resource.TestCheckResourceAttrSet(dataSource, "p2c_vpn_gateways.0.eip.0.bandwidth_billing_info"),
					resource.TestCheckResourceAttrSet(dataSource, "p2c_vpn_gateways.0.eip.0.share_type"),
					resource.TestCheckResourceAttrSet(dataSource, "p2c_vpn_gateways.0.enterprise_project_id"),
					resource.TestCheckResourceAttrSet(dataSource, "p2c_vpn_gateways.0.availability_zone_ids.#"),
				),
			},
		},
	})
}

const testDataSourceVpnP2cGateways_basic = `data "huaweicloud_vpn_p2c_gateways" "test" {}`

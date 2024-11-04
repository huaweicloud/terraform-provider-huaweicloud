package vpn

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceVpnP2cGatewayConnections_basic(t *testing.T) {
	dataSource := "data.huaweicloud_vpn_p2c_gateway_connections.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckVPNP2cGatewayId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceVpnP2cGatewayConnections_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "connections.0.connection_id"),
					resource.TestCheckResourceAttrSet(dataSource, "connections.0.client_ip"),
					resource.TestCheckResourceAttrSet(dataSource, "connections.0.inbound_packets"),
					resource.TestCheckResourceAttrSet(dataSource, "connections.0.inbound_bytes"),
					resource.TestCheckResourceAttrSet(dataSource, "connections.0.outbound_bytes"),
					resource.TestCheckResourceAttrSet(dataSource, "connections.0.connection_established_time"),
					resource.TestCheckResourceAttrSet(dataSource, "connections.0.timestamp"),
				),
			},
		},
	})
}

func testDataSourceVpnP2cGatewayConnections_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_vpn_p2c_gateway_connections" "test" {
  p2c_gateway_id = "%[1]s"
}
`, acceptance.HW_VPN_P2C_GATEWAY_ID)
}

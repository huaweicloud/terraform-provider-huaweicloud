package vpn

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceVpnGatewayRouteTables_basic(t *testing.T) {
	dataSource := "data.huaweicloud_vpn_gateway_route_tables.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckVPNGatewayId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceVpnGatewayRouteTables_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "routing_table.#"),
					resource.TestCheckResourceAttrSet(dataSource, "routing_table.0.destination"),
					resource.TestCheckResourceAttrSet(dataSource, "routing_table.0.nexthop"),
					resource.TestCheckResourceAttrSet(dataSource, "routing_table.0.outbound_interface_ip"),
					resource.TestCheckResourceAttrSet(dataSource, "routing_table.0.origin"),
					resource.TestCheckResourceAttrSet(dataSource, "routing_table.0.med"),
					resource.TestCheckResourceAttrSet(dataSource, "routing_table.0.nexthop_resource.#"),
					resource.TestCheckResourceAttrSet(dataSource, "routing_table.0.nexthop_resource.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "routing_table.0.nexthop_resource.0.type"),
				),
			},
		},
	})
}

func testDataSourceVpnGatewayRouteTables_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_vpn_gateway_route_tables" "test" {
  vgw_id                      = "%[1]s"
  is_include_nexthop_resource = true
}
`, acceptance.HW_VPN_GATEWAY_ID)
}

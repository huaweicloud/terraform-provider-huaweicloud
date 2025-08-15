package dc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDcGlobalGatewayRouteTables_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()
	dataSource := "data.huaweicloud_dc_global_gateway_route_tables.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDcDirectConnection(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDcGlobalGatewayRouteTables_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "gdgw_routetables.#"),
					resource.TestCheckResourceAttrSet(dataSource, "gdgw_routetables.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "gdgw_routetables.0.gateway_id"),
					resource.TestCheckResourceAttrSet(dataSource, "gdgw_routetables.0.nexthop"),
					resource.TestCheckResourceAttrSet(dataSource, "gdgw_routetables.0.obtain_mode"),
					resource.TestCheckResourceAttrSet(dataSource, "gdgw_routetables.0.destination"),
					resource.TestCheckResourceAttrSet(dataSource, "gdgw_routetables.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "gdgw_routetables.0.address_family"),
					resource.TestCheckResourceAttrSet(dataSource, "gdgw_routetables.0.type"),

					resource.TestCheckOutput("nexthop_filter_is_useful", "true"),
					resource.TestCheckOutput("destination_filter_is_useful", "true"),
					resource.TestCheckOutput("address_family_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceDcGlobalGatewayRouteTables_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_dc_global_gateway_route_tables" "test" {
  depends_on = [huaweicloud_dc_global_gateway_route_table.test]

  gdgw_id = huaweicloud_dc_global_gateway.test.id
}

locals {
  nexthop = data.huaweicloud_dc_global_gateway_route_tables.test.gdgw_routetables[0].nexthop
}
data "huaweicloud_dc_global_gateway_route_tables" "nexthop_filter" {
  depends_on = [huaweicloud_dc_global_gateway_route_table.test]

  gdgw_id = huaweicloud_dc_global_gateway.test.id
  nexthop = [local.nexthop]
}
output "nexthop_filter_is_useful" {
  value = length(data.huaweicloud_dc_global_gateway_route_tables.nexthop_filter.gdgw_routetables) > 0 && alltrue(
  [for v in data.huaweicloud_dc_global_gateway_route_tables.nexthop_filter.gdgw_routetables[*].nexthop : v == local.nexthop]
  )
}

locals {
  destination = data.huaweicloud_dc_global_gateway_route_tables.test.gdgw_routetables[0].destination
}
data "huaweicloud_dc_global_gateway_route_tables" "destination_filter" {
  depends_on = [huaweicloud_dc_global_gateway_route_table.test]

  gdgw_id     = huaweicloud_dc_global_gateway.test.id
  destination = [local.destination]
}
output "destination_filter_is_useful" {
  value = length(data.huaweicloud_dc_global_gateway_route_tables.destination_filter.gdgw_routetables) > 0 && alltrue(
  [for v in data.huaweicloud_dc_global_gateway_route_tables.destination_filter.gdgw_routetables[*].destination :
  v == local.destination]
  )
}

locals {
  address_family = data.huaweicloud_dc_global_gateway_route_tables.test.gdgw_routetables[0].address_family
}
data "huaweicloud_dc_global_gateway_route_tables" "address_family_filter" {
  depends_on = [huaweicloud_dc_global_gateway_route_table.test]

  gdgw_id        = huaweicloud_dc_global_gateway.test.id
  address_family = [local.address_family]
}
output "address_family_filter_is_useful" {
  value = length(data.huaweicloud_dc_global_gateway_route_tables.address_family_filter.gdgw_routetables) > 0 && alltrue(
  [for v in data.huaweicloud_dc_global_gateway_route_tables.address_family_filter.gdgw_routetables[*].address_family :
  v == local.address_family]
  )
}
`, testResourceDcGlobalGatewayRouteTable_basic(name))
}

package dc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDcConnectGateways_basic(t *testing.T) {
	dataSource := "data.huaweicloud_dc_connect_gateways.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDcConnectGateways_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "connect_gateways.#"),
					resource.TestCheckResourceAttrSet(dataSource, "connect_gateways.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "connect_gateways.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "connect_gateways.0.description"),
					resource.TestCheckResourceAttrSet(dataSource, "connect_gateways.0.address_family"),
					resource.TestCheckResourceAttrSet(dataSource, "connect_gateways.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "connect_gateways.0.bgp_asn"),
					resource.TestCheckResourceAttrSet(dataSource, "connect_gateways.0.current_geip_count"),
					resource.TestCheckResourceAttrSet(dataSource, "connect_gateways.0.created_time"),
					resource.TestCheckOutput("id_filter_is_useful", "true"),
					resource.TestCheckOutput("name_filter_is_useful", "true"),
					resource.TestCheckOutput("sort_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceDcConnectGateways_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_dc_connect_gateways" "test" {
  depends_on = [huaweicloud_dc_connect_gateway.test]
}

locals {
  id = huaweicloud_dc_connect_gateway.test.id
}
data "huaweicloud_dc_connect_gateways" "id_filter" {
  connect_gateway_id = [huaweicloud_dc_connect_gateway.test.id]
}
output "id_filter_is_useful" {
  value = length(data.huaweicloud_dc_connect_gateways.id_filter.connect_gateways) > 0 && alltrue(
  [for v in data.huaweicloud_dc_connect_gateways.id_filter.connect_gateways[*].id : v == local.id]
  )
}

data "huaweicloud_dc_connect_gateways" "name_filter" {
  depends_on = [huaweicloud_dc_connect_gateway.test]

  name = ["%[2]s"]
}
output "name_filter_is_useful" {
  value = length(data.huaweicloud_dc_connect_gateways.name_filter.connect_gateways) > 0 && alltrue(
  [for v in data.huaweicloud_dc_connect_gateways.name_filter.connect_gateways[*].name : v == "%[2]s"]
  )
}

data "huaweicloud_dc_connect_gateways" "sort_filter" {
  depends_on = [huaweicloud_dc_connect_gateway.test]

  sort_key = "name"
  sort_dir = ["desc"]
}
output "sort_filter_is_useful" {
  value = length(data.huaweicloud_dc_connect_gateways.sort_filter.connect_gateways) > 0
}
`, testResourceDcConnectGateway_basic(name), name)
}

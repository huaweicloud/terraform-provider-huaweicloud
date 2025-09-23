package dc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDcConnectGatewayGeips_basic(t *testing.T) {
	dataSource := "data.huaweicloud_dc_connect_gateway_geips.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDcConnectGatewayId(t)
			acceptance.TestAccPreCheckGlobalEipId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDcConnectGatewayGeips_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "global_eips.#"),
					resource.TestCheckResourceAttrSet(dataSource, "global_eips.0.global_eip_id"),
					resource.TestCheckResourceAttrSet(dataSource, "global_eips.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "global_eips.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "global_eips.0.cidr"),
					resource.TestCheckResourceAttrSet(dataSource, "global_eips.0.address_family"),
					resource.TestCheckResourceAttrSet(dataSource, "global_eips.0.ie_vtep_ip"),
					resource.TestCheckResourceAttrSet(dataSource, "global_eips.0.created_time"),
					resource.TestCheckOutput("global_eip_id_filter_is_useful", "true"),
					resource.TestCheckOutput("status_filter_is_useful", "true"),
					resource.TestCheckOutput("sort_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceDcConnectGatewayGeips_basic() string {
	return fmt.Sprintf(`
%[1]s


data "huaweicloud_dc_connect_gateway_geips" "test" {
  depends_on = [huaweicloud_dc_connect_gateway_geip_associate.test]

  connect_gateway_id = "%[2]s"
}

locals {
  global_eip_id = "%[3]s"
}
data "huaweicloud_dc_connect_gateway_geips" "global_eip_id_filter" {
  depends_on = [huaweicloud_dc_connect_gateway_geip_associate.test]

  connect_gateway_id = "%[2]s"
  global_eip_id      = ["%[3]s"]
}
output "global_eip_id_filter_is_useful" {
  value = length(data.huaweicloud_dc_connect_gateway_geips.global_eip_id_filter.global_eips) > 0 && alltrue(
  [for v in data.huaweicloud_dc_connect_gateway_geips.global_eip_id_filter.global_eips[*].global_eip_id : v == local.global_eip_id]
  )
}

locals {
  status = data.huaweicloud_dc_connect_gateway_geips.test.global_eips[0].status
}
data "huaweicloud_dc_connect_gateway_geips" "status_filter" {
  depends_on = [huaweicloud_dc_connect_gateway_geip_associate.test]

  connect_gateway_id = "%[2]s"
  status             = [data.huaweicloud_dc_connect_gateway_geips.test.global_eips[0].status]
}
output "status_filter_is_useful" {
  value = length(data.huaweicloud_dc_connect_gateway_geips.status_filter.global_eips) > 0 && alltrue(
  [for v in data.huaweicloud_dc_connect_gateway_geips.status_filter.global_eips[*].status : v == local.status]
  )
}

data "huaweicloud_dc_connect_gateway_geips" "sort_filter" {
  depends_on = [huaweicloud_dc_connect_gateway_geip_associate.test]

  connect_gateway_id = "%[2]s"
  sort_key           = "created_time"
  sort_dir           = ["desc"]
}
output "sort_filter_is_useful" {
  value = length(data.huaweicloud_dc_connect_gateway_geips.sort_filter.global_eips) > 0
}
`, testResourceDcConnectGatewayGeipAssociate_basic(), acceptance.HW_DC_CONNECT_GATEWAY_ID, acceptance.HW_GLOBAL_EIP_ID)
}

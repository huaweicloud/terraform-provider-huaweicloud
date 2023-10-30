package vpn

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceGateway_basic(t *testing.T) {
	rName := "data.huaweicloud_vpn_gateways.test"
	dc := acceptance.InitDataSourceCheck(rName)
	name := acceptance.RandomAccResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceGateway_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "gateways.0.id"),
					resource.TestCheckResourceAttrSet(rName, "gateways.0.name"),
					resource.TestCheckResourceAttrSet(rName, "gateways.0.network_type"),
					resource.TestCheckResourceAttrSet(rName, "gateways.0.status"),
					resource.TestCheckResourceAttrSet(rName, "gateways.0.attachment_type"),
					resource.TestCheckResourceAttrSet(rName, "gateways.0.vpc_id"),
					resource.TestCheckResourceAttrSet(rName, "gateways.0.local_subnets.#"),
					resource.TestCheckResourceAttrSet(rName, "gateways.0.connect_subnet"),
					resource.TestCheckResourceAttrSet(rName, "gateways.0.bgp_asn"),
					resource.TestCheckResourceAttrSet(rName, "gateways.0.flavor"),
					resource.TestCheckResourceAttrSet(rName, "gateways.0.availability_zones.#"),
					resource.TestCheckResourceAttrSet(rName, "gateways.0.connection_number"),
					resource.TestCheckResourceAttrSet(rName, "gateways.0.used_connection_number"),
					resource.TestCheckResourceAttrSet(rName, "gateways.0.used_connection_group"),
					resource.TestCheckResourceAttrSet(rName, "gateways.0.enterprise_project_id"),
					resource.TestCheckResourceAttrSet(rName, "gateways.0.eips.#"),
					resource.TestCheckResourceAttrSet(rName, "gateways.0.created_at"),
					resource.TestCheckResourceAttrSet(rName, "gateways.0.updated_at"),
					resource.TestCheckResourceAttrSet(rName, "gateways.0.access_vpc_id"),
					resource.TestCheckResourceAttrSet(rName, "gateways.0.access_subnet_id"),
					resource.TestCheckResourceAttrSet(rName, "gateways.0.ha_mode"),

					resource.TestCheckOutput("name_filter_is_useful", "true"),

					resource.TestCheckOutput("gateway_id_filter_is_useful", "true"),

					resource.TestCheckOutput("network_type_filter_is_useful", "true"),

					resource.TestCheckOutput("attachment_type_filter_is_useful", "true"),

					resource.TestCheckOutput("enterprise_project_id_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDatasourceGateway_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_vpn_gateways" "test" {
  depends_on = [huaweicloud_vpn_gateway.test]
}

data "huaweicloud_vpn_gateways" "name_filter" {
  name = "%[2]s"

  depends_on = [huaweicloud_vpn_gateway.test]
}

output "name_filter_is_useful" {
  value = length(data.huaweicloud_vpn_gateways.name_filter.gateways) > 0 && alltrue(
    [for v in data.huaweicloud_vpn_gateways.name_filter.gateways[*].name : v == "%[2]s"]
  )  
}

data "huaweicloud_vpn_gateways" "gateway_id_filter" {
  gateway_id = huaweicloud_vpn_gateway.test.id
}

locals {
  gateway_id = huaweicloud_vpn_gateway.test.id
}

output "gateway_id_filter_is_useful" {
  value = length(data.huaweicloud_vpn_gateways.gateway_id_filter.gateways) > 0 && alltrue(
    [for v in data.huaweicloud_vpn_gateways.gateway_id_filter.gateways[*].id : v == local.gateway_id]
  )  
}

data "huaweicloud_vpn_gateways" "network_type_filter" {
  network_type = "public"

  depends_on = [huaweicloud_vpn_gateway.test]
}
output "network_type_filter_is_useful" {
  value = length(data.huaweicloud_vpn_gateways.network_type_filter.gateways) > 0 && alltrue(
    [for v in data.huaweicloud_vpn_gateways.network_type_filter.gateways[*].network_type : v == "public"]
  )  
}

data "huaweicloud_vpn_gateways" "attachment_type_filter" {
  attachment_type = "vpc"

  depends_on = [huaweicloud_vpn_gateway.test]
}
output "attachment_type_filter_is_useful" {
  value = length(data.huaweicloud_vpn_gateways.attachment_type_filter.gateways) > 0 && alltrue(
    [for v in data.huaweicloud_vpn_gateways.attachment_type_filter.gateways[*].attachment_type : v == "vpc"]
  )  
}

data "huaweicloud_vpn_gateways" "enterprise_project_id_filter" {
  enterprise_project_id = "0"

  depends_on = [huaweicloud_vpn_gateway.test]
}
output "enterprise_project_id_filter_is_useful" {
  value = length(data.huaweicloud_vpn_gateways.enterprise_project_id_filter.gateways) > 0 && alltrue(
    [for v in data.huaweicloud_vpn_gateways.enterprise_project_id_filter.gateways[*].enterprise_project_id : v == "0"]
  )  
}
`, testGateway_basic(name), name)
}

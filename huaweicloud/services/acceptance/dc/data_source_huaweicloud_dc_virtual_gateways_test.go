package dc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceDCVirtualGateways_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()
	rName := "data.huaweicloud_dc_virtual_gateways.test"
	dc := acceptance.InitDataSourceCheck(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceDCVirtualGateways_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "virtual_gateways.0.id"),
					resource.TestCheckResourceAttrSet(rName, "virtual_gateways.0.vpc_id"),
					resource.TestCheckResourceAttrSet(rName, "virtual_gateways.0.name"),
					resource.TestCheckResourceAttrSet(rName, "virtual_gateways.0.type"),
					resource.TestCheckResourceAttrSet(rName, "virtual_gateways.0.status"),
					resource.TestCheckResourceAttrSet(rName, "virtual_gateways.0.asn"),
					resource.TestCheckResourceAttrSet(rName, "virtual_gateways.0.local_ep_group.#"),
					resource.TestCheckResourceAttrSet(rName, "virtual_gateways.0.description"),
					resource.TestCheckResourceAttrSet(rName, "virtual_gateways.0.enterprise_project_id"),

					resource.TestCheckOutput("virtual_gateway_id_filter_is_useful", "true"),

					resource.TestCheckOutput("name_filter_is_useful", "true"),

					resource.TestCheckOutput("vpc_id_filter_is_useful", "true"),

					resource.TestCheckOutput("enterprise_project_id_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDatasourceDCVirtualGateways_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_dc_virtual_gateways" "test" {
  depends_on = [huaweicloud_dc_virtual_gateway.test]
}

data "huaweicloud_dc_virtual_gateways" "virtual_gateway_id_filter" {
  virtual_gateway_id = huaweicloud_dc_virtual_gateway.test.id
}

locals {
  virtual_gateway_id = huaweicloud_dc_virtual_gateway.test.id
}

output "virtual_gateway_id_filter_is_useful" {
  value = length(data.huaweicloud_dc_virtual_gateways.virtual_gateway_id_filter.virtual_gateways) > 0 && alltrue(
    [for v in data.huaweicloud_dc_virtual_gateways.virtual_gateway_id_filter.virtual_gateways[*].id : v == 
  local.virtual_gateway_id]
  )  
}

data "huaweicloud_dc_virtual_gateways" "name_filter" {
  name = huaweicloud_dc_virtual_gateway.test.name
}

locals {
  name = huaweicloud_dc_virtual_gateway.test.name
}

output "name_filter_is_useful" {
  value = length(data.huaweicloud_dc_virtual_gateways.name_filter.virtual_gateways) > 0 && alltrue(
    [for v in data.huaweicloud_dc_virtual_gateways.name_filter.virtual_gateways[*].name : v == local.name]
  )
}

data "huaweicloud_dc_virtual_gateways" "vpc_id_filter" {
  vpc_id = huaweicloud_dc_virtual_gateway.test.vpc_id
}

locals {
  vpc_id = huaweicloud_dc_virtual_gateway.test.vpc_id
}

output "vpc_id_filter_is_useful" {
  value = length(data.huaweicloud_dc_virtual_gateways.vpc_id_filter.virtual_gateways) > 0 && alltrue(
    [for v in data.huaweicloud_dc_virtual_gateways.vpc_id_filter.virtual_gateways[*].vpc_id : v == local.vpc_id]
  )
}

data "huaweicloud_dc_virtual_gateways" "enterprise_project_id_filter" {
  enterprise_project_id = huaweicloud_dc_virtual_gateway.test.enterprise_project_id
}

locals {
  enterprise_project_id = huaweicloud_dc_virtual_gateway.test.enterprise_project_id
}

output "enterprise_project_id_filter_is_useful" {
  value = length(data.huaweicloud_dc_virtual_gateways.enterprise_project_id_filter.virtual_gateways) > 0 && alltrue(
    [for v in data.huaweicloud_dc_virtual_gateways.enterprise_project_id_filter.virtual_gateways[*].
  enterprise_project_id : v == local.enterprise_project_id]
  )
}
`, testAccVirtualGateway_basic(name, acceptance.RandomCidr()))
}

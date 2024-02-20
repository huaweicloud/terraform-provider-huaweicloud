package eip

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccVPCInternetGateways_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_vpc_internet_gateways.all"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVPCInternetGatewaysDataSource_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "vpc_igws.#"),

					resource.TestCheckOutput("is_name_filter_useful", "true"),
					resource.TestCheckOutput("is_id_filter_useful", "true"),
					resource.TestCheckOutput("is_vpc_id_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccVPCInternetGatewaysDataSource_basic() string {
	rNameWithDash := acceptance.RandomAccResourceNameWithDash()

	return fmt.Sprintf(`
%[1]s

data "huaweicloud_vpc_internet_gateways" "all" {
  depends_on = [huaweicloud_vpc_internet_gateway.test]
}

// filter by name
data "huaweicloud_vpc_internet_gateways" "filter_by_name" {
  igw_name = huaweicloud_vpc_internet_gateway.test.name
}

locals {
  filter_result_by_name = [for v in data.huaweicloud_vpc_internet_gateways.filter_by_name.vpc_igws[*].name : 
    v == "%[2]s"]
}

output "is_name_filter_useful" {
  value = length(local.filter_result_by_name) > 0 && alltrue(local.filter_result_by_name) 
}

// filter by ID
data "huaweicloud_vpc_internet_gateways" "filter_by_id" {
  igw_id = huaweicloud_vpc_internet_gateway.test.id
}

locals {
  filter_result_by_id = [for v in data.huaweicloud_vpc_internet_gateways.filter_by_id.vpc_igws[*].id : 
    v == huaweicloud_vpc_internet_gateway.test.id]
}

output "is_id_filter_useful" {
  value = length(local.filter_result_by_id) == 1 && alltrue(local.filter_result_by_id) 
}

// filter by vpc ID
data "huaweicloud_vpc_internet_gateways" "filter_by_vpc_id" {
  vpc_id = huaweicloud_vpc_internet_gateway.test.vpc_id
}

locals {
  filter_result_by_vpc_id = [for v in data.huaweicloud_vpc_internet_gateways.filter_by_vpc_id.vpc_igws[*].id : 
    v == huaweicloud_vpc_internet_gateway.test.id]
}

output "is_vpc_id_filter_useful" {
  value = length(local.filter_result_by_vpc_id) == 1 && alltrue(local.filter_result_by_vpc_id) 
}
`, testAccIGW_basic(rNameWithDash), rNameWithDash)
}

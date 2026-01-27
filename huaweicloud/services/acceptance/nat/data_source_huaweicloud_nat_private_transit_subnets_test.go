package nat

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccDatasourcePrivateTransitSubnets_basic(t *testing.T) {
	var (
		name           = acceptance.RandomAccResourceName()
		dataSourceName = "data.huaweicloud_nat_private_transit_subnets.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckProjectID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourcePrivateTransitSubnets_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),

					resource.TestCheckOutput("transit_subnet_id_filter_is_useful", "true"),
					resource.TestCheckOutput("transit_subnet_name_filter_is_useful", "true"),
					resource.TestCheckOutput("transit_subnet_description_filter_is_useful", "true"),
					resource.TestCheckOutput("subnet_id_filter_is_useful", "true"),
					resource.TestCheckOutput("vpc_id_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourcePrivateTransitSubnets_base(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_nat_private_transit_subnet" "test" {
  name                 = "%[2]s"
  description          = "Created by acc test"
  virsubnet_id         = huaweicloud_vpc_subnet.test.id
  virsubnet_project_id = "%[3]s"

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, common.TestVpc(name), name, acceptance.HW_PROJECT_ID)
}

func testAccDatasourcePrivateTransitSubnets_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_nat_private_transit_subnets" "test" {
  depends_on = [
    huaweicloud_nat_private_transit_subnet.test
  ]
}

locals {
  transit_subnet_id = data.huaweicloud_nat_private_transit_subnets.test.transit_subnets[0].id
}

data "huaweicloud_nat_private_transit_subnets" "filter_by_transit_subnet_id" {
  ids = [local.transit_subnet_id]
}

locals {
  transit_subnet_id_filter_result = [
    for v in data.huaweicloud_nat_private_transit_subnets.filter_by_transit_subnet_id.transit_subnets[*].id : 
    v == local.transit_subnet_id
  ]
}

output "transit_subnet_id_filter_is_useful" {
  value = alltrue(local.transit_subnet_id_filter_result) && length(local.transit_subnet_id_filter_result) > 0
}

locals {
  transit_subnet_name = data.huaweicloud_nat_private_transit_subnets.test.transit_subnets[0].name
}

data "huaweicloud_nat_private_transit_subnets" "filter_by_transit_subnet_name" {
  names = [local.transit_subnet_name]
}

locals {
  transit_subnet_name_filter_result = [
    for v in data.huaweicloud_nat_private_transit_subnets.filter_by_transit_subnet_name.transit_subnets[*].name : 
    v == local.transit_subnet_name
  ]
}

output "transit_subnet_name_filter_is_useful" {
  value = alltrue(local.transit_subnet_name_filter_result) && length(local.transit_subnet_name_filter_result) > 0
}

locals {
  transit_subnet_description = data.huaweicloud_nat_private_transit_subnets.test.transit_subnets[0].description
}

data "huaweicloud_nat_private_transit_subnets" "filter_by_transit_subnet_description" {
  descriptions = [local.transit_subnet_description]
}

locals {
  transit_subnet_description_filter_result = [
    for v in data.huaweicloud_nat_private_transit_subnets.filter_by_transit_subnet_description.transit_subnets[*].description : 
    v == local.transit_subnet_description
  ]
}

output "transit_subnet_description_filter_is_useful" {
  value = alltrue(local.transit_subnet_description_filter_result) && length(local.transit_subnet_description_filter_result) > 0
}

locals {
  subnet_id = data.huaweicloud_nat_private_transit_subnets.test.transit_subnets[0].virsubnet_id
}

data "huaweicloud_nat_private_transit_subnets" "filter_by_subnet_id" {
  virsubnet_ids = [local.subnet_id]
}

locals {
  subnet_id_filter_result = [
    for v in data.huaweicloud_nat_private_transit_subnets.filter_by_subnet_id.transit_subnets[*].virsubnet_id : 
    v == local.subnet_id
  ]
}

output "subnet_id_filter_is_useful" {
  value = alltrue(local.subnet_id_filter_result) && length(local.subnet_id_filter_result) > 0
}

locals {
  vpc_id = data.huaweicloud_nat_private_transit_subnets.test.transit_subnets[0].virsubnet_project_id
}

data "huaweicloud_nat_private_transit_subnets" "filter_by_vpc_id" {
  vpc_ids = [local.vpc_id]
}

locals {
  vpc_id_filter_result = [
    for v in data.huaweicloud_nat_private_transit_subnets.filter_by_vpc_id.transit_subnets[*].project_id : 
    v == local.vpc_id
  ]
}

output "vpc_id_filter_is_useful" {
  value = alltrue(local.vpc_id_filter_result) && length(local.vpc_id_filter_result) > 0
}
`, testDataSourcePrivateTransitSubnets_base(name))
}

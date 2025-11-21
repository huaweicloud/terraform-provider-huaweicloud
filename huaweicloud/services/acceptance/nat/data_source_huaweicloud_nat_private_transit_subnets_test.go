package nat

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourcePrivateTransitSubnets_basic(t *testing.T) {
	var (
		name           = acceptance.RandomAccResourceName()
		dataSourceName = "data.huaweicloud_nat_private_transit_subnets.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)

		filterByTransitSubnetId          = "data.huaweicloud_nat_private_transit_subnets.filter_by_transit_subnet_id"
		filterByTransitSubnetName        = "data.huaweicloud_nat_private_transit_subnets.filter_by_transit_subnet_name"
		filterByTransitSubnetDescription = "data.huaweicloud_nat_private_transit_subnets.filter_by_transit_subnet_description"
		filterByTransitSubnetSubnetId    = "data.huaweicloud_nat_private_transit_subnets.filter_by_transit_subnet_subnet_id"
		filterByTransitSubnetProjectId   = "data.huaweicloud_nat_private_transit_subnets.filter_by_transit_subnet_project_id"

		dcFilterByTransitSubnetId          = acceptance.InitDataSourceCheck(filterByTransitSubnetId)
		dcFilterByTransitSubnetName        = acceptance.InitDataSourceCheck(filterByTransitSubnetName)
		dcFilterByTransitSubnetDescription = acceptance.InitDataSourceCheck(filterByTransitSubnetDescription)
		dcFilterByTransitSubnetSubnetId    = acceptance.InitDataSourceCheck(filterByTransitSubnetSubnetId)
		dcFilterByTransitSubnetProjectId   = acceptance.InitDataSourceCheck(filterByTransitSubnetProjectId)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourcePrivateTransitSubnets_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),

					dcFilterByTransitSubnetId.CheckResourceExists(),
					resource.TestCheckOutput("transit_subnet_id_filter_is_useful", "true"),

					dcFilterByTransitSubnetName.CheckResourceExists(),
					resource.TestCheckOutput("transit_subnet_name_filter_is_useful", "true"),

					dcFilterByTransitSubnetDescription.CheckResourceExists(),
					resource.TestCheckOutput("transit_subnet_description_filter_is_useful", "true"),

					dcFilterByTransitSubnetSubnetId.CheckResourceExists(),
					resource.TestCheckOutput("transit_subnet_subnet_id_filter_is_useful", "true"),

					dcFilterByTransitSubnetProjectId.CheckResourceExists(),
					resource.TestCheckOutput("transit_subnet_project_id_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceNatPrivateTransitSubnets_base(name string) string {
	return fmt.Sprintf(`

data "huaweicloud_vpc_subnet" "test" {
  name = "subnet-default"
}

resource "huaweicloud_nat_private_transit_subnet" "test" {
  name                 = "%[1]s"
  description          = "Created by acc test"
  virsubnet_id         = data.huaweicloud_vpc_subnet.test.id
  virsubnet_project_id = "%[2]s"

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, name, acceptance.HW_PROJECT_ID)
}

func testAccDatasourcePrivateTransitSubnets_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_nat_private_transit_subnets" "test" {
  depends_on = [
    resource.huaweicloud_nat_private_transit_subnet.test,
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
  transit_subnet_subnet_id = data.huaweicloud_nat_private_transit_subnets.test.transit_subnets[0].virsubnet_id
}

data "huaweicloud_nat_private_transit_subnets" "filter_by_transit_subnet_subnet_id" {
  virsubnet_ids = [local.transit_subnet_subnet_id]
}

locals {
  transit_subnet_subnet_id_filter_result = [
    for v in data.huaweicloud_nat_private_transit_subnets.filter_by_transit_subnet_subnet_id.transit_subnets[*].virsubnet_id : 
    v == local.transit_subnet_subnet_id
  ]
}

output "transit_subnet_subnet_id_filter_is_useful" {
  value = alltrue(local.transit_subnet_subnet_id_filter_result) && length(local.transit_subnet_subnet_id_filter_result) > 0
}

locals {
  transit_subnet_project_id = data.huaweicloud_nat_private_transit_subnets.test.transit_subnets[0].virsubnet_project_id
}

data "huaweicloud_nat_private_transit_subnets" "filter_by_transit_subnet_project_id" {
  virsubnet_project_ids = [local.transit_subnet_project_id]
}

locals {
  transit_subnet_project_id_filter_result = [
    for v in data.huaweicloud_nat_private_transit_subnets.filter_by_transit_subnet_project_id.transit_subnets[*].project_id : 
    v == local.transit_subnet_project_id
  ]
}

output "transit_subnet_project_id_filter_is_useful" {
  value = alltrue(local.transit_subnet_project_id_filter_result) && length(local.transit_subnet_project_id_filter_result) > 0
}
`, testDataSourceNatPrivateTransitSubnets_base(name))
}

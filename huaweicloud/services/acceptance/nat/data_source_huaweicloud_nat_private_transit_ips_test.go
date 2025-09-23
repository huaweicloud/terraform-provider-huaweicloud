package nat

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccDatasourcePrivateTransitIps_basic(t *testing.T) {
	var (
		name           = acceptance.RandomAccResourceName()
		dataSourceName = "data.huaweicloud_nat_private_transit_ips.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)

		byTransitIpId   = "data.huaweicloud_nat_private_transit_ips.filter_by_transit_ip_id"
		dcByTransitIpId = acceptance.InitDataSourceCheck(byTransitIpId)

		byIpAddress   = "data.huaweicloud_nat_private_transit_ips.filter_by_ip_address"
		dcByIpAddress = acceptance.InitDataSourceCheck(byIpAddress)

		bySubnetId   = "data.huaweicloud_nat_private_transit_ips.filter_by_subnet_id"
		dcBySubnetId = acceptance.InitDataSourceCheck(bySubnetId)

		byTags   = "data.huaweicloud_nat_private_transit_ips.filter_by_tags"
		dcByTags = acceptance.InitDataSourceCheck(byTags)

		byNetWorkInterfaceId   = "data.huaweicloud_nat_private_transit_ips.filter_by_network_interface_id"
		dcByNetWorkInterfaceId = acceptance.InitDataSourceCheck(byNetWorkInterfaceId)

		byEps   = "data.huaweicloud_nat_private_transit_ips.filter_by_eps"
		dcByEps = acceptance.InitDataSourceCheck(byEps)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourcePrivateTransitIps_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					dcByTransitIpId.CheckResourceExists(),
					resource.TestCheckOutput("transit_ip_id_filter_is_useful", "true"),

					dcByIpAddress.CheckResourceExists(),
					resource.TestCheckOutput("ip_address_filter_is_useful", "true"),

					dcBySubnetId.CheckResourceExists(),
					resource.TestCheckOutput("subnet_id_filter_is_useful", "true"),

					dcByTags.CheckResourceExists(),
					resource.TestCheckOutput("tags_filter_is_useful", "true"),

					dcByNetWorkInterfaceId.CheckResourceExists(),
					resource.TestCheckOutput("network_interface_id_filter_is_useful", "true"),

					dcByEps.CheckResourceExists(),
					resource.TestCheckOutput("eps_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccPrivateTransitIpsDataSource_base(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_nat_private_transit_ip" "test" {
  subnet_id             = huaweicloud_vpc_subnet.test.id
  enterprise_project_id = "0"

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, common.TestBaseNetwork(name))
}

func testAccDatasourcePrivateTransitIps_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_nat_private_transit_ips" "test" {
  depends_on = [
    huaweicloud_nat_private_transit_ip.test
  ]  
}

locals {
  transit_ip_id = data.huaweicloud_nat_private_transit_ips.test.transit_ips[0].id
}

data "huaweicloud_nat_private_transit_ips" "filter_by_transit_ip_id" {
  transit_ip_id = local.transit_ip_id
}

locals {
  transit_ip_id_filter_result = [
    for v in data.huaweicloud_nat_private_transit_ips.filter_by_transit_ip_id.transit_ips[*].id : 
    v == local.transit_ip_id
  ]
}

output "transit_ip_id_filter_is_useful" {
  value = alltrue(local.transit_ip_id_filter_result) && length(local.transit_ip_id_filter_result) > 0
}

locals {
  ip_address = data.huaweicloud_nat_private_transit_ips.test.transit_ips[0].ip_address
}

data "huaweicloud_nat_private_transit_ips" "filter_by_ip_address" {
  ip_address = local.ip_address
}

locals {
  ip_address_filter_result = [
    for v in data.huaweicloud_nat_private_transit_ips.filter_by_ip_address.transit_ips[*].ip_address : 
    v == local.ip_address
  ]
}

output "ip_address_filter_is_useful" {
  value = alltrue(local.ip_address_filter_result) && length(local.ip_address_filter_result) > 0
}

locals {
  subnet_id = data.huaweicloud_nat_private_transit_ips.test.transit_ips[0].subnet_id
}

data "huaweicloud_nat_private_transit_ips" "filter_by_subnet_id" {
  subnet_id = local.subnet_id
}

locals {
  subnet_id_filter_result = [
    for v in data.huaweicloud_nat_private_transit_ips.filter_by_subnet_id.transit_ips[*].subnet_id : 
    v == local.subnet_id
  ]
}

output "subnet_id_filter_is_useful" {
  value = alltrue(local.subnet_id_filter_result) && length(local.subnet_id_filter_result) > 0
}

locals {
  network_interface_id = data.huaweicloud_nat_private_transit_ips.test.transit_ips[0].network_interface_id
}

data "huaweicloud_nat_private_transit_ips" "filter_by_network_interface_id" {
  network_interface_id = local.network_interface_id
}

locals {
  network_interface_id_filter_result = [
    for v in data.huaweicloud_nat_private_transit_ips.filter_by_network_interface_id.transit_ips[*].
    network_interface_id : v == local.network_interface_id
  ]
}

output "network_interface_id_filter_is_useful" {
  value = alltrue(local.network_interface_id_filter_result) && length(local.network_interface_id_filter_result) > 0
}

locals {
  enterprise_project_id = data.huaweicloud_nat_private_transit_ips.test.transit_ips[0].enterprise_project_id
}

data "huaweicloud_nat_private_transit_ips" "filter_by_eps" {
  enterprise_project_id = local.enterprise_project_id
}

locals {
  eps_filter_result = [
    for v in data.huaweicloud_nat_private_transit_ips.filter_by_eps.transit_ips[*].enterprise_project_id : 
    v == local.enterprise_project_id
  ]
}

output "eps_filter_is_useful" {
  value = alltrue(local.eps_filter_result) && length(local.eps_filter_result) > 0
}

locals {
  tags = data.huaweicloud_nat_private_transit_ips.test.transit_ips[0].tags
}

data "huaweicloud_nat_private_transit_ips" "filter_by_tags" {
  tags = local.tags
}

locals {
  tags_filter_result = [
    for tagMap in data.huaweicloud_nat_private_transit_ips.filter_by_tags.transit_ips[*].tags : alltrue([
      for k, v in local.tags : tagMap[k] == v
    ])
  ]
}

output "tags_filter_is_useful" {
  value = alltrue(local.tags_filter_result) && length(local.tags_filter_result) > 0
}
`, testAccPrivateTransitIpsDataSource_base(name))
}

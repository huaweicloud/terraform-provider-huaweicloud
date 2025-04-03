package nat

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccDatasourcePrivateSnatRules_basic(t *testing.T) {
	var (
		name           = acceptance.RandomAccResourceName()
		baseConfig     = testAccPrivateSnatRulesDataSource_base(name)
		dataSourceName = "data.huaweicloud_nat_private_snat_rules.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)

		bySnatId   = "data.huaweicloud_nat_private_snat_rules.filter_by_rule_id"
		dcBySnatId = acceptance.InitDataSourceCheck(bySnatId)

		byGatewayId   = "data.huaweicloud_nat_private_snat_rules.filter_by_gateway_id"
		dcByGatewayId = acceptance.InitDataSourceCheck(byGatewayId)

		byCidr   = "data.huaweicloud_nat_private_snat_rules.filter_by_cidr"
		dcByCidr = acceptance.InitDataSourceCheck(byCidr)

		bySubnetId   = "data.huaweicloud_nat_private_snat_rules.filter_by_subnet_id"
		dcBySubnetId = acceptance.InitDataSourceCheck(bySubnetId)

		byTransitIpId   = "data.huaweicloud_nat_private_snat_rules.filter_by_transit_ip_id"
		dcByTransitIpId = acceptance.InitDataSourceCheck(byTransitIpId)

		byTransitIpAddress   = "data.huaweicloud_nat_private_snat_rules.filter_by_transit_ip_address"
		dcByTransitIpAddress = acceptance.InitDataSourceCheck(byTransitIpAddress)

		byEps   = "data.huaweicloud_nat_private_snat_rules.filter_by_eps"
		dcByEps = acceptance.InitDataSourceCheck(byEps)

		byDescription   = "data.huaweicloud_nat_private_snat_rules.filter_by_description"
		dcByDescription = acceptance.InitDataSourceCheck(byDescription)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourcePrivateSnatRules_basic(baseConfig),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "rules.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "rules.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "rules.0.gateway_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "rules.0.subnet_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "rules.0.status"),
					resource.TestCheckResourceAttrSet(dataSourceName, "rules.0.transit_ip_associations.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "rules.0.transit_ip_associations.0.transit_ip_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "rules.0.transit_ip_associations.0.transit_ip_address"),

					dcBySnatId.CheckResourceExists(),
					resource.TestCheckOutput("rule_id_filter_is_useful", "true"),

					dcByGatewayId.CheckResourceExists(),
					resource.TestCheckOutput("gateway_id_filter_is_useful", "true"),

					dcByCidr.CheckResourceExists(),
					resource.TestCheckOutput("cidr_filter_is_useful", "true"),

					dcBySubnetId.CheckResourceExists(),
					resource.TestCheckOutput("subnet_id_filter_is_useful", "true"),

					dcByTransitIpId.CheckResourceExists(),
					resource.TestCheckOutput("transit_ip_id_filter_is_useful", "true"),

					dcByTransitIpAddress.CheckResourceExists(),
					resource.TestCheckOutput("transit_ip_address_filter_is_useful", "true"),

					dcByEps.CheckResourceExists(),
					resource.TestCheckOutput("eps_filter_is_useful", "true"),

					dcByDescription.CheckResourceExists(),
					resource.TestCheckOutput("description_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccPrivateSnatRulesDataSource_base(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_vpc" "transit_ip_used" {
  name = "%[2]s-transit"
  cidr = "192.168.0.0/24"
}

resource "huaweicloud_vpc_subnet" "transit_ip_used" {
  vpc_id     = huaweicloud_vpc.transit_ip_used.id
  name       = "%[2]s-transit"
  cidr       = "192.168.0.0/24"
  gateway_ip = "192.168.0.1"
}

resource "huaweicloud_nat_private_transit_ip" "test" {
  subnet_id             = huaweicloud_vpc_subnet.transit_ip_used.id
  enterprise_project_id = "0"
}

resource "huaweicloud_nat_private_gateway" "test" {
  subnet_id             = huaweicloud_vpc_subnet.test.id
  name                  = "%[2]s"
  enterprise_project_id = "0"
}

resource "huaweicloud_nat_private_snat_rule" "test" {
  gateway_id    = huaweicloud_nat_private_gateway.test.id
  description   = "Created by acc test"
  transit_ip_id = huaweicloud_nat_private_transit_ip.test.id
  subnet_id     = huaweicloud_vpc_subnet.test.id
}
`, common.TestBaseNetwork(name), name)
}

func testAccDatasourcePrivateSnatRules_basic(baseConfig string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_nat_private_snat_rules" "test" {
  depends_on = [
    huaweicloud_nat_private_snat_rule.test
  ]
}


locals {
  rule_id = data.huaweicloud_nat_private_snat_rules.test.rules[0].id
}

data "huaweicloud_nat_private_snat_rules" "filter_by_rule_id" {
  rule_id = local.rule_id
}

locals {
  rule_id_filter_result = [
    for v in data.huaweicloud_nat_private_snat_rules.filter_by_rule_id.rules[*].id : v == local.rule_id
  ]
}

output "rule_id_filter_is_useful" {
  value = alltrue(local.rule_id_filter_result) && length(local.rule_id_filter_result) > 0
}

locals {
  gateway_id = data.huaweicloud_nat_private_snat_rules.test.rules[0].gateway_id
}

data "huaweicloud_nat_private_snat_rules" "filter_by_gateway_id" {
  gateway_id = local.gateway_id
}

locals {
  gateway_id_filter_result = [
    for v in data.huaweicloud_nat_private_snat_rules.filter_by_gateway_id.rules[*].gateway_id : 
    v == local.gateway_id
  ]
}

output "gateway_id_filter_is_useful" {
  value = alltrue(local.gateway_id_filter_result) && length(local.gateway_id_filter_result) > 0
}

locals {
  cidr = data.huaweicloud_nat_private_snat_rules.test.rules[0].cidr
}

data "huaweicloud_nat_private_snat_rules" "filter_by_cidr" {
  cidr = local.cidr
}

locals {
  cidr_filter_result = [
    for v in data.huaweicloud_nat_private_snat_rules.filter_by_cidr.rules[*].cidr : v == local.cidr
  ]
}

output "cidr_filter_is_useful" {
  value = alltrue(local.cidr_filter_result) && length(local.cidr_filter_result) > 0
}

locals {
  transit_ip_id = data.huaweicloud_nat_private_snat_rules.test.rules[0].transit_ip_id
}

data "huaweicloud_nat_private_snat_rules" "filter_by_transit_ip_id" {
  transit_ip_id = local.transit_ip_id
}

locals {
  transit_ip_id_filter_result = [
    for v in data.huaweicloud_nat_private_snat_rules.filter_by_transit_ip_id.rules[*].transit_ip_id : 
    v == local.transit_ip_id
  ]
}

output "transit_ip_id_filter_is_useful" {
  value = alltrue(local.transit_ip_id_filter_result) && length(local.transit_ip_id_filter_result) > 0
}

locals {
  transit_ip_address = data.huaweicloud_nat_private_snat_rules.test.rules[0].transit_ip_address
}

data "huaweicloud_nat_private_snat_rules" "filter_by_transit_ip_address" {
  transit_ip_address = local.transit_ip_address
}

locals {
  transit_ip_address_filter_result = [
    for v in data.huaweicloud_nat_private_snat_rules.filter_by_transit_ip_address.rules[*].transit_ip_address : 
    v == local.transit_ip_address
  ]
}

output "transit_ip_address_filter_is_useful" {
  value = alltrue(local.transit_ip_address_filter_result) && length(local.transit_ip_address_filter_result) > 0
}

locals {
  enterprise_project_id = data.huaweicloud_nat_private_snat_rules.test.rules[0].enterprise_project_id
}

data "huaweicloud_nat_private_snat_rules" "filter_by_eps" {
  enterprise_project_id = local.enterprise_project_id
}

locals {
  eps_filter_result = [
    for v in data.huaweicloud_nat_private_snat_rules.filter_by_eps.rules[*].enterprise_project_id : 
    v == local.enterprise_project_id
  ]
}

output "eps_filter_is_useful" {
  value = alltrue(local.eps_filter_result) && length(local.eps_filter_result) > 0
}

locals {
  subnet_id = [
    for v in data.huaweicloud_nat_private_snat_rules.test.rules[*].subnet_id : v  if v != ""
  ]
}

data "huaweicloud_nat_private_snat_rules" "filter_by_subnet_id" {
  subnet_id = local.subnet_id[0]
}

locals {
  subnet_id_filter_result = [
    for v in data.huaweicloud_nat_private_snat_rules.filter_by_subnet_id.rules[*].subnet_id : 
    v == local.subnet_id[0]
  ]
}

output "subnet_id_filter_is_useful" {
  value = alltrue(local.subnet_id_filter_result) && length(local.subnet_id_filter_result) > 0
}

locals {
  description = data.huaweicloud_nat_private_snat_rules.test.rules[0].description
}

data "huaweicloud_nat_private_snat_rules" "filter_by_description" {
  description = [local.description]
}

locals {
  description_filter_result = [
    for v in data.huaweicloud_nat_private_snat_rules.filter_by_description.rules[*].description : 
    v == local.description
  ]
}

output "description_filter_is_useful" {
  value = alltrue(local.description_filter_result) && length(local.description_filter_result) > 0
}
`, baseConfig)
}

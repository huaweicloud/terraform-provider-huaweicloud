package nat

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccDatasourceSnatRules_basic(t *testing.T) {
	var (
		name           = acceptance.RandomAccResourceName()
		baseConfig     = testAccSnatRulesDataSource_base(name)
		dataSourceName = "data.huaweicloud_nat_snat_rules.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)

		byRuleId   = "data.huaweicloud_nat_snat_rules.filter_by_rule_id"
		dcByRuleId = acceptance.InitDataSourceCheck(byRuleId)

		byGatewayId   = "data.huaweicloud_nat_snat_rules.filter_by_gateway_id"
		dcByGatewayId = acceptance.InitDataSourceCheck(byGatewayId)

		bySourceType   = "data.huaweicloud_nat_snat_rules.filter_by_source_type"
		dcBySourceType = acceptance.InitDataSourceCheck(bySourceType)

		bySubnetId   = "data.huaweicloud_nat_snat_rules.filter_by_subnet_id"
		dcBySubnetId = acceptance.InitDataSourceCheck(bySubnetId)

		byEipId   = "data.huaweicloud_nat_snat_rules.filter_by_floating_ip_id"
		dcByEipId = acceptance.InitDataSourceCheck(byEipId)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceSnatRules_basic(baseConfig),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					dcByRuleId.CheckResourceExists(),
					resource.TestCheckOutput("rule_id_filter_is_useful", "true"),

					dcByGatewayId.CheckResourceExists(),
					resource.TestCheckOutput("gateway_id_filter_is_useful", "true"),

					dcBySourceType.CheckResourceExists(),
					resource.TestCheckOutput("source_type_filter_is_useful", "true"),

					dcBySubnetId.CheckResourceExists(),
					resource.TestCheckOutput("subnet_id_filter_is_useful", "true"),

					dcByEipId.CheckResourceExists(),
					resource.TestCheckOutput("floating_ip_id_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccSnatRulesDataSource_base(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_nat_gateway" "test" {
  name                  = "%[2]s"
  description           = "created by terraform"
  spec                  = "1"
  vpc_id                = huaweicloud_vpc.test.id
  subnet_id             = huaweicloud_vpc_subnet.test.id
  enterprise_project_id = "0"
}

resource "huaweicloud_vpc_eip" "test" {
  publicip {
    type = "5_bgp"
  }

  bandwidth {
    name        = "%[2]s"
    size        = 5
    share_type  = "PER"
    charge_mode = "traffic"
  }
}

resource "huaweicloud_nat_snat_rule" "test" {
  nat_gateway_id = huaweicloud_nat_gateway.test.id
  floating_ip_id = huaweicloud_vpc_eip.test.id
  subnet_id      = huaweicloud_vpc_subnet.test.id
}
`, common.TestBaseNetwork(name), name)
}

func testAccDatasourceSnatRules_basic(baseConfig string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_nat_snat_rules" "test" {
  depends_on = [
    huaweicloud_nat_snat_rule.test
  ]
}

locals {
  rule_id = data.huaweicloud_nat_snat_rules.test.rules[0].id
}

data "huaweicloud_nat_snat_rules" "filter_by_rule_id" {
  rule_id = local.rule_id
}

locals {
  rule_id_filter_result = [
    for v in data.huaweicloud_nat_snat_rules.filter_by_rule_id.rules[*].id : v == local.rule_id
  ]
}

output "rule_id_filter_is_useful" {
  value = alltrue(local.rule_id_filter_result) && length(local.rule_id_filter_result) > 0
}

locals {
  gateway_id = data.huaweicloud_nat_snat_rules.test.rules[0].gateway_id
}

data "huaweicloud_nat_snat_rules" "filter_by_gateway_id" {
  gateway_id = local.gateway_id
}

locals {
  gateway_id_filter_result = [
    for v in data.huaweicloud_nat_snat_rules.filter_by_gateway_id.rules[*].gateway_id : v == local.gateway_id
  ]
}

output "gateway_id_filter_is_useful" {
  value = alltrue(local.gateway_id_filter_result) && length(local.gateway_id_filter_result) > 0
}

locals {
  source_type = data.huaweicloud_nat_snat_rules.test.rules[0].source_type
}

data "huaweicloud_nat_snat_rules" "filter_by_source_type" {
  source_type = local.source_type
}

locals {
  source_type_filter_result = [
    for v in data.huaweicloud_nat_snat_rules.filter_by_source_type.rules[*].source_type : v == local.source_type
  ]
}

output "source_type_filter_is_useful" {
  value = alltrue(local.source_type_filter_result) && length(local.source_type_filter_result) > 0
}

locals {
  subnet_id = data.huaweicloud_nat_snat_rules.test.rules[0].subnet_id
}

data "huaweicloud_nat_snat_rules" "filter_by_subnet_id" {
  subnet_id = local.subnet_id
}

locals {
  subnet_id_filter_result = [
    for v in data.huaweicloud_nat_snat_rules.filter_by_subnet_id.rules[*].subnet_id : v == local.subnet_id
  ]
}

output "subnet_id_filter_is_useful" {
  value = alltrue(local.subnet_id_filter_result) && length(local.subnet_id_filter_result) > 0
}

locals {
  floating_ip_id = data.huaweicloud_nat_snat_rules.test.rules[0].floating_ip_id
}

data "huaweicloud_nat_snat_rules" "filter_by_floating_ip_id" {
  floating_ip_id = local.floating_ip_id
}

locals {
  floating_ip_id_filter_result = [
    for v in data.huaweicloud_nat_snat_rules.filter_by_floating_ip_id.rules[*].floating_ip_id : 
    v == local.floating_ip_id
  ]
}

output "floating_ip_id_filter_is_useful" {
  value = alltrue(local.floating_ip_id_filter_result) && length(local.floating_ip_id_filter_result) > 0
}
`, baseConfig)
}

func TestAccDatasourceSnatRules_direct(t *testing.T) {
	var (
		name           = acceptance.RandomAccResourceName()
		baseConfig     = testAccSnatRulesDataSource_direct_base(name)
		dataSourceName = "data.huaweicloud_nat_snat_rules.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)

		byCidr   = "data.huaweicloud_nat_snat_rules.filter_by_cidr"
		dcByCidr = acceptance.InitDataSourceCheck(byCidr)

		bySourceType   = "data.huaweicloud_nat_snat_rules.filter_by_source_type"
		dcBySourceType = acceptance.InitDataSourceCheck(bySourceType)

		byStatus   = "data.huaweicloud_nat_snat_rules.filter_by_status"
		dcByStatus = acceptance.InitDataSourceCheck(byStatus)

		byEipAddress   = "data.huaweicloud_nat_snat_rules.filter_by_floating_ip_address"
		dcByEipAddress = acceptance.InitDataSourceCheck(byEipAddress)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceSnatRules_direct(baseConfig),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					dcByCidr.CheckResourceExists(),
					resource.TestCheckOutput("cidr_filter_is_useful", "true"),

					dcBySourceType.CheckResourceExists(),
					resource.TestCheckOutput("source_type_filter_is_useful", "true"),

					dcByStatus.CheckResourceExists(),
					resource.TestCheckOutput("status_filter_is_useful", "true"),

					dcByEipAddress.CheckResourceExists(),
					resource.TestCheckOutput("floating_ip_address_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccSnatRulesDataSource_direct_base(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_nat_gateway" "test" {
  name                  = "%[2]s"
  description           = "created by terraform"
  spec                  = "1"
  vpc_id                = huaweicloud_vpc.test.id
  subnet_id             = huaweicloud_vpc_subnet.test.id
  enterprise_project_id = "0"
}

resource "huaweicloud_vpc_eip" "test" {
  publicip {
    type = "5_bgp"
  }

  bandwidth {
    name        = "%[2]s"
    size        = 5
    share_type  = "PER"
    charge_mode = "traffic"
  }
}

resource "huaweicloud_nat_snat_rule" "direct" {
  nat_gateway_id = huaweicloud_nat_gateway.test.id
  floating_ip_id = huaweicloud_vpc_eip.test.id
  source_type    = 1
  cidr           = "192.168.1.0/24"  
}
`, common.TestBaseNetwork(name), name)
}

func testAccDatasourceSnatRules_direct(baseConfig string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_nat_snat_rules" "test" {
  depends_on = [
    huaweicloud_nat_snat_rule.direct
  ]
}

locals {
  cidr = data.huaweicloud_nat_snat_rules.test.rules[0].cidr
}

data "huaweicloud_nat_snat_rules" "filter_by_cidr" {
  cidr = local.cidr
}

locals {
  cidr_filter_result = [
    for v in data.huaweicloud_nat_snat_rules.filter_by_cidr.rules[*].cidr : 
    v == local.cidr
  ]
}

output "cidr_filter_is_useful" {
  value = alltrue(local.cidr_filter_result) && length(local.cidr_filter_result) > 0
}

locals {
  source_type = data.huaweicloud_nat_snat_rules.test.rules[0].source_type
}

data "huaweicloud_nat_snat_rules" "filter_by_source_type" {
  source_type = local.source_type
}

locals {
  source_type_filter_result = [
    for v in data.huaweicloud_nat_snat_rules.filter_by_source_type.rules[*].source_type : 
    v == local.source_type
  ]
}

output "source_type_filter_is_useful" {
  value = alltrue(local.source_type_filter_result) && length(local.source_type_filter_result) > 0
}

locals {
  status = data.huaweicloud_nat_snat_rules.test.rules[0].status
}

data "huaweicloud_nat_snat_rules" "filter_by_status" {
  status = local.status
}

locals {
  status_filter_result = [
    for v in data.huaweicloud_nat_snat_rules.filter_by_status.rules[*].status : 
    v == local.status
  ]
}

output "status_filter_is_useful" {
  value = alltrue(local.status_filter_result) && length(local.status_filter_result) > 0
}

locals {
  floating_ip_address = data.huaweicloud_nat_snat_rules.test.rules[0].floating_ip_address
}

data "huaweicloud_nat_snat_rules" "filter_by_floating_ip_address" {
  floating_ip_address = local.floating_ip_address
}

locals {
  floating_ip_address_filter_result = [
    for v in data.huaweicloud_nat_snat_rules.filter_by_floating_ip_address.rules[*].floating_ip_address : 
    v == local.floating_ip_address
  ]
}

output "floating_ip_address_filter_is_useful" {
  value = alltrue(local.floating_ip_address_filter_result) && length(local.floating_ip_address_filter_result) > 0
}
`, baseConfig)
}

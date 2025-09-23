package nat

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccDatasourceDnatRules_basic(t *testing.T) {
	var (
		name           = acceptance.RandomAccResourceName()
		baseConfig     = testAccDnatRulesDataSource_base(name)
		dataSourceName = "data.huaweicloud_nat_dnat_rules.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)

		byRuleId   = "data.huaweicloud_nat_dnat_rules.filter_by_rule_id"
		dcByRuleId = acceptance.InitDataSourceCheck(byRuleId)

		byGatewayId   = "data.huaweicloud_nat_dnat_rules.filter_by_gateway_id"
		dcByGatewayId = acceptance.InitDataSourceCheck(byGatewayId)

		byProtocol   = "data.huaweicloud_nat_dnat_rules.filter_by_protocol"
		dcByProtocol = acceptance.InitDataSourceCheck(byProtocol)

		byInternalServicePort   = "data.huaweicloud_nat_dnat_rules.filter_by_internal_service_port"
		dcByInternalServicePort = acceptance.InitDataSourceCheck(byInternalServicePort)

		byPortId   = "data.huaweicloud_nat_dnat_rules.filter_by_port_id"
		dcByPortId = acceptance.InitDataSourceCheck(byPortId)

		byPrivateIp   = "data.huaweicloud_nat_dnat_rules.filter_by_private_ip"
		dcByPrivateIp = acceptance.InitDataSourceCheck(byPrivateIp)

		byStatus   = "data.huaweicloud_nat_dnat_rules.filter_by_status"
		dcByStatus = acceptance.InitDataSourceCheck(byStatus)

		byFloatingIpAddress   = "data.huaweicloud_nat_dnat_rules.filter_by_floating_ip_address"
		dcByFloatingIpAddress = acceptance.InitDataSourceCheck(byFloatingIpAddress)

		byDescription   = "data.huaweicloud_nat_dnat_rules.filter_by_description"
		dcByDescription = acceptance.InitDataSourceCheck(byDescription)

		byCreatedAt   = "data.huaweicloud_nat_dnat_rules.filter_by_created_at"
		dcByCreatedAt = acceptance.InitDataSourceCheck(byCreatedAt)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceDnatRules_basic(baseConfig),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "rules.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "rules.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "rules.0.protocol"),
					resource.TestCheckResourceAttrSet(dataSourceName, "rules.0.description"),
					resource.TestCheckResourceAttrSet(dataSourceName, "rules.0.status"),
					resource.TestCheckResourceAttrSet(dataSourceName, "rules.0.floating_ip_address"),
					resource.TestCheckResourceAttrSet(dataSourceName, "rules.0.internal_service_port"),
					resource.TestCheckResourceAttrSet(dataSourceName, "rules.0.created_at"),

					dcByRuleId.CheckResourceExists(),
					resource.TestCheckOutput("rule_id_filter_is_useful", "true"),

					dcByGatewayId.CheckResourceExists(),
					resource.TestCheckOutput("gateway_id_filter_is_useful", "true"),

					dcByProtocol.CheckResourceExists(),
					resource.TestCheckOutput("protocol_filter_is_useful", "true"),

					dcByInternalServicePort.CheckResourceExists(),
					resource.TestCheckOutput("internal_service_port_filter_is_useful", "true"),

					dcByPortId.CheckResourceExists(),
					resource.TestCheckOutput("port_id_filter_is_useful", "true"),

					dcByPrivateIp.CheckResourceExists(),
					resource.TestCheckOutput("private_ip_filter_is_useful", "true"),

					dcByStatus.CheckResourceExists(),
					resource.TestCheckOutput("status_filter_is_useful", "true"),

					dcByFloatingIpAddress.CheckResourceExists(),
					resource.TestCheckOutput("floating_ip_address_filter_is_useful", "true"),

					dcByDescription.CheckResourceExists(),
					resource.TestCheckOutput("description_filter_is_useful", "true"),

					dcByCreatedAt.CheckResourceExists(),
					resource.TestCheckOutput("created_at_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDnatRulesDataSource_base(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_vpc_eip" "test" {
  publicip {
    type = "5_bgp"
  }

  bandwidth {
    name        = "eip-test"
    size        = 5
    share_type  = "PER"
    charge_mode = "traffic"
  }
}

resource "huaweicloud_compute_instance" "test" {
  name              = "%[2]s"
  flavor_id         = data.huaweicloud_compute_flavors.test.ids[0]
  image_id          = data.huaweicloud_images_image.test.id
  security_groups   = [huaweicloud_networking_secgroup.test.name]
  availability_zone = data.huaweicloud_availability_zones.test.names[0]

  network {
    uuid = huaweicloud_vpc_subnet.test.id
  }
}

resource "huaweicloud_nat_gateway" "test" {
  name                  = "%[2]s"
  spec                  = "1"
  vpc_id                = huaweicloud_vpc.test.id
  subnet_id             = huaweicloud_vpc_subnet.test.id
  enterprise_project_id = "0"
}

resource "huaweicloud_nat_dnat_rule" "test" {
  nat_gateway_id        = huaweicloud_nat_gateway.test.id
  floating_ip_id        = huaweicloud_vpc_eip.test.id
  private_ip            = huaweicloud_compute_instance.test.network[0].fixed_ip_v4
  description           = "Created by acc test"
  protocol              = "tcp"
  internal_service_port = 60
  external_service_port = 2000
}
`, common.TestBaseComputeResources(name), name)
}

func testAccDatasourceDnatRules_basic(baseConfig string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_nat_dnat_rules" "test" {
  depends_on = [
    huaweicloud_nat_dnat_rule.test
  ]
}

locals {
  rule_id = data.huaweicloud_nat_dnat_rules.test.rules[0].id
}

data "huaweicloud_nat_dnat_rules" "filter_by_rule_id" {
  rule_id = local.rule_id
}

locals {
  rule_id_filter_result = [
    for v in data.huaweicloud_nat_dnat_rules.filter_by_rule_id.rules[*].id : v == local.rule_id
  ]
}

output "rule_id_filter_is_useful" {
  value = alltrue(local.rule_id_filter_result) && length(local.rule_id_filter_result) > 0
}

locals {
  gateway_id = data.huaweicloud_nat_dnat_rules.test.rules[0].gateway_id
}

data "huaweicloud_nat_dnat_rules" "filter_by_gateway_id" {
  gateway_id = local.gateway_id
}

locals {
  gateway_id_filter_result = [
    for v in data.huaweicloud_nat_dnat_rules.filter_by_gateway_id.rules[*].gateway_id : 
    v == local.gateway_id
  ]
}

output "gateway_id_filter_is_useful" {
  value = alltrue(local.gateway_id_filter_result) && length(local.gateway_id_filter_result) > 0
}

locals {
  protocol = data.huaweicloud_nat_dnat_rules.test.rules[0].protocol
}

data "huaweicloud_nat_dnat_rules" "filter_by_protocol" {
  protocol = local.protocol
}

locals {
  protocol_filter_result = [
    for v in data.huaweicloud_nat_dnat_rules.filter_by_protocol.rules[*].protocol : v == local.protocol
  ]
}

output "protocol_filter_is_useful" {
  value = alltrue(local.protocol_filter_result) && length(local.protocol_filter_result) > 0
}

locals {
  internal_service_port = data.huaweicloud_nat_dnat_rules.test.rules[0].internal_service_port
}

data "huaweicloud_nat_dnat_rules" "filter_by_internal_service_port" {
  internal_service_port = local.internal_service_port
}

locals {
  internal_service_port_filter_result = [
    for v in data.huaweicloud_nat_dnat_rules.filter_by_internal_service_port.rules[*].internal_service_port : 
    v == local.internal_service_port
  ]
}

output "internal_service_port_filter_is_useful" {
  value = alltrue(local.internal_service_port_filter_result) && length(local.internal_service_port_filter_result) > 0
}

locals {
  port_id = data.huaweicloud_nat_dnat_rules.test.rules[0].port_id
}

data "huaweicloud_nat_dnat_rules" "filter_by_port_id" {
  port_id = local.port_id
}

locals {
  port_id_filter_result = [
    for v in data.huaweicloud_nat_dnat_rules.filter_by_port_id.rules[*].port_id : v == local.port_id
  ]
}

output "port_id_filter_is_useful" {
  value = alltrue(local.port_id_filter_result) && length(local.port_id_filter_result) > 0
}

locals {
  private_ip = data.huaweicloud_nat_dnat_rules.test.rules[0].private_ip
}

data "huaweicloud_nat_dnat_rules" "filter_by_private_ip" {
  private_ip = local.private_ip
}

locals {
  private_ip_filter_result = [
    for v in data.huaweicloud_nat_dnat_rules.filter_by_private_ip.rules[*].private_ip : v == local.private_ip
  ]
}

output "private_ip_filter_is_useful" {
  value = alltrue(local.private_ip_filter_result) && length(local.private_ip_filter_result) > 0
}

locals {
  status = data.huaweicloud_nat_dnat_rules.test.rules[0].status
}

data "huaweicloud_nat_dnat_rules" "filter_by_status" {
  status = local.status
}

locals {
  status_filter_result = [
    for v in data.huaweicloud_nat_dnat_rules.filter_by_status.rules[*].status : v == local.status
  ]
}

output "status_filter_is_useful" {
  value = alltrue(local.status_filter_result) && length(local.status_filter_result) > 0
}


locals {
  floating_ip_address = data.huaweicloud_nat_dnat_rules.test.rules[0].floating_ip_address
}

data "huaweicloud_nat_dnat_rules" "filter_by_floating_ip_address" {
  floating_ip_address = local.floating_ip_address
}

locals {
  floating_ip_address_filter_result = [
    for v in data.huaweicloud_nat_dnat_rules.filter_by_floating_ip_address.rules[*].floating_ip_address : 
    v == local.floating_ip_address
  ]
}

output "floating_ip_address_filter_is_useful" {
  value = alltrue(local.floating_ip_address_filter_result) && length(local.floating_ip_address_filter_result) > 0
}

locals {
  description = data.huaweicloud_nat_dnat_rules.test.rules[0].description
}

data "huaweicloud_nat_dnat_rules" "filter_by_description" {
  description = local.description
}

locals {
  description_filter_result = [
    for v in data.huaweicloud_nat_dnat_rules.filter_by_description.rules[*].description : v == local.description
  ]
}

output "description_filter_is_useful" {
  value = alltrue(local.description_filter_result) && length(local.description_filter_result) > 0
}

locals {
  created_at = data.huaweicloud_nat_dnat_rules.test.rules[0].created_at
}

data "huaweicloud_nat_dnat_rules" "filter_by_created_at" {
  created_at = local.created_at
}

locals {
  created_at_filter_result = [
    for v in data.huaweicloud_nat_dnat_rules.filter_by_created_at.rules[*].created_at : v == local.created_at
  ]
}

output "created_at_filter_is_useful" {
  value = alltrue(local.created_at_filter_result) && length(local.created_at_filter_result) > 0
}
`, baseConfig)
}

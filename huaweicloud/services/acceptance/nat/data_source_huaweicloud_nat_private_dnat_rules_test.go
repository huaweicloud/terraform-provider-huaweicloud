package nat

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccDatasourcePrivateDnatRules_basic(t *testing.T) {
	var (
		name           = acceptance.RandomAccResourceName()
		baseConfig     = testAccPrivateDnatRulesDataSource_base(name)
		dataSourceName = "data.huaweicloud_nat_private_dnat_rules.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)

		byRuleId   = "data.huaweicloud_nat_private_dnat_rules.filter_by_rule_id"
		dcByRuleId = acceptance.InitDataSourceCheck(byRuleId)

		byGatewayId   = "data.huaweicloud_nat_private_dnat_rules.filter_by_gateway_id"
		dcByGatewayId = acceptance.InitDataSourceCheck(byGatewayId)

		byBackendType   = "data.huaweicloud_nat_private_dnat_rules.filter_by_backend_type"
		dcByBackendType = acceptance.InitDataSourceCheck(byBackendType)

		byProtocol   = "data.huaweicloud_nat_private_dnat_rules.filter_by_protocol"
		dcByProtocol = acceptance.InitDataSourceCheck(byProtocol)

		byInternalServicePort   = "data.huaweicloud_nat_private_dnat_rules.filter_by_internal_service_port"
		dcByInternalServicePort = acceptance.InitDataSourceCheck(byInternalServicePort)

		byBackendInterfaceId   = "data.huaweicloud_nat_private_dnat_rules.filter_by_backend_interface_id"
		dcByBackendInterfaceId = acceptance.InitDataSourceCheck(byBackendInterfaceId)

		byTransitIpId   = "data.huaweicloud_nat_private_dnat_rules.filter_by_transit_ip_id"
		dcByTransitIpId = acceptance.InitDataSourceCheck(byTransitIpId)

		byTransitServicePort   = "data.huaweicloud_nat_private_dnat_rules.filter_by_transit_service_port"
		dcByTransitServicePort = acceptance.InitDataSourceCheck(byTransitServicePort)

		byBackendPrivateIp   = "data.huaweicloud_nat_private_dnat_rules.filter_by_backend_private_ip"
		dcByBackendPrivateIp = acceptance.InitDataSourceCheck(byBackendPrivateIp)

		byEps   = "data.huaweicloud_nat_private_dnat_rules.filter_by_eps"
		dcByEps = acceptance.InitDataSourceCheck(byEps)

		byDescription   = "data.huaweicloud_nat_private_dnat_rules.filter_by_description"
		dcByDescription = acceptance.InitDataSourceCheck(byDescription)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourcePrivateDnatRules_basic(baseConfig),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					dcByRuleId.CheckResourceExists(),
					resource.TestCheckOutput("rule_id_filter_is_useful", "true"),

					dcByGatewayId.CheckResourceExists(),
					resource.TestCheckOutput("gateway_id_filter_is_useful", "true"),

					dcByBackendType.CheckResourceExists(),
					resource.TestCheckOutput("backend_type_filter_is_useful", "true"),

					dcByProtocol.CheckResourceExists(),
					resource.TestCheckOutput("protocol_filter_is_useful", "true"),

					dcByInternalServicePort.CheckResourceExists(),
					resource.TestCheckOutput("internal_service_port_filter_is_useful", "true"),

					dcByBackendInterfaceId.CheckResourceExists(),
					resource.TestCheckOutput("backend_interface_id_filter_is_useful", "true"),

					dcByTransitIpId.CheckResourceExists(),
					resource.TestCheckOutput("transit_ip_id_filter_is_useful", "true"),

					dcByTransitServicePort.CheckResourceExists(),
					resource.TestCheckOutput("transit_service_port_filter_is_useful", "true"),

					dcByBackendPrivateIp.CheckResourceExists(),
					resource.TestCheckOutput("backend_private_ip_filter_is_useful", "true"),

					dcByEps.CheckResourceExists(),
					resource.TestCheckOutput("eps_filter_is_useful", "true"),

					dcByDescription.CheckResourceExists(),
					resource.TestCheckOutput("description_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccPrivateDnatRulesDataSource_base(name string) string {
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

resource "huaweicloud_compute_instance" "test" {
  name              = "%[2]s-ecs"
  flavor_id         = data.huaweicloud_compute_flavors.test.ids[0]
  image_id          = data.huaweicloud_images_image.test.id
  security_groups   = [huaweicloud_networking_secgroup.test.name]
  availability_zone = data.huaweicloud_availability_zones.test.names[0]

  network {
    uuid = huaweicloud_vpc_subnet.test.id
  }
}

resource "huaweicloud_nat_private_gateway" "test" {
  subnet_id             = huaweicloud_vpc_subnet.test.id
  name                  = "%[2]s"
  enterprise_project_id = "0"
}

resource "huaweicloud_nat_private_dnat_rule" "test" {
  gateway_id            = huaweicloud_nat_private_gateway.test.id
  protocol              = "tcp"
  description           = "Created by terraform"
  transit_ip_id         = huaweicloud_nat_private_transit_ip.test.id
  transit_service_port  = 1000
  backend_interface_id  = huaweicloud_compute_instance.test.network[0].port
  internal_service_port = 2000
}
`, common.TestBaseComputeResources(name), name)
}

func testAccDatasourcePrivateDnatRules_basic(baseConfig string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_nat_private_dnat_rules" "test" {
  depends_on = [
    huaweicloud_nat_private_dnat_rule.test
  ]
}

locals {
  rule_id = data.huaweicloud_nat_private_dnat_rules.test.rules[0].id
}

data "huaweicloud_nat_private_dnat_rules" "filter_by_rule_id" {
  rule_id = local.rule_id
}

locals {
  rule_id_filter_result = [
    for v in data.huaweicloud_nat_private_dnat_rules.filter_by_rule_id.rules[*].id : v == local.rule_id
  ]
}

output "rule_id_filter_is_useful" {
  value = alltrue(local.rule_id_filter_result) && length(local.rule_id_filter_result) > 0
}

locals {
  gateway_id = data.huaweicloud_nat_private_dnat_rules.test.rules[0].gateway_id
}

data "huaweicloud_nat_private_dnat_rules" "filter_by_gateway_id" {
  gateway_id = local.gateway_id
}

locals {
  gateway_id_filter_result = [
    for v in data.huaweicloud_nat_private_dnat_rules.filter_by_gateway_id.rules[*].gateway_id : 
    v == local.gateway_id
  ]
}

output "gateway_id_filter_is_useful" {
  value = alltrue(local.gateway_id_filter_result) && length(local.gateway_id_filter_result) > 0
}

locals {
  backend_type = data.huaweicloud_nat_private_dnat_rules.test.rules[0].backend_type
}

data "huaweicloud_nat_private_dnat_rules" "filter_by_backend_type" {
  backend_type = local.backend_type
}

locals {
  backend_type_filter_result = [
    for v in data.huaweicloud_nat_private_dnat_rules.filter_by_backend_type.rules[*].backend_type : 
    v == local.backend_type
  ]
}

output "backend_type_filter_is_useful" {
  value = alltrue(local.backend_type_filter_result) && length(local.backend_type_filter_result) > 0
}

locals {
  protocol = data.huaweicloud_nat_private_dnat_rules.test.rules[0].protocol
}

data "huaweicloud_nat_private_dnat_rules" "filter_by_protocol" {
  protocol = local.protocol
}

locals {
  protocol_filter_result = [
    for v in data.huaweicloud_nat_private_dnat_rules.filter_by_protocol.rules[*].protocol : v == local.protocol
  ]
}

output "protocol_filter_is_useful" {
  value = alltrue(local.protocol_filter_result) && length(local.protocol_filter_result) > 0
}

locals {
  internal_service_port = data.huaweicloud_nat_private_dnat_rules.test.rules[0].internal_service_port
}

data "huaweicloud_nat_private_dnat_rules" "filter_by_internal_service_port" {
  internal_service_port = local.internal_service_port
}

locals {
  internal_service_port_filter_result = [
    for v in data.huaweicloud_nat_private_dnat_rules.filter_by_internal_service_port.rules[*].internal_service_port : 
    v == local.internal_service_port
  ]
}

output "internal_service_port_filter_is_useful" {
  value = alltrue(local.internal_service_port_filter_result) && length(local.internal_service_port_filter_result) > 0
}

locals {
  backend_interface_id = data.huaweicloud_nat_private_dnat_rules.test.rules[0].backend_interface_id
}

data "huaweicloud_nat_private_dnat_rules" "filter_by_backend_interface_id" {
  backend_interface_id = local.backend_interface_id
}

locals {
  backend_interface_id_filter_result = [
    for v in data.huaweicloud_nat_private_dnat_rules.filter_by_backend_interface_id.rules[*].backend_interface_id : 
    v == local.backend_interface_id
  ]
}

output "backend_interface_id_filter_is_useful" {
  value = alltrue(local.backend_interface_id_filter_result) && length(local.backend_interface_id_filter_result) > 0
}

locals {
  transit_ip_id = data.huaweicloud_nat_private_dnat_rules.test.rules[0].transit_ip_id
}

data "huaweicloud_nat_private_dnat_rules" "filter_by_transit_ip_id" {
  transit_ip_id = local.transit_ip_id
}

locals {
  transit_ip_id_filter_result = [
    for v in data.huaweicloud_nat_private_dnat_rules.filter_by_transit_ip_id.rules[*].transit_ip_id : 
    v == local.transit_ip_id
  ]
}

output "transit_ip_id_filter_is_useful" {
  value = alltrue(local.transit_ip_id_filter_result) && length(local.transit_ip_id_filter_result) > 0
}

locals {
  transit_service_port = data.huaweicloud_nat_private_dnat_rules.test.rules[0].transit_service_port
}

data "huaweicloud_nat_private_dnat_rules" "filter_by_transit_service_port" {
  transit_service_port = local.transit_service_port
}

locals {
  transit_service_port_filter_result = [
    for v in data.huaweicloud_nat_private_dnat_rules.filter_by_transit_service_port.rules[*].transit_service_port : 
    v == local.transit_service_port
  ]
}

output "transit_service_port_filter_is_useful" {
  value = alltrue(local.transit_service_port_filter_result) && length(local.transit_service_port_filter_result) > 0
}


locals {
  backend_private_ip = data.huaweicloud_nat_private_dnat_rules.test.rules[0].backend_private_ip
}

data "huaweicloud_nat_private_dnat_rules" "filter_by_backend_private_ip" {
  backend_private_ip = local.backend_private_ip
}

locals {
  backend_private_ip_filter_result = [
    for v in data.huaweicloud_nat_private_dnat_rules.filter_by_backend_private_ip.rules[*].backend_private_ip : 
    v == local.backend_private_ip
  ]
}

output "backend_private_ip_filter_is_useful" {
  value = alltrue(local.backend_private_ip_filter_result) && length(local.backend_private_ip_filter_result) > 0
}

locals {
  enterprise_project_id = data.huaweicloud_nat_private_dnat_rules.test.rules[0].enterprise_project_id
}

data "huaweicloud_nat_private_dnat_rules" "filter_by_eps" {
  enterprise_project_id = local.enterprise_project_id
}

locals {
  eps_filter_result = [
    for v in data.huaweicloud_nat_private_dnat_rules.filter_by_eps.rules[*].enterprise_project_id : 
    v == local.enterprise_project_id
  ]
}

output "eps_filter_is_useful" {
  value = alltrue(local.eps_filter_result) && length(local.eps_filter_result) > 0
}

locals {
  description = data.huaweicloud_nat_private_dnat_rules.test.rules[0].description
}

data "huaweicloud_nat_private_dnat_rules" "filter_by_description" {
  description = [local.description]
}

locals {
  description_filter_result = [
    for v in data.huaweicloud_nat_private_dnat_rules.filter_by_description.rules[*].description : 
    v == local.description
  ]
}

output "description_filter_is_useful" {
  value = alltrue(local.description_filter_result) && length(local.description_filter_result) > 0
}
`, baseConfig)
}

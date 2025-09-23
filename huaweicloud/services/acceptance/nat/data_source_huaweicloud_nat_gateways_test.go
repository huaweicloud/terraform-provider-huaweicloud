package nat

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccDatasourcePublicGateways_basic(t *testing.T) {
	var (
		name           = acceptance.RandomAccResourceName()
		baseConfig     = testAccPublicGatewaysDataSource_base()
		dataSourceName = "data.huaweicloud_nat_gateways.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)

		byName         = "data.huaweicloud_nat_gateways.filter_by_name"
		nameNotFound   = "data.huaweicloud_nat_gateways.not_found"
		dcByName       = acceptance.InitDataSourceCheck(byName)
		dcNameNotFound = acceptance.InitDataSourceCheck(nameNotFound)

		byGatewayId   = "data.huaweicloud_nat_gateways.filter_by_gateway_id"
		dcByGatewayId = acceptance.InitDataSourceCheck(byGatewayId)

		bySpec   = "data.huaweicloud_nat_gateways.filter_by_spec"
		dcBySpec = acceptance.InitDataSourceCheck(bySpec)

		byStatus   = "data.huaweicloud_nat_gateways.filter_by_status"
		dcByStatus = acceptance.InitDataSourceCheck(byStatus)

		byVpcId   = "data.huaweicloud_nat_gateways.filter_by_vpc_id"
		dcByVpcId = acceptance.InitDataSourceCheck(byVpcId)

		bySubnetId   = "data.huaweicloud_nat_gateways.filter_by_subnet_id"
		dcBySubnetId = acceptance.InitDataSourceCheck(bySubnetId)

		byEps   = "data.huaweicloud_nat_gateways.filter_by_eps"
		dcByEps = acceptance.InitDataSourceCheck(byEps)

		byDescription   = "data.huaweicloud_nat_gateways.filter_by_description"
		dcByDescription = acceptance.InitDataSourceCheck(byDescription)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourcePublicGateways_basic(baseConfig, name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "gateways.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "gateways.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "gateways.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "gateways.0.status"),
					resource.TestCheckResourceAttrSet(dataSourceName, "gateways.0.description"),
					resource.TestCheckResourceAttrSet(dataSourceName, "gateways.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSourceName, "gateways.0.session_conf.0.tcp_session_expire_time"),
					resource.TestCheckResourceAttrSet(dataSourceName, "gateways.0.session_conf.0.udp_session_expire_time"),
					resource.TestCheckResourceAttrSet(dataSourceName, "gateways.0.dnat_rules_limit"),
					resource.TestCheckResourceAttrSet(dataSourceName, "gateways.0.bps_max"),

					dcByName.CheckResourceExists(),
					resource.TestCheckOutput("name_filter_is_useful", "true"),
					dcNameNotFound.CheckResourceExists(),
					resource.TestCheckOutput("not_found_validation_pass", "true"),

					dcByGatewayId.CheckResourceExists(),
					resource.TestCheckOutput("gateway_id_filter_is_useful", "true"),

					dcBySpec.CheckResourceExists(),
					resource.TestCheckOutput("spec_filter_is_useful", "true"),

					dcByStatus.CheckResourceExists(),
					resource.TestCheckOutput("status_filter_is_useful", "true"),

					dcByVpcId.CheckResourceExists(),
					resource.TestCheckOutput("vpc_id_filter_is_useful", "true"),

					dcBySubnetId.CheckResourceExists(),
					resource.TestCheckOutput("subnet_id_filter_is_useful", "true"),

					dcByEps.CheckResourceExists(),
					resource.TestCheckOutput("eps_filter_is_useful", "true"),

					dcByDescription.CheckResourceExists(),
					resource.TestCheckOutput("description_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccPublicGatewaysDataSource_base() string {
	name := acceptance.RandomAccResourceNameWithDash()

	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_nat_gateway" "test" {
  name                  = "%[2]s"
  spec                  = "1"
  description           = "Created by acc test"
  vpc_id                = huaweicloud_vpc.test.id
  subnet_id             = huaweicloud_vpc_subnet.test.id
  enterprise_project_id = "0"

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, common.TestBaseNetwork(name), name)
}

func testAccDatasourcePublicGateways_basic(baseConfig, name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_nat_gateways" "test" {
  depends_on = [
    huaweicloud_nat_gateway.test
  ]  
}

locals {
  name = data.huaweicloud_nat_gateways.test.gateways[0].name
}

data "huaweicloud_nat_gateways" "filter_by_name" {
  name = local.name
}

locals {
  name_filter_result = [
    for v in data.huaweicloud_nat_gateways.filter_by_name.gateways[*].name : v == local.name
  ]
}

output "name_filter_is_useful" {
  value = alltrue(local.name_filter_result) && length(local.name_filter_result) > 0
}

data "huaweicloud_nat_gateways" "not_found" {
  name = "not_found"
}

output "not_found_validation_pass" {
  value = length(data.huaweicloud_nat_gateways.not_found.gateways) == 0
}

locals {
  gateway_id = data.huaweicloud_nat_gateways.test.gateways[0].id
}

data "huaweicloud_nat_gateways" "filter_by_gateway_id" {
  gateway_id = local.gateway_id
}

locals {
  gateway_id_filter_result = [
    for v in data.huaweicloud_nat_gateways.filter_by_gateway_id.gateways[*].id : v == local.gateway_id
  ]
}

output "gateway_id_filter_is_useful" {
  value = alltrue(local.gateway_id_filter_result) && length(local.gateway_id_filter_result) > 0
}

locals {
  spec = data.huaweicloud_nat_gateways.test.gateways[0].spec
}

data "huaweicloud_nat_gateways" "filter_by_spec" {
  spec = local.spec
}

locals {
  spec_filter_result = [for v in data.huaweicloud_nat_gateways.filter_by_spec.gateways[*].spec : v == local.spec]
}

output "spec_filter_is_useful" {
  value = alltrue(local.spec_filter_result) && length(local.spec_filter_result) > 0
}

locals {
  status = data.huaweicloud_nat_gateways.test.gateways[0].status
}

data "huaweicloud_nat_gateways" "filter_by_status" {
  status = local.status
}

locals {
  status_filter_result = [ 
    for v in data.huaweicloud_nat_gateways.filter_by_status.gateways[*].status : v == local.status
  ]
}

output "status_filter_is_useful" {
  value = alltrue(local.status_filter_result) && length(local.status_filter_result) > 0
}

locals {
  vpc_id = data.huaweicloud_nat_gateways.test.gateways[0].vpc_id
}

data "huaweicloud_nat_gateways" "filter_by_vpc_id" {
  vpc_id = local.vpc_id
}

locals {
  vpc_id_filter_result = [
    for v in data.huaweicloud_nat_gateways.filter_by_vpc_id.gateways[*].vpc_id : v == local.vpc_id
  ]
}

output "vpc_id_filter_is_useful" {
  value = alltrue(local.vpc_id_filter_result) && length(local.vpc_id_filter_result) > 0
}

locals {
  subnet_id = data.huaweicloud_nat_gateways.test.gateways[0].subnet_id
}

data "huaweicloud_nat_gateways" "filter_by_subnet_id" {
  subnet_id = local.subnet_id
}

locals {
  subnet_id_filter_result = [
    for v in data.huaweicloud_nat_gateways.filter_by_subnet_id.gateways[*].subnet_id : v == local.subnet_id
  ]
}

output "subnet_id_filter_is_useful" {
  value = alltrue(local.subnet_id_filter_result) && length(local.subnet_id_filter_result) > 0
}

locals {
  enterprise_project_id = data.huaweicloud_nat_gateways.test.gateways[0].enterprise_project_id
}

data "huaweicloud_nat_gateways" "filter_by_eps" {
  enterprise_project_id = local.enterprise_project_id
}

locals {
  eps_filter_result = [
    for v in data.huaweicloud_nat_gateways.filter_by_eps.gateways[*].enterprise_project_id : 
    v == local.enterprise_project_id
  ]
}

output "eps_filter_is_useful" {
  value = alltrue(local.eps_filter_result) && length(local.eps_filter_result) > 0
}

locals {
  description = data.huaweicloud_nat_gateways.test.gateways[0].description
}

data "huaweicloud_nat_gateways" "filter_by_description" {
  description = local.description
}

locals {
  description_filter_result = [
    for v in data.huaweicloud_nat_gateways.filter_by_description.gateways[*].description : 
    v == local.description
  ]
}

output "description_filter_is_useful" {
  value = alltrue(local.description_filter_result) && length(local.description_filter_result) > 0
}
`, baseConfig, name)
}

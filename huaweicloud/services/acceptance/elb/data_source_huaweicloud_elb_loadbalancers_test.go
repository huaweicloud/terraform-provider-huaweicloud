package elb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccDatasourceLoadBalancers_basic(t *testing.T) {
	rName := "data.huaweicloud_elb_loadbalancers.test"
	dc := acceptance.InitDataSourceCheck(rName)
	name := acceptance.RandomAccResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceLoadBalancers_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "loadbalancers.#"),
					resource.TestCheckResourceAttrSet(rName, "loadbalancers.0.name"),
					resource.TestCheckResourceAttrSet(rName, "loadbalancers.0.id"),
					resource.TestCheckResourceAttrSet(rName, "loadbalancers.0.ipv4_address"),
					resource.TestCheckResourceAttrSet(rName, "loadbalancers.0.ipv4_port_id"),
					resource.TestCheckResourceAttrSet(rName, "loadbalancers.0.l4_flavor_id"),
					resource.TestCheckResourceAttrSet(rName, "loadbalancers.0.l7_flavor_id"),
					resource.TestCheckResourceAttrSet(rName, "loadbalancers.0.vpc_id"),
					resource.TestCheckResourceAttrSet(rName, "loadbalancers.0.enterprise_project_id"),
					resource.TestCheckOutput("name_filter_is_useful", "true"),
					resource.TestCheckOutput("vpc_id_filter_is_useful", "true"),
					resource.TestCheckOutput("ipv4_subnet_id_filter_is_useful", "true"),
					resource.TestCheckOutput("description_filter_is_useful", "true"),
					resource.TestCheckOutput("l4_flavor_id_filter_is_useful", "true"),
					resource.TestCheckOutput("l7_flavor_id_filter_is_useful", "true"),
					resource.TestCheckOutput("type_is_useful", "true"),
					resource.TestCheckOutput("enterprise_project_id_is_useful", "true"),
				),
			},
		},
	})
}

func TestAccDatasourceLoadBalancers_gateway(t *testing.T) {
	rName := "data.huaweicloud_elb_loadbalancers.test"
	dc := acceptance.InitDataSourceCheck(rName)
	name := acceptance.RandomAccResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckElbGatewayType(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceLoadBalancers_gateway(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "loadbalancers.#"),
					resource.TestCheckResourceAttr(rName, "loadbalancers.0.name", name),
					resource.TestCheckResourceAttr(rName, "loadbalancers.0.loadbalancer_type", "gateway"),
					resource.TestCheckResourceAttrSet(rName, "loadbalancers.0.gw_flavor_id"),
				),
			},
		},
	})
}

func testAccDatasourceLoadBalancers_base(rName string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_elb_flavors" "l4flavors" {
  type            = "L4"
  max_connections = 1000000
  cps             = 20000
  bandwidth       = 100
}

data "huaweicloud_elb_flavors" "l7flavors" {
  type            = "L7"
  max_connections = 400000
  cps             = 4000
  bandwidth       = 100
}

resource "huaweicloud_elb_loadbalancer" "test" {
  name           = "%[2]s"
  vpc_id         = huaweicloud_vpc.test.id
  ipv4_subnet_id = huaweicloud_vpc_subnet.test.ipv4_subnet_id
  description    = "update flavors"
  l4_flavor_id   = data.huaweicloud_elb_flavors.l4flavors.ids[0]
  l7_flavor_id   = data.huaweicloud_elb_flavors.l7flavors.ids[0]
  
  availability_zone = [
    data.huaweicloud_availability_zones.test.names[0]
  ]

  backend_subnets = [
    huaweicloud_vpc_subnet.test.id
  ]

  tags = {
    key   = "value"
    owner = "terraform"
  }
}
`, common.TestVpc(rName), rName)
}

func testAccDatasourceLoadBalancers_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_elb_loadbalancers" "test" {
  depends_on = [huaweicloud_elb_loadbalancer.test]
}

data "huaweicloud_elb_loadbalancers" "name_filter" {
  name       = "%[2]s"
  depends_on = [huaweicloud_elb_loadbalancer.test]
}

output "name_filter_is_useful" {
  value = length(data.huaweicloud_elb_loadbalancers.name_filter.loadbalancers) > 0 && alltrue(
  [for v in data.huaweicloud_elb_loadbalancers.name_filter.loadbalancers[*].name :v == "%[2]s"]
  )  
}

data "huaweicloud_elb_loadbalancers" "vpc_id_filter" {
  vpc_id     = huaweicloud_elb_loadbalancer.test.vpc_id
  depends_on = [huaweicloud_elb_loadbalancer.test]
}
locals {
  vpc_id = huaweicloud_elb_loadbalancer.test.vpc_id
}
output "vpc_id_filter_is_useful" {
  value = length(data.huaweicloud_elb_loadbalancers.vpc_id_filter.loadbalancers) > 0 && alltrue(
  [for v in data.huaweicloud_elb_loadbalancers.vpc_id_filter.loadbalancers[*].vpc_id : v == local.vpc_id]
  )  
}

data "huaweicloud_elb_loadbalancers" "ipv4_subnet_id_filter" {
  ipv4_subnet_id = huaweicloud_elb_loadbalancer.test.ipv4_subnet_id
  depends_on     = [huaweicloud_elb_loadbalancer.test]
}
locals {
  ipv4_subnet_id = huaweicloud_elb_loadbalancer.test.ipv4_subnet_id
}
output "ipv4_subnet_id_filter_is_useful" {
  value = length(data.huaweicloud_elb_loadbalancers.vpc_id_filter.loadbalancers) > 0 && alltrue(
  [for v in data.huaweicloud_elb_loadbalancers.ipv4_subnet_id_filter.loadbalancers[*].ipv4_subnet_id : 
  v == local.ipv4_subnet_id]
  )  
}

data "huaweicloud_elb_loadbalancers" "description_filter" {
  description = huaweicloud_elb_loadbalancer.test.description
  depends_on  = [huaweicloud_elb_loadbalancer.test]
}
locals {
  description = huaweicloud_elb_loadbalancer.test.description
}
output "description_filter_is_useful" {
  value = length(data.huaweicloud_elb_loadbalancers.description_filter.loadbalancers) > 0 && alltrue(
  [for v in data.huaweicloud_elb_loadbalancers.description_filter.loadbalancers[*].description : 
  v == local.description]
  )  
}


data "huaweicloud_elb_loadbalancers" "l4_flavor_id_filter" {
  l4_flavor_id = huaweicloud_elb_loadbalancer.test.l4_flavor_id
  depends_on   = [huaweicloud_elb_loadbalancer.test]
}
locals {
  l4_flavor_id = huaweicloud_elb_loadbalancer.test.l4_flavor_id
}
output "l4_flavor_id_filter_is_useful" {
  value = length(data.huaweicloud_elb_loadbalancers.l4_flavor_id_filter.loadbalancers) > 0 && alltrue(
  [for v in data.huaweicloud_elb_loadbalancers.l4_flavor_id_filter.loadbalancers[*].l4_flavor_id : 
  v == local.l4_flavor_id]
  )  
}

data "huaweicloud_elb_loadbalancers" "l7_flavor_id_filter" {
  l7_flavor_id = huaweicloud_elb_loadbalancer.test.l7_flavor_id
  depends_on   = [huaweicloud_elb_loadbalancer.test]
}
locals {
  l7_flavor_id = huaweicloud_elb_loadbalancer.test.l7_flavor_id
}
output "l7_flavor_id_filter_is_useful" {
  value = length(data.huaweicloud_elb_loadbalancers.l7_flavor_id_filter.loadbalancers) > 0 && alltrue(
  [for v in data.huaweicloud_elb_loadbalancers.l7_flavor_id_filter.loadbalancers[*].l7_flavor_id : 
  v == local.l7_flavor_id]
  ) 
}


data "huaweicloud_elb_loadbalancers" "type_filter" {
  type       = "dedicated"
  depends_on = [huaweicloud_elb_loadbalancer.test]
}
locals {
  type       = "dedicated"
}
output "type_is_useful" {
  value = length(data.huaweicloud_elb_loadbalancers.type_filter.loadbalancers) > 0 && alltrue(
  [for v in data.huaweicloud_elb_loadbalancers.type_filter.loadbalancers[*].type : 
  v == local.type]
  ) 
}

data "huaweicloud_elb_loadbalancers" "enterprise_project_id_filter" {
  depends_on = [huaweicloud_elb_loadbalancer.test]
}
locals {
  enterprise_project_id = huaweicloud_elb_loadbalancer.test.enterprise_project_id
}
output "enterprise_project_id_is_useful" {
  value = length(data.huaweicloud_elb_loadbalancers.enterprise_project_id_filter.loadbalancers) > 0 && alltrue(
  [for v in data.huaweicloud_elb_loadbalancers.enterprise_project_id_filter.loadbalancers[*].enterprise_project_id : 
  v == local.enterprise_project_id]
  ) 
}

`, testAccDatasourceLoadBalancers_base(name), name)
}

func testAccDatasourceLoadBalancers_gateway(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_vpc" "test" {
  name = "%[1]s"
  cidr = "192.168.0.0/16"
}

resource "huaweicloud_vpc_subnet" "test" {
  name        = "%[1]s"
  vpc_id      = huaweicloud_vpc.test.id
  cidr        = "192.168.0.0/24"
  gateway_ip  = "192.168.0.1"
  ipv6_enable = true
}

resource "huaweicloud_elb_loadbalancer" "test" {
  name              = "%[1]s"
  vpc_id            = huaweicloud_vpc.test.id
  ipv4_subnet_id    = huaweicloud_vpc_subnet.test.ipv4_subnet_id
  ipv6_network_id   = huaweicloud_vpc_subnet.test.id
  loadbalancer_type = "gateway"
  description       = "test gateway description"
  availability_zone = [
    data.huaweicloud_availability_zones.test.names[0]
  ]
}

data "huaweicloud_elb_loadbalancers" "test" {
  depends_on = [huaweicloud_elb_loadbalancer.test]
  name       = "%[1]s"
}
`, name)
}

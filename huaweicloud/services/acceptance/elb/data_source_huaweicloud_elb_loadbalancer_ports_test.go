package elb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceElbLoadBalancerPorts_basic(t *testing.T) {
	dataSource := "data.huaweicloud_elb_loadbalancer_ports.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckElbLoadbalancerID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceElbLoadBalancerPorts_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "ports.#"),
					resource.TestCheckResourceAttrSet(dataSource, "ports.0.port_id"),
					resource.TestCheckResourceAttrSet(dataSource, "ports.0.ip_address"),
					resource.TestCheckResourceAttrSet(dataSource, "ports.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "ports.0.virsubnet_id"),
					resource.TestCheckOutput("port_id_filter_is_useful", "true"),
					resource.TestCheckOutput("ip_address_filter_is_useful", "true"),
					resource.TestCheckOutput("type_filter_is_useful", "true"),
					resource.TestCheckOutput("virsubnet_id_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceElbLoadBalancerPorts_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_elb_loadbalancer_ports" "test" {
  loadbalancer_id = "%[1]s"
}

locals {
  port_id = data.huaweicloud_elb_loadbalancer_ports.test.ports[0].port_id
}
data "huaweicloud_elb_loadbalancer_ports" "port_id_filter" {
  loadbalancer_id = "%[1]s"
  port_id         = [data.huaweicloud_elb_loadbalancer_ports.test.ports[0].port_id]
}
output "port_id_filter_is_useful" {
  value = length(data.huaweicloud_elb_loadbalancer_ports.port_id_filter.ports) > 0 && alltrue(
    [for v in data.huaweicloud_elb_loadbalancer_ports.port_id_filter.ports[*].port_id : v == local.port_id]
  )
}

locals {
  ip_address = data.huaweicloud_elb_loadbalancer_ports.test.ports[0].ip_address
}
data "huaweicloud_elb_loadbalancer_ports" "ip_address_filter" {
  loadbalancer_id = "%[1]s"
  ip_address      = [data.huaweicloud_elb_loadbalancer_ports.test.ports[0].ip_address]
}
output "ip_address_filter_is_useful" {
  value = length(data.huaweicloud_elb_loadbalancer_ports.ip_address_filter.ports) > 0 && alltrue(
    [for v in data.huaweicloud_elb_loadbalancer_ports.ip_address_filter.ports[*].ip_address : v == local.ip_address]
  )
}

locals {
  type = data.huaweicloud_elb_loadbalancer_ports.test.ports[0].type
}
data "huaweicloud_elb_loadbalancer_ports" "type_filter" {
  loadbalancer_id = "%[1]s"
  type            = [data.huaweicloud_elb_loadbalancer_ports.test.ports[0].type]
}
output "type_filter_is_useful" {
  value = length(data.huaweicloud_elb_loadbalancer_ports.type_filter.ports) > 0 && alltrue(
    [for v in data.huaweicloud_elb_loadbalancer_ports.type_filter.ports[*].type : v == local.type]
  )
}

locals {
  virsubnet_id = data.huaweicloud_elb_loadbalancer_ports.test.ports[0].virsubnet_id
}
data "huaweicloud_elb_loadbalancer_ports" "virsubnet_id_filter" {
  loadbalancer_id = "%[1]s"
  virsubnet_id    = [data.huaweicloud_elb_loadbalancer_ports.test.ports[0].virsubnet_id]
}
output "virsubnet_id_filter_is_useful" {
  value = length(data.huaweicloud_elb_loadbalancer_ports.virsubnet_id_filter.ports) > 0 && alltrue(
    [for v in data.huaweicloud_elb_loadbalancer_ports.virsubnet_id_filter.ports[*].virsubnet_id : v == local.virsubnet_id]
  )
}
`, acceptance.HW_ELB_LOADBALANCER_ID)
}

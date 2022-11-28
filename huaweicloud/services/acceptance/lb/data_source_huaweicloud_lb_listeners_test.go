package lb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceListeners_basic(t *testing.T) {
	var (
		rName            = acceptance.RandomAccResourceNameWithDash()
		dcByName         = acceptance.InitDataSourceCheck("data.huaweicloud_lb_listeners.by_name")
		dcByProtocol     = acceptance.InitDataSourceCheck("data.huaweicloud_lb_listeners.by_protocol")
		dcByProtocolPort = acceptance.InitDataSourceCheck("data.huaweicloud_lb_listeners.by_protocol_port")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceListeners_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dcByName.CheckResourceExists(),
					resource.TestCheckOutput("name_query_result_validation", "true"),
					resource.TestCheckResourceAttrSet("data.huaweicloud_lb_listeners.by_name",
						"listeners.0.name"),
					resource.TestCheckResourceAttrSet("data.huaweicloud_lb_listeners.by_name",
						"listeners.0.protocol"),
					resource.TestCheckResourceAttrSet("data.huaweicloud_lb_listeners.by_name",
						"listeners.0.protocol_port"),
					resource.TestCheckResourceAttrSet("data.huaweicloud_lb_listeners.by_name",
						"listeners.0.connection_limit"),
					resource.TestCheckResourceAttrSet("data.huaweicloud_lb_listeners.by_name",
						"listeners.0.http2_enable"),
					resource.TestCheckResourceAttr("data.huaweicloud_lb_listeners.by_name",
						"listeners.0.loadbalancers.#", "1"),
					dcByProtocol.CheckResourceExists(),
					resource.TestCheckOutput("protocol_query_result_validation", "true"),
					dcByProtocolPort.CheckResourceExists(),
					resource.TestCheckOutput("protocol_port_query_result_validation", "true"),
				),
			},
		},
	})
}

func testAccDatasourceListeners_base(rName string) string {
	rCidr := acceptance.RandomCidr()
	return fmt.Sprintf(`
variable "listener_configuration" {
  type = list(object({
    protocol_port = number
    protocol      = string
  }))
  default = [
    {protocol_port = 306, protocol = "TCP"},
    {protocol_port = 406, protocol = "UDP"},
    {protocol_port = 506, protocol = "HTTP"},
  ]
}

resource "huaweicloud_vpc" "test" {
  name = "%[1]s"
  cidr = "%[2]s"
}

resource "huaweicloud_vpc_subnet" "test" {
  vpc_id = huaweicloud_vpc.test.id

  name       = "%[1]s"
  cidr       = cidrsubnet(huaweicloud_vpc.test.cidr, 4, 1)
  gateway_ip = cidrhost(cidrsubnet(huaweicloud_vpc.test.cidr, 4, 1), 1)
}

resource "huaweicloud_lb_loadbalancer" "test" {
  name          = "%[1]s"
  vip_subnet_id = huaweicloud_vpc_subnet.test.ipv4_subnet_id
}

resource "huaweicloud_lb_listener" "test" {
  count = length(var.listener_configuration)

  loadbalancer_id = huaweicloud_lb_loadbalancer.test.id

  name          = "%[1]s-${count.index}"
  protocol      = var.listener_configuration[count.index]["protocol"]
  protocol_port = var.listener_configuration[count.index]["protocol_port"]
}
`, rName, rCidr)
}

func testAccDatasourceListeners_basic(rName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_lb_listeners" "by_name" {
  depends_on = [huaweicloud_lb_listener.test]

  name = huaweicloud_lb_listener.test[0].name
}

data "huaweicloud_lb_listeners" "by_protocol" {
  depends_on = [huaweicloud_lb_listener.test]

  protocol = huaweicloud_lb_listener.test[1].protocol
}

data "huaweicloud_lb_listeners" "by_protocol_port" {
  depends_on = [huaweicloud_lb_listener.test]

  protocol_port = huaweicloud_lb_listener.test[2].protocol_port
}

output "name_query_result_validation" {
  value = contains(data.huaweicloud_lb_listeners.by_name.listeners[*].id,
  huaweicloud_lb_listener.test[0].id) && !contains(data.huaweicloud_lb_listeners.by_name.listeners[*].id,
  huaweicloud_lb_listener.test[1].id) && !contains(data.huaweicloud_lb_listeners.by_name.listeners[*].id,
  huaweicloud_lb_listener.test[2].id)
}

output "protocol_query_result_validation" {
  value = contains(data.huaweicloud_lb_listeners.by_protocol.listeners[*].id,
  huaweicloud_lb_listener.test[1].id) && !contains(data.huaweicloud_lb_listeners.by_protocol.listeners[*].id,
  huaweicloud_lb_listener.test[0].id) && !contains(data.huaweicloud_lb_listeners.by_protocol.listeners[*].id,
  huaweicloud_lb_listener.test[2].id)
}

output "protocol_port_query_result_validation" {
  value = contains(data.huaweicloud_lb_listeners.by_protocol_port.listeners[*].id,
  huaweicloud_lb_listener.test[2].id) && !contains(data.huaweicloud_lb_listeners.by_protocol_port.listeners[*].id,
  huaweicloud_lb_listener.test[0].id) && !contains(data.huaweicloud_lb_listeners.by_protocol_port.listeners[*].id,
  huaweicloud_lb_listener.test[1].id)
}
`, testAccDatasourceListeners_base(rName))
}

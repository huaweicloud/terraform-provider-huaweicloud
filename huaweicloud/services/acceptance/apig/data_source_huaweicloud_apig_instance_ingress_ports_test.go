package apig

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataInstanceIngressPorts_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_apig_instance_ingress_ports.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)

		byProtocol   = "data.huaweicloud_apig_instance_ingress_ports.filter_by_protocol"
		dcByProtocol = acceptance.InitDataSourceCheck(byProtocol)

		byPort   = "data.huaweicloud_apig_instance_ingress_ports.filter_by_port"
		dcByPort = acceptance.InitDataSourceCheck(byPort)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckApigSubResourcesRelatedInfo(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataInstanceIngressPorts_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dataSource, "ingress_ports.#", regexp.MustCompile(`^[0-9]+$`)),
					resource.TestCheckResourceAttrSet(dataSource, "ingress_ports.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "ingress_ports.0.protocol"),
					resource.TestCheckResourceAttrSet(dataSource, "ingress_ports.0.port"),
					resource.TestCheckResourceAttrSet(dataSource, "ingress_ports.0.status"),
					dcByProtocol.CheckResourceExists(),
					resource.TestCheckOutput("protocol_filter_is_useful", "true"),
					dcByPort.CheckResourceExists(),
					resource.TestCheckOutput("port_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDataInstanceIngressPorts_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_apig_instances" "test" {
  instance_id = "%[1]s"
}

locals {
  instance_id = data.huaweicloud_apig_instances.test.instances[0].id
}

resource "huaweicloud_apig_instance_ingress_port" "http" {
  instance_id = local.instance_id
  protocol    = "HTTP"
  port        = 8080
}

resource "huaweicloud_apig_instance_ingress_port" "https" {
  instance_id = local.instance_id
  protocol    = "HTTPS"
  port        = 8443
}

data "huaweicloud_apig_instance_ingress_ports" "test" {
  depends_on = [
    huaweicloud_apig_instance_ingress_port.http,
    huaweicloud_apig_instance_ingress_port.https
  ]

  instance_id = local.instance_id
}

# Filter by protocol
locals {
  protocol = huaweicloud_apig_instance_ingress_port.http.protocol
}

data "huaweicloud_apig_instance_ingress_ports" "filter_by_protocol" {
  depends_on = [
    huaweicloud_apig_instance_ingress_port.http,
    huaweicloud_apig_instance_ingress_port.https
  ]

  instance_id = local.instance_id
  protocol    = local.protocol
}

locals {
  protocol_filter_result = [
    for v in data.huaweicloud_apig_instance_ingress_ports.filter_by_protocol.ingress_ports[*].protocol : v == local.protocol
  ]
}

output "protocol_filter_is_useful" {
  value = length(local.protocol_filter_result) > 0 && alltrue(local.protocol_filter_result)
}

# Filter by port
locals {
  port = huaweicloud_apig_instance_ingress_port.http.port
}

data "huaweicloud_apig_instance_ingress_ports" "filter_by_port" {
  depends_on = [
    huaweicloud_apig_instance_ingress_port.http,
    huaweicloud_apig_instance_ingress_port.https
  ]

  instance_id = local.instance_id
  port        = local.port
}

locals {
  port_filter_result = [
    for v in data.huaweicloud_apig_instance_ingress_ports.filter_by_port.ingress_ports[*].port : v == local.port
  ]
}

output "port_filter_is_useful" {
  value = length(local.port_filter_result) > 0 && alltrue(local.port_filter_result)
}
`, acceptance.HW_APIG_DEDICATED_INSTANCE_ID)
}
